package http

import (
	"io"
	"strings"
	"testing"
	"time"
)

func TestNewClient_TimeoutClamp(t *testing.T) {
	if got := NewClient(Config{Timeout: 0}).httpClient.Timeout; got != DefaultTimeoutSeconds*time.Second {
		t.Errorf("zero Timeout should clamp to default, got %v", got)
	}
	if got := NewClient(Config{Timeout: -1}).httpClient.Timeout; got != 0 {
		t.Errorf("negative Timeout is explicit-unlimited, got %v", got)
	}
	if got := NewClient(Config{Timeout: 7}).httpClient.Timeout; got != 7*time.Second {
		t.Errorf("explicit Timeout not honored, got %v", got)
	}
	if got := NewClient(Config{}).config.MaxResponseBytes; got != DefaultMaxResponseBytes {
		t.Errorf("zero MaxResponseBytes should clamp to default, got %d", got)
	}
}

func TestMaxBytesReadCloser(t *testing.T) {
	// A body exactly at the limit reads fine.
	body := io.NopCloser(strings.NewReader("hello")) // 5 bytes
	out, err := io.ReadAll(newLimitedBody(body, 5))
	if err != nil || string(out) != "hello" {
		t.Errorf("at-limit read failed: out=%q err=%v", out, err)
	}

	// One byte over the limit fails with ErrResponseTooLarge.
	body = io.NopCloser(strings.NewReader("hello!")) // 6 bytes
	if _, err := io.ReadAll(newLimitedBody(body, 5)); err != ErrResponseTooLarge {
		t.Errorf("over-limit read should return ErrResponseTooLarge, got %v", err)
	}

	// limit <= 0 means unlimited.
	body = io.NopCloser(strings.NewReader("anything"))
	if out, err := io.ReadAll(newLimitedBody(body, 0)); err != nil || string(out) != "anything" {
		t.Errorf("unlimited read failed: out=%q err=%v", out, err)
	}
}
