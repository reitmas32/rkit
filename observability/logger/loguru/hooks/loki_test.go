package hooks_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/reitmas32/rkit/observability/logger/loguru/hooks"
	"github.com/sirupsen/logrus"
)

func entry(msg string) *logrus.Entry {
	return &logrus.Entry{Time: time.Now(), Level: logrus.InfoLevel, Message: msg, Data: logrus.Fields{}}
}

func TestLokiHook_BufferCapDropsOldest(t *testing.T) {
	// Large batch size so nothing auto-flushes; a small cap so the buffer is bounded.
	h := hooks.NewLokiBufferedHook("http://127.0.0.1:1/loki/api/v1/push", 1000, nil)
	h.MaxBufferSize = 3
	for i := 0; i < 5; i++ {
		_ = h.Fire(entry("m"))
	}
	if got := h.Dropped(); got != 2 {
		t.Errorf("Dropped() = %d, want 2 (buffer capped at 3 after 5 entries)", got)
	}
}

func TestLokiHook_FlushSendsBatch(t *testing.T) {
	var posts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&posts, 1)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	h := hooks.NewLokiBufferedHook(srv.URL+"/loki/api/v1/push", 2, map[string]string{"app": "test"})
	_ = h.Fire(entry("a"))
	_ = h.Fire(entry("b")) // reaches batch size 2 -> synchronous auto-flush

	if got := atomic.LoadInt32(&posts); got != 1 {
		t.Errorf("expected exactly 1 POST after a full batch, got %d", got)
	}
}

func TestLokiHook_OnErrorCalledOnFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	var called atomic.Bool
	h := hooks.NewLokiBufferedHook(srv.URL+"/loki/api/v1/push", 100, nil)
	h.OnError = func(error) { called.Store(true) }

	_ = h.Fire(entry("x"))
	if err := h.Flush(); err == nil {
		t.Error("Flush should return an error on a non-2xx response")
	}
	if !called.Load() {
		t.Error("OnError should have been invoked on failure")
	}
	// The failed batch is re-queued, so a later flush can retry it.
	if h.Dropped() != 0 {
		t.Errorf("no entries should be dropped below MaxBufferSize, got %d", h.Dropped())
	}
}
