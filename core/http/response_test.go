package http

import (
	"io"
	"strings"
	"testing"
)

func TestResponse_ReadBody(t *testing.T) {
	testData := "test response body"
	body := io.NopCloser(strings.NewReader(testData))
	resp := &Response{
		StatusCode: 200,
		Body:       body,
	}

	data, err := resp.ReadBody()
	if err != nil {
		t.Fatalf("ReadBody() error = %v", err)
	}

	if string(data) != testData {
		t.Errorf("ReadBody() = %v, want %v", string(data), testData)
	}

	// Body should be closed after ReadBody
	_, err = resp.Body.Read(make([]byte, 1))
	if err == nil {
		t.Error("Body should be closed after ReadBody()")
	}
}

func TestResponse_ReadBodyString(t *testing.T) {
	testData := "test response body"
	body := io.NopCloser(strings.NewReader(testData))
	resp := &Response{
		StatusCode: 200,
		Body:       body,
	}

	str, err := resp.ReadBodyString()
	if err != nil {
		t.Fatalf("ReadBodyString() error = %v", err)
	}

	if str != testData {
		t.Errorf("ReadBodyString() = %v, want %v", str, testData)
	}
}

func TestResponse_Close(t *testing.T) {
	body := io.NopCloser(strings.NewReader("test"))
	resp := &Response{
		StatusCode: 200,
		Body:       body,
	}

	err := resp.Close()
	if err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Test closing nil body
	respNil := &Response{
		StatusCode: 200,
		Body:       nil,
	}

	err = respNil.Close()
	if err != nil {
		t.Fatalf("Close() with nil body error = %v", err)
	}
}

func TestResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"200 OK", 200, true},
		{"201 Created", 201, true},
		{"204 No Content", 204, true},
		{"299 Custom", 299, true},
		{"199 Not 2xx", 199, false},
		{"300 Redirect", 300, false},
		{"400 Bad Request", 400, false},
		{"500 Internal Server Error", 500, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &Response{StatusCode: tt.statusCode}
			if got := resp.IsSuccess(); got != tt.want {
				t.Errorf("IsSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_IsClientError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"400 Bad Request", 400, true},
		{"401 Unauthorized", 401, true},
		{"404 Not Found", 404, true},
		{"499 Custom", 499, true},
		{"399 Not 4xx", 399, false},
		{"500 Internal Server Error", 500, false},
		{"200 OK", 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &Response{StatusCode: tt.statusCode}
			if got := resp.IsClientError(); got != tt.want {
				t.Errorf("IsClientError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_IsServerError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"500 Internal Server Error", 500, true},
		{"502 Bad Gateway", 502, true},
		{"503 Service Unavailable", 503, true},
		{"599 Custom", 599, true},
		{"499 Not 5xx", 499, false},
		{"600 Not 5xx", 600, false},
		{"200 OK", 200, false},
		{"400 Bad Request", 400, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &Response{StatusCode: tt.statusCode}
			if got := resp.IsServerError(); got != tt.want {
				t.Errorf("IsServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_ReadBody_Empty(t *testing.T) {
	body := io.NopCloser(strings.NewReader(""))
	resp := &Response{
		StatusCode: 200,
		Body:       body,
	}

	data, err := resp.ReadBody()
	if err != nil {
		t.Fatalf("ReadBody() error = %v", err)
	}

	if len(data) != 0 {
		t.Errorf("ReadBody() = %v, want empty", data)
	}
}

func TestResponse_ReadBodyString_Empty(t *testing.T) {
	body := io.NopCloser(strings.NewReader(""))
	resp := &Response{
		StatusCode: 200,
		Body:       body,
	}

	str, err := resp.ReadBodyString()
	if err != nil {
		t.Fatalf("ReadBodyString() error = %v", err)
	}

	if str != "" {
		t.Errorf("ReadBodyString() = %v, want empty", str)
	}
}
