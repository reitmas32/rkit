package http

import (
	"errors"
	"io"
)

// ErrResponseTooLarge is returned when a response body exceeds the client's
// configured MaxResponseBytes while being read. It guards against memory
// exhaustion from oversized or malicious responses.
var ErrResponseTooLarge = errors.New("http: response body exceeds configured MaxResponseBytes limit")

// maxBytesReadCloser wraps a ReadCloser and fails with ErrResponseTooLarge once
// more than limit bytes have been read. At most one extra byte beyond the limit
// is read, so memory stays bounded.
type maxBytesReadCloser struct {
	r     io.ReadCloser
	limit int64
	read  int64
}

// newLimitedBody wraps body so that reads beyond limit fail. A limit <= 0 means
// unlimited and returns the original body unchanged.
func newLimitedBody(body io.ReadCloser, limit int64) io.ReadCloser {
	if body == nil || limit <= 0 {
		return body
	}
	return &maxBytesReadCloser{r: body, limit: limit}
}

func (m *maxBytesReadCloser) Read(p []byte) (int, error) {
	if m.read > m.limit {
		return 0, ErrResponseTooLarge
	}
	// Read at most one byte past the limit so we can detect overflow without
	// pulling in a large excess.
	if max := m.limit - m.read + 1; int64(len(p)) > max {
		p = p[:max]
	}
	n, err := m.r.Read(p)
	m.read += int64(n)
	if m.read > m.limit {
		return n, ErrResponseTooLarge
	}
	return n, err
}

func (m *maxBytesReadCloser) Close() error { return m.r.Close() }
