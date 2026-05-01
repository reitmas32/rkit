package http

import (
	"context"
	"strings"
	"testing"
)

func TestWithHeader(t *testing.T) {
	req := NewRequest("GET", "https://example.com", WithHeader("Authorization", "Bearer token"))

	if req.Headers["Authorization"] != "Bearer token" {
		t.Errorf("WithHeader() = %v, want 'Bearer token'", req.Headers["Authorization"])
	}
}

func TestWithHeader_Multiple(t *testing.T) {
	req := NewRequest("GET", "https://example.com",
		WithHeader("Authorization", "Bearer token"),
		WithHeader("Accept", "application/json"),
	)

	if req.Headers["Authorization"] != "Bearer token" {
		t.Errorf("WithHeader() Authorization = %v, want 'Bearer token'", req.Headers["Authorization"])
	}
	if req.Headers["Accept"] != "application/json" {
		t.Errorf("WithHeader() Accept = %v, want 'application/json'", req.Headers["Accept"])
	}
}

func TestWithHeaders(t *testing.T) {
	headers := map[string]string{
		"Authorization": "Bearer token",
		"Accept":        "application/json",
	}
	req := NewRequest("GET", "https://example.com", WithHeaders(headers))

	if req.Headers["Authorization"] != "Bearer token" {
		t.Errorf("WithHeaders() Authorization = %v, want 'Bearer token'", req.Headers["Authorization"])
	}
	if req.Headers["Accept"] != "application/json" {
		t.Errorf("WithHeaders() Accept = %v, want 'application/json'", req.Headers["Accept"])
	}
}

func TestWithQueryParam(t *testing.T) {
	req := NewRequest("GET", "https://example.com", WithQueryParam("page", "1"))

	if req.QueryParams["page"] != "1" {
		t.Errorf("WithQueryParam() = %v, want '1'", req.QueryParams["page"])
	}
}

func TestWithQueryParams(t *testing.T) {
	params := map[string]string{
		"page":  "1",
		"limit": "10",
	}
	req := NewRequest("GET", "https://example.com", WithQueryParams(params))

	if req.QueryParams["page"] != "1" {
		t.Errorf("WithQueryParams() page = %v, want '1'", req.QueryParams["page"])
	}
	if req.QueryParams["limit"] != "10" {
		t.Errorf("WithQueryParams() limit = %v, want '10'", req.QueryParams["limit"])
	}
}

func TestWithContentType(t *testing.T) {
	req := NewRequest("POST", "https://example.com", WithContentType("application/json"))

	if req.ContentType != "application/json" {
		t.Errorf("WithContentType() = %v, want 'application/json'", req.ContentType)
	}
}

func TestWithTimeout(t *testing.T) {
	timeout := 30
	req := NewRequest("GET", "https://example.com", WithTimeout(timeout))

	if req.Timeout == nil {
		t.Error("WithTimeout() Timeout should not be nil")
	}
	if *req.Timeout != timeout {
		t.Errorf("WithTimeout() = %v, want %v", *req.Timeout, timeout)
	}
}

func TestNewRequest(t *testing.T) {
	req := NewRequest("GET", "https://example.com")

	if req.Method != "GET" {
		t.Errorf("NewRequest() Method = %v, want 'GET'", req.Method)
	}
	if req.URL != "https://example.com" {
		t.Errorf("NewRequest() URL = %v, want 'https://example.com'", req.URL)
	}
	if req.Headers == nil {
		t.Error("NewRequest() Headers should not be nil")
	}
	if req.QueryParams == nil {
		t.Error("NewRequest() QueryParams should not be nil")
	}
}

func TestNewRequest_WithOptions(t *testing.T) {
	req := NewRequest("POST", "https://example.com",
		WithHeader("Authorization", "Bearer token"),
		WithQueryParam("page", "1"),
		WithContentType("application/json"),
		WithTimeout(30),
	)

	if req.Method != "POST" {
		t.Errorf("NewRequest() Method = %v, want 'POST'", req.Method)
	}
	if req.Headers["Authorization"] != "Bearer token" {
		t.Errorf("NewRequest() Headers = %v, want 'Bearer token'", req.Headers["Authorization"])
	}
	if req.QueryParams["page"] != "1" {
		t.Errorf("NewRequest() QueryParams = %v, want '1'", req.QueryParams["page"])
	}
	if req.ContentType != "application/json" {
		t.Errorf("NewRequest() ContentType = %v, want 'application/json'", req.ContentType)
	}
	if req.Timeout == nil || *req.Timeout != 30 {
		t.Errorf("NewRequest() Timeout = %v, want 30", req.Timeout)
	}
}

func TestRequest_ToHTTPRequest(t *testing.T) {
	req := NewRequest("GET", "https://example.com", WithQueryParam("page", "1"))

	httpReq, err := req.ToHTTPRequest(context.Background())
	if err != nil {
		t.Fatalf("ToHTTPRequest() error = %v", err)
	}

	if httpReq.Method != "GET" {
		t.Errorf("ToHTTPRequest() Method = %v, want 'GET'", httpReq.Method)
	}
	if !strings.Contains(httpReq.URL.String(), "page=1") {
		t.Errorf("ToHTTPRequest() URL = %v, should contain 'page=1'", httpReq.URL.String())
	}
}

func TestRequest_ToHTTPRequest_WithMultipleQueryParams(t *testing.T) {
	req := NewRequest("GET", "https://example.com",
		WithQueryParam("page", "1"),
		WithQueryParam("limit", "10"),
	)

	httpReq, err := req.ToHTTPRequest(context.Background())
	if err != nil {
		t.Fatalf("ToHTTPRequest() error = %v", err)
	}

	url := httpReq.URL.String()
	if !strings.Contains(url, "page=1") {
		t.Errorf("ToHTTPRequest() URL = %v, should contain 'page=1'", url)
	}
	if !strings.Contains(url, "limit=10") {
		t.Errorf("ToHTTPRequest() URL = %v, should contain 'limit=10'", url)
	}
}

func TestRequest_ToHTTPRequest_WithHeaders(t *testing.T) {
	req := NewRequest("GET", "https://example.com",
		WithHeader("Authorization", "Bearer token"),
		WithContentType("application/json"),
	)

	httpReq, err := req.ToHTTPRequest(context.Background())
	if err != nil {
		t.Fatalf("ToHTTPRequest() error = %v", err)
	}

	if httpReq.Header.Get("Authorization") != "Bearer token" {
		t.Errorf("ToHTTPRequest() Authorization header = %v, want 'Bearer token'", httpReq.Header.Get("Authorization"))
	}
	if httpReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("ToHTTPRequest() Content-Type header = %v, want 'application/json'", httpReq.Header.Get("Content-Type"))
	}
}

func TestRequest_ToHTTPRequest_WithBody(t *testing.T) {
	body := strings.NewReader("test body")
	req := &Request{
		Method:      "POST",
		URL:         "https://example.com",
		Body:        body,
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}

	httpReq, err := req.ToHTTPRequest(context.Background())
	if err != nil {
		t.Fatalf("ToHTTPRequest() error = %v", err)
	}

	if httpReq.Body == nil {
		t.Error("ToHTTPRequest() Body should not be nil")
	}
}

func TestRequest_ToHTTPRequest_InvalidURL(t *testing.T) {
	req := &Request{
		Method:      "GET",
		URL:         "://invalid-url",
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}

	_, err := req.ToHTTPRequest(context.Background())
	if err == nil {
		t.Error("ToHTTPRequest() should return error for invalid URL")
	}
}

func TestWithHeader_NilHeaders(t *testing.T) {
	req := &Request{
		Method:      "GET",
		URL:         "https://example.com",
		Headers:     nil,
		QueryParams: make(map[string]string),
	}

	WithHeader("Authorization", "Bearer token")(req)

	if req.Headers == nil {
		t.Error("WithHeader() should initialize Headers if nil")
	}
	if req.Headers["Authorization"] != "Bearer token" {
		t.Errorf("WithHeader() = %v, want 'Bearer token'", req.Headers["Authorization"])
	}
}

func TestWithQueryParam_NilQueryParams(t *testing.T) {
	req := &Request{
		Method:      "GET",
		URL:         "https://example.com",
		Headers:     make(map[string]string),
		QueryParams: nil,
	}

	WithQueryParam("page", "1")(req)

	if req.QueryParams == nil {
		t.Error("WithQueryParam() should initialize QueryParams if nil")
	}
	if req.QueryParams["page"] != "1" {
		t.Errorf("WithQueryParam() = %v, want '1'", req.QueryParams["page"])
	}
}
