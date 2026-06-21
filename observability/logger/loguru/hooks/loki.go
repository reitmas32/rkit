// loki_buffered_hook.go
package hooks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// LokiBufferedHook buffers log entries and sends them to Loki in batches.
// It posts to Loki's push endpoint: /loki/api/v1/push
//
// Behavior:
// - Every logrus.Entry received in Fire() is appended to an in-memory buffer.
// - When the buffer reaches BatchSize, it flushes and POSTs a single payload to Loki.
// - If the POST fails (network or non-2xx), the buffer is kept (no data loss).
//
// Notes:
//   - This implementation groups entries into Loki "streams" based on labels.
//   - BaseLabels are always included, and a "level" label is automatically added.
//   - Optionally, if the entry has field "caller", it can be added as a label.
//     (Be careful: high-cardinality labels can hurt Loki performance.)
type LokiBufferedHook struct {
	// URL is the full Loki push endpoint, e.g. http://loki:3100/loki/api/v1/push
	URL string

	// BatchSize is how many buffered logs trigger an immediate POST flush.
	// If BatchSize <= 0, it defaults to DefaultBatchSize.
	BatchSize int

	// MaxBufferSize caps the number of buffered entries. When the buffer would
	// exceed it (typically because Loki is unreachable and flushes keep failing),
	// the oldest entries are dropped so memory stays bounded. A value <= 0 means
	// unbounded (not recommended). Drops are counted (see Dropped) and reported
	// through OnError.
	MaxBufferSize int

	// Timeout is the HTTP timeout applied to each push request.
	// If <= 0, DefaultTimeout is used.
	Timeout time.Duration

	// BaseLabels are added to every stream (e.g., app/env/service).
	BaseLabels map[string]string

	// Optional headers (e.g. X-Scope-OrgID, Authorization).
	Headers map[string]string

	// OnError, if set, is called with any flush error (network/non-2xx) and with a
	// notice whenever entries are dropped due to MaxBufferSize. It replaces stdout
	// debug printing; if nil, failures are silent (Flush still returns the error).
	OnError func(error)

	// IncludeCallerAsLabel: add entry.Data["caller"] as a Loki label (high cardinality risk).
	IncludeCallerAsLabel bool

	// IncludeFieldsInLine: include entry.Data serialized into the log line as JSON.
	IncludeFieldsInLine bool

	client  *http.Client
	mu      sync.Mutex
	buffer  []lokiEntry
	dropped int64
}

// lokiEntry is the internal buffered representation.
type lokiEntry struct {
	Time   time.Time
	Level  string
	Msg    string
	Fields logrus.Fields
}

// Defaults applied by NewLokiBufferedHook (all fields remain configurable).
const (
	// DefaultBatchSize is used when BatchSize <= 0.
	DefaultBatchSize = 50
	// DefaultMaxBufferSize bounds the buffer to avoid unbounded memory growth
	// when Loki is unreachable.
	DefaultMaxBufferSize = 10000
	// DefaultTimeout is the per-push HTTP timeout when Timeout <= 0.
	DefaultTimeout = 5 * time.Second
)

// NewLokiBufferedHook creates a buffered Loki hook with safe, overridable
// defaults. After construction you may tune any exported field (MaxBufferSize,
// Timeout, OnError, Headers, ...) before adding the hook to a logger.
//
// url: full push endpoint.
// batchSize: number of logs per batch (DefaultBatchSize if <= 0).
// baseLabels: fixed labels for all streams.
func NewLokiBufferedHook(url string, batchSize int, baseLabels map[string]string) *LokiBufferedHook {
	if batchSize <= 0 {
		batchSize = DefaultBatchSize
	}
	cp := make(map[string]string, len(baseLabels))
	for k, v := range baseLabels {
		cp[k] = v
	}

	return &LokiBufferedHook{
		URL:           url,
		BatchSize:     batchSize,
		MaxBufferSize: DefaultMaxBufferSize,
		Timeout:       DefaultTimeout,
		BaseLabels:    cp,
		Headers:       map[string]string{},
		client:        &http.Client{},
		buffer:        make([]lokiEntry, 0, batchSize),
	}
}

// Levels implements logrus.Hook.
func (h *LokiBufferedHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Dropped returns the total number of entries dropped because the buffer reached
// MaxBufferSize (typically while Loki was unreachable).
func (h *LokiBufferedHook) Dropped() int64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.dropped
}

// Fire implements logrus.Hook. It appends the entry under the lock and, once the
// batch size is reached, flushes. The HTTP POST happens outside the lock, so a
// slow or unreachable Loki never blocks the goroutines that are logging.
func (h *LokiBufferedHook) Fire(e *logrus.Entry) error {
	h.mu.Lock()
	h.buffer = append(h.buffer, lokiEntry{
		Time:   e.Time,
		Level:  e.Level.String(),
		Msg:    e.Message,
		Fields: copyFields(e.Data),
	})
	h.enforceCapLocked()
	full := len(h.buffer) >= h.BatchSize
	h.mu.Unlock()

	if full {
		_ = h.flush() // best-effort; errors are surfaced via OnError
	}
	return nil
}

// Flush sends any buffered entries to Loki (best-effort). It is safe to call at
// the end of a request or on shutdown, and it does not hold the lock during the
// network call.
func (h *LokiBufferedHook) Flush() error {
	return h.flush()
}

// enforceCapLocked drops the oldest entries when the buffer exceeds MaxBufferSize.
// Caller must hold h.mu.
func (h *LokiBufferedHook) enforceCapLocked() {
	if h.MaxBufferSize <= 0 || len(h.buffer) <= h.MaxBufferSize {
		return
	}
	drop := len(h.buffer) - h.MaxBufferSize
	h.buffer = append(h.buffer[:0], h.buffer[drop:]...)
	h.dropped += int64(drop)
	if h.OnError != nil {
		h.OnError(fmt.Errorf("loki hook: buffer full (max %d), dropped %d oldest entries", h.MaxBufferSize, drop))
	}
}

// flush snapshots the buffer under the lock, releases the lock, then POSTs the
// batch. On failure the batch is re-queued in front of newer entries (respecting
// MaxBufferSize) so it is retried on the next flush.
func (h *LokiBufferedHook) flush() error {
	h.mu.Lock()
	if len(h.buffer) == 0 {
		h.mu.Unlock()
		return nil
	}
	if h.URL == "" {
		h.mu.Unlock()
		return errors.New("loki hook: URL is empty")
	}
	batch := h.buffer
	h.buffer = make([]lokiEntry, 0, h.BatchSize)
	h.mu.Unlock()

	if err := h.send(batch); err != nil {
		h.mu.Lock()
		h.buffer = append(batch, h.buffer...)
		h.enforceCapLocked()
		h.mu.Unlock()
		if h.OnError != nil {
			h.OnError(err)
		}
		return err
	}
	return nil
}

// send POSTs a batch to Loki. It neither holds the lock nor touches the buffer.
func (h *LokiBufferedHook) send(batch []lokiEntry) error {
	body, err := json.Marshal(h.buildPayload(batch))
	if err != nil {
		return fmt.Errorf("loki hook: marshal payload: %w", err)
	}

	timeout := h.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("loki hook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range h.Headers {
		req.Header.Set(k, v)
	}

	client := h.client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("loki hook: post %d entries to %s: %w", len(batch), h.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return fmt.Errorf("loki hook: non-2xx from %s: %s body=%s", h.URL, resp.Status, string(b))
	}
	return nil
}

// copyFields returns a copy of the logrus fields so later mutation of the entry
// cannot race with buffered data.
func copyFields(in logrus.Fields) logrus.Fields {
	cp := make(logrus.Fields, len(in))
	for k, v := range in {
		cp[k] = v
	}
	return cp
}

// buildPayload groups entries by labels into Loki streams.
func (h *LokiBufferedHook) buildPayload(entries []lokiEntry) map[string]any {
	// Map[labelsKey] -> values
	// values: [["ts","line"], ...]
	streams := make(map[string][][]string)

	// Also keep parsed labels per key
	labelsByKey := make(map[string]map[string]string)

	for _, en := range entries {
		// Start with base labels
		labels := make(map[string]string, len(h.BaseLabels)+len(en.Fields)+4)
		for k, v := range h.BaseLabels {
			labels[k] = v
		}

		// Add level as label
		labels["level"] = en.Level

		// Add all fields as labels (except caller if IncludeCallerAsLabel is false)
		for k, v := range en.Fields {
			// Skip caller if IncludeCallerAsLabel is false
			if k == "caller" && !h.IncludeCallerAsLabel {
				continue
			}
			// Skip empty values
			if v == nil {
				continue
			}
			// Convert to string and add as label
			labels[k] = toString(v)
		}

		// Stable key: JSON of labels map.
		// Note: json.Marshal on a map is stable in Go? It sorts keys since Go 1.20+ encoder sorts map keys.
		// We'll rely on that; if you want 100% control, build your own canonical string.
		keyBytes, _ := json.Marshal(labels)
		key := string(keyBytes)

		labelsByKey[key] = labels

		ts := strconv.FormatInt(en.Time.UnixNano(), 10)
		line := en.Msg

		if h.IncludeFieldsInLine && len(en.Fields) > 0 {
			// Attach fields as compact JSON at end
			if b, err := json.Marshal(en.Fields); err == nil {
				line = line + " " + string(b)
			}
		}

		streams[key] = append(streams[key], []string{ts, line})
	}

	outStreams := make([]map[string]any, 0, len(streams))
	for key, values := range streams {
		outStreams = append(outStreams, map[string]any{
			"stream": labelsByKey[key],
			"values": values,
		})
	}

	return map[string]any{
		"streams": outStreams,
	}
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
