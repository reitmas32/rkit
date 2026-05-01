package http

import (
	"io"
)

// Response represents an HTTP response.
type Response struct {
	// StatusCode is the HTTP status code (200, 404, 500, etc.)
	StatusCode int

	// Status is the status text (e.g., "200 OK")
	Status string

	// Headers contains the response headers
	Headers map[string]string

	// Body is the response body
	Body io.ReadCloser

	// ContentLength is the length of the response body
	ContentLength int64
}

// ReadBody reads the entire response body and returns it as bytes.
// After calling this method, the Body will be closed.
func (r *Response) ReadBody() ([]byte, error) {
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

// ReadBodyString reads the entire response body and returns it as a string.
// After calling this method, the Body will be closed.
func (r *Response) ReadBodyString() (string, error) {
	body, err := r.ReadBody()
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Close closes the response body.
func (r *Response) Close() error {
	if r.Body != nil {
		return r.Body.Close()
	}
	return nil
}

// IsSuccess returns true if the status code is in the 2xx range.
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError returns true if the status code is in the 4xx range.
func (r *Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError returns true if the status code is in the 5xx range.
func (r *Response) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}
