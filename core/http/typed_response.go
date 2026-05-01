package http

import (
	"encoding/json"
	"io"
	"time"
)

// TypedResponse represents a typed HTTP response with parsed body.
// T is the type of the response body.
type TypedResponse[T any] struct {
	// StatusCode is the HTTP status code (200, 404, 500, etc.)
	StatusCode int

	// Status is the status text (e.g., "200 OK")
	Status string

	// Headers contains the response headers
	Headers map[string]string

	// Body is the parsed response body of type T
	Body T

	// RawBody contains the raw response body as bytes (useful for binary data or debugging)
	RawBody []byte

	// ContentType is the Content-Type header value
	ContentType string

	// ContentLength is the length of the response body
	ContentLength int64

	// RequestTime is when the request was initiated
	RequestTime time.Time

	// ResponseTime is when the response was received
	ResponseTime time.Time

	// Duration is the time it took to complete the request
	Duration time.Duration

	// ExpectedStatusCode is an optional expected status code.
	// If set, IsSuccess() will return true only if StatusCode matches this value.
	ExpectedStatusCode *int

	// SuccessStatusCodeRange is an optional range of status codes considered successful.
	// If set, IsSuccess() will return true if StatusCode is within [Min, Max] (inclusive).
	// If both ExpectedStatusCode and SuccessStatusCodeRange are set, ExpectedStatusCode takes precedence.
	SuccessStatusCodeRange *StatusCodeRange
}

// StatusCodeRange represents a range of HTTP status codes.
type StatusCodeRange struct {
	Min int // Minimum status code (inclusive)
	Max int // Maximum status code (inclusive)
}

// IsSuccess returns true if the status code is considered successful.
// It checks in the following order:
// 1. If ExpectedStatusCode is set, returns true only if StatusCode matches it.
// 2. If SuccessStatusCodeRange is set, returns true if StatusCode is within [Min, Max] (inclusive).
// 3. Otherwise, returns true if StatusCode is in the 2xx range (default behavior).
func (r *TypedResponse[T]) IsSuccess() bool {
	// Check expected status code first (highest priority)
	if r.ExpectedStatusCode != nil {
		return r.StatusCode == *r.ExpectedStatusCode
	}

	// Check status code range
	if r.SuccessStatusCodeRange != nil {
		return r.StatusCode >= r.SuccessStatusCodeRange.Min && r.StatusCode <= r.SuccessStatusCodeRange.Max
	}

	// Default: 2xx range
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError returns true if the status code is in the 4xx range.
func (r *TypedResponse[T]) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError returns true if the status code is in the 5xx range.
func (r *TypedResponse[T]) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}

// ParseResponse parses a Response into a TypedResponse with the given type T.
// It attempts to parse the body as JSON. If T is []byte, it returns the raw body.
// requestTime is the time when the request was initiated (used to calculate duration).
func ParseResponse[T any](resp *Response, requestTime time.Time) (*TypedResponse[T], error) {
	responseTime := time.Now()
	duration := responseTime.Sub(requestTime)

	// Read raw body
	rawBody, err := resp.ReadBody()
	if err != nil {
		return nil, err
	}

	// Get content type
	contentType := resp.Headers["Content-Type"]
	if contentType == "" {
		contentType = resp.Headers["content-type"]
	}

	var body T

	// Handle []byte type specially (for binary data)
	// Use type assertion to check if T is []byte
	var zeroT T
	if _, ok := any(zeroT).([]byte); ok {
		// T is []byte, assign raw body directly
		body = any(rawBody).(T)
	} else {
		// Try to parse as JSON
		if len(rawBody) > 0 {
			if err := json.Unmarshal(rawBody, &body); err != nil {
				return nil, err
			}
		}
	}

	return &TypedResponse[T]{
		StatusCode:            resp.StatusCode,
		Status:                resp.Status,
		Headers:               resp.Headers,
		Body:                  body,
		RawBody:               rawBody,
		ContentType:           contentType,
		ContentLength:         resp.ContentLength,
		RequestTime:           requestTime,
		ResponseTime:          responseTime,
		Duration:              duration,
		ExpectedStatusCode:    nil,
		SuccessStatusCodeRange: nil,
	}, nil
}

// ParseResponseFromReader parses a response from a reader with the given status code and headers.
// Useful for custom parsing scenarios.
// requestTime is the time when the request was initiated (used to calculate duration).
func ParseResponseFromReader[T any](reader io.ReadCloser, statusCode int, status string, headers map[string]string, requestTime time.Time) (*TypedResponse[T], error) {
	responseTime := time.Now()
	duration := responseTime.Sub(requestTime)

	rawBody, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	contentType := headers["Content-Type"]
	if contentType == "" {
		contentType = headers["content-type"]
	}

	var body T

	// Handle []byte type specially (for binary data)
	// Use type assertion to check if T is []byte
	var zeroT T
	if _, ok := any(zeroT).([]byte); ok {
		// T is []byte, assign raw body directly
		body = any(rawBody).(T)
	} else {
		// Try to parse as JSON
		if len(rawBody) > 0 {
			if err := json.Unmarshal(rawBody, &body); err != nil {
				return nil, err
			}
		}
	}

	return &TypedResponse[T]{
		StatusCode:            statusCode,
		Status:                status,
		Headers:               headers,
		Body:                  body,
		RawBody:               rawBody,
		ContentType:           contentType,
		ContentLength:         int64(len(rawBody)),
		RequestTime:           requestTime,
		ResponseTime:          responseTime,
		Duration:              duration,
		ExpectedStatusCode:    nil,
		SuccessStatusCodeRange: nil,
	}, nil
}
