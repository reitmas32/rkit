// loki_buffered_hook.go
package hooks

import (
	"bytes"
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

	// BatchSize is how many logs trigger an immediate POST flush.
	// If BatchSize <= 0, it defaults to 50.
	BatchSize int

	// BaseLabels are added to every stream (e.g., app/env/service).
	BaseLabels map[string]string

	// Optional headers (e.g. X-Scope-OrgID, Authorization).
	Headers map[string]string

	// Optional: if true, add entry.Data["caller"] as a Loki label (high cardinality risk).
	IncludeCallerAsLabel bool

	// Optional: if true, include entry.Data serialized into the log line as JSON.
	IncludeFieldsInLine bool

	client *http.Client

	mu     sync.Mutex
	buffer []lokiEntry
}

// lokiEntry is the internal buffered representation.
type lokiEntry struct {
	Time   time.Time
	Level  string
	Msg    string
	Fields logrus.Fields
}

// NewLokiBufferedHook creates a buffered Loki hook.
// url: full push endpoint.
// batchSize: number of logs per batch (default 50 if <= 0).
// baseLabels: fixed labels for all streams.
func NewLokiBufferedHook(url string, batchSize int, baseLabels map[string]string) *LokiBufferedHook {
	if batchSize <= 0 {
		batchSize = 50
	}
	if baseLabels == nil {
		baseLabels = map[string]string{}
	}

	return &LokiBufferedHook{
		URL:       url,
		BatchSize: batchSize,
		BaseLabels: func() map[string]string {
			cp := make(map[string]string, len(baseLabels))
			for k, v := range baseLabels {
				cp[k] = v
			}
			return cp
		}(),
		Headers: map[string]string{},
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		buffer: make([]lokiEntry, 0, batchSize),
	}
}

// Levels implements logrus.Hook.
func (h *LokiBufferedHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implements logrus.Hook.
func (h *LokiBufferedHook) Fire(e *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.buffer = append(h.buffer, lokiEntry{
		Time:  e.Time,
		Level: e.Level.String(),
		Msg:   e.Message,
		Fields: func() logrus.Fields {
			// Copy map to avoid mutation issues
			cp := make(logrus.Fields, len(e.Data))
			for k, v := range e.Data {
				cp[k] = v
			}
			return cp
		}(),
	})

	if len(h.buffer) >= h.BatchSize {
		// Best-effort flush; keep buffer if it fails.
		_ = h.flushLocked()
	}

	return nil
}

// Flush triggers a manual flush (best-effort).
func (h *LokiBufferedHook) Flush() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.flushLocked()
}

// flushLocked sends current buffer to Loki.
// Caller must hold h.mu.
func (h *LokiBufferedHook) flushLocked() error {
	if len(h.buffer) == 0 {
		return nil
	}
	if h.URL == "" {
		return errors.New("loki hook: URL is empty")
	}

	logCount := len(h.buffer)
	flushTime := time.Now()
	fmt.Printf("[loki] flushing %d logs at %s to %s\n", logCount, flushTime.Format(time.RFC3339), h.URL)

	payload := h.buildPayload(h.buffer)

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("[loki] flush failed before request: %v\n", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, h.URL, bytes.NewReader(body))
	if err != nil {
		fmt.Printf("[loki] flush failed creating request: %v\n", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range h.Headers {
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		// keep buffer, let next flush retry
		fmt.Printf("[loki] flush failed sending %d logs to %s: %v\n", logCount, h.URL, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// read small body for debugging
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		// keep buffer
		err := errors.New("loki hook: non-2xx response: " + resp.Status + " body=" + string(b))
		fmt.Printf("[loki] flush failed sending %d logs to %s: status=%s body=%s\n", logCount, h.URL, resp.Status, string(b))
		return err
	}

	// success: clear buffer
	h.buffer = h.buffer[:0]
	fmt.Printf("[loki] flush succeeded sending %d logs to %s: status=%s\n", logCount, h.URL, resp.Status)
	return nil
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
