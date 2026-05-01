package http

import (
	"context"
	"io"
	"net/http"
)

// Request represents an HTTP request.
type Request struct {
	// Method is the HTTP method (GET, POST, PUT, DELETE, etc.)
	Method string

	// URL is the request URL
	URL string

	// Headers contains the request headers
	Headers map[string]string

	// Body is the request body
	Body io.Reader

	// ContentType specifies the Content-Type header
	ContentType string

	// QueryParams contains query parameters
	QueryParams map[string]string

	// Timeout specifies the request timeout (optional, uses client default if not set)
	Timeout *int // in seconds
}

// RequestOption is a function that modifies a Request.
type RequestOption func(*Request)

// WithHeader adds a header to the request.
func WithHeader(key, value string) RequestOption {
	return func(r *Request) {
		if r.Headers == nil {
			r.Headers = make(map[string]string)
		}
		r.Headers[key] = value
	}
}

// WithHeaders adds multiple headers to the request.
func WithHeaders(headers map[string]string) RequestOption {
	return func(r *Request) {
		if r.Headers == nil {
			r.Headers = make(map[string]string)
		}
		for k, v := range headers {
			r.Headers[k] = v
		}
	}
}

// WithQueryParam adds a query parameter to the request.
func WithQueryParam(key, value string) RequestOption {
	return func(r *Request) {
		if r.QueryParams == nil {
			r.QueryParams = make(map[string]string)
		}
		r.QueryParams[key] = value
	}
}

// WithQueryParams adds multiple query parameters to the request.
func WithQueryParams(params map[string]string) RequestOption {
	return func(r *Request) {
		if r.QueryParams == nil {
			r.QueryParams = make(map[string]string)
		}
		for k, v := range params {
			r.QueryParams[k] = v
		}
	}
}

// WithContentType sets the Content-Type header.
func WithContentType(contentType string) RequestOption {
	return func(r *Request) {
		r.ContentType = contentType
	}
}

// WithTimeout sets the request timeout in seconds.
func WithTimeout(seconds int) RequestOption {
	return func(r *Request) {
		r.Timeout = &seconds
	}
}

// NewRequest creates a new Request with the given method and URL.
func NewRequest(method, url string, opts ...RequestOption) *Request {
	req := &Request{
		Method:      method,
		URL:         url,
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}

// ToHTTPRequest converts this Request to a standard library http.Request.
func (r *Request) ToHTTPRequest(ctx context.Context) (*http.Request, error) {
	// Build URL with query parameters
	url := r.URL
	if len(r.QueryParams) > 0 {
		url += "?"
		first := true
		for k, v := range r.QueryParams {
			if !first {
				url += "&"
			}
			url += k + "=" + v
			first = false
		}
	}

	// Create http.Request
	httpReq, err := http.NewRequestWithContext(ctx, r.Method, url, r.Body)
	if err != nil {
		return nil, err
	}

	// Set headers
	if r.ContentType != "" {
		httpReq.Header.Set("Content-Type", r.ContentType)
	}
	for k, v := range r.Headers {
		httpReq.Header.Set(k, v)
	}

	return httpReq, nil
}
