package http

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"
)

func TestTypedResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		response *TypedResponse[string]
		expected bool
	}{
		{"200 OK", &TypedResponse[string]{StatusCode: 200}, true},
		{"201 Created", &TypedResponse[string]{StatusCode: 201}, true},
		{"299 Custom", &TypedResponse[string]{StatusCode: 299}, true},
		{"400 Bad Request", &TypedResponse[string]{StatusCode: 400}, false},
		{"500 Internal Server Error", &TypedResponse[string]{StatusCode: 500}, false},
		{"with expected status code match", &TypedResponse[string]{StatusCode: 201, ExpectedStatusCode: intPtr(201)}, true},
		{"with expected status code mismatch", &TypedResponse[string]{StatusCode: 200, ExpectedStatusCode: intPtr(201)}, false},
		{"with status code range in range", &TypedResponse[string]{StatusCode: 201, SuccessStatusCodeRange: &StatusCodeRange{Min: 200, Max: 299}}, true},
		{"with status code range out of range", &TypedResponse[string]{StatusCode: 300, SuccessStatusCodeRange: &StatusCodeRange{Min: 200, Max: 299}}, false},
		{"expected status code takes precedence", &TypedResponse[string]{StatusCode: 201, ExpectedStatusCode: intPtr(200), SuccessStatusCodeRange: &StatusCodeRange{Min: 200, Max: 299}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.IsSuccess(); got != tt.expected {
				t.Errorf("IsSuccess() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTypedResponse_IsClientError(t *testing.T) {
	tests := []struct {
		name     string
		response *TypedResponse[string]
		expected bool
	}{
		{"400 Bad Request", &TypedResponse[string]{StatusCode: 400}, true},
		{"401 Unauthorized", &TypedResponse[string]{StatusCode: 401}, true},
		{"404 Not Found", &TypedResponse[string]{StatusCode: 404}, true},
		{"499 Custom", &TypedResponse[string]{StatusCode: 499}, true},
		{"200 OK", &TypedResponse[string]{StatusCode: 200}, false},
		{"500 Internal Server Error", &TypedResponse[string]{StatusCode: 500}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.IsClientError(); got != tt.expected {
				t.Errorf("IsClientError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTypedResponse_IsServerError(t *testing.T) {
	tests := []struct {
		name     string
		response *TypedResponse[string]
		expected bool
	}{
		{"500 Internal Server Error", &TypedResponse[string]{StatusCode: 500}, true},
		{"502 Bad Gateway", &TypedResponse[string]{StatusCode: 502}, true},
		{"503 Service Unavailable", &TypedResponse[string]{StatusCode: 503}, true},
		{"599 Custom", &TypedResponse[string]{StatusCode: 599}, true},
		{"200 OK", &TypedResponse[string]{StatusCode: 200}, false},
		{"400 Bad Request", &TypedResponse[string]{StatusCode: 400}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.IsServerError(); got != tt.expected {
				t.Errorf("IsServerError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseResponse(t *testing.T) {
	testData := map[string]string{"key": "value"}
	jsonData, _ := json.Marshal(testData)
	body := io.NopCloser(strings.NewReader(string(jsonData)))

	resp := &Response{
		StatusCode: 200,
		Status:     "200 OK",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       body,
	}

	requestTime := time.Now()
	typedResp, err := ParseResponse[map[string]string](resp, requestTime)
	if err != nil {
		t.Fatalf("ParseResponse() error = %v", err)
	}

	if typedResp.StatusCode != 200 {
		t.Errorf("ParseResponse() StatusCode = %v, want 200", typedResp.StatusCode)
	}
	if typedResp.Body["key"] != "value" {
		t.Errorf("ParseResponse() Body = %v, want map with key='value'", typedResp.Body)
	}
}

func TestParseResponse_WithBytes(t *testing.T) {
	testData := []byte("test data")
	body := io.NopCloser(strings.NewReader(string(testData)))

	resp := &Response{
		StatusCode: 200,
		Status:     "200 OK",
		Headers:    map[string]string{"Content-Type": "application/octet-stream"},
		Body:       body,
	}

	requestTime := time.Now()
	typedResp, err := ParseResponse[[]byte](resp, requestTime)
	if err != nil {
		t.Fatalf("ParseResponse() error = %v", err)
	}

	if string(typedResp.Body) != string(testData) {
		t.Errorf("ParseResponse() Body = %v, want %v", typedResp.Body, testData)
	}
}

func TestParseResponse_EmptyBody(t *testing.T) {
	body := io.NopCloser(strings.NewReader(""))

	resp := &Response{
		StatusCode: 204,
		Status:     "204 No Content",
		Headers:    map[string]string{},
		Body:       body,
	}

	requestTime := time.Now()
	typedResp, err := ParseResponse[map[string]interface{}](resp, requestTime)
	if err != nil {
		t.Fatalf("ParseResponse() error = %v", err)
	}

	if typedResp.StatusCode != 204 {
		t.Errorf("ParseResponse() StatusCode = %v, want 204", typedResp.StatusCode)
	}
}

func TestParseResponse_InvalidJSON(t *testing.T) {
	body := io.NopCloser(strings.NewReader("invalid json"))

	resp := &Response{
		StatusCode: 200,
		Status:     "200 OK",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       body,
	}

	requestTime := time.Now()
	_, err := ParseResponse[map[string]interface{}](resp, requestTime)
	if err == nil {
		t.Error("ParseResponse() should return error for invalid JSON")
	}
}

func TestParseResponse_ContentTypeCaseInsensitive(t *testing.T) {
	testData := map[string]string{"key": "value"}
	jsonData, _ := json.Marshal(testData)
	body := io.NopCloser(strings.NewReader(string(jsonData)))

	resp := &Response{
		StatusCode: 200,
		Status:     "200 OK",
		Headers:    map[string]string{"content-type": "application/json"},
		Body:       body,
	}

	requestTime := time.Now()
	typedResp, err := ParseResponse[map[string]string](resp, requestTime)
	if err != nil {
		t.Fatalf("ParseResponse() error = %v", err)
	}

	if typedResp.ContentType != "application/json" {
		t.Errorf("ParseResponse() ContentType = %v, want 'application/json'", typedResp.ContentType)
	}
}

func TestParseResponseFromReader(t *testing.T) {
	testData := map[string]string{"key": "value"}
	jsonData, _ := json.Marshal(testData)
	reader := io.NopCloser(strings.NewReader(string(jsonData)))

	headers := map[string]string{"Content-Type": "application/json"}
	requestTime := time.Now()

	typedResp, err := ParseResponseFromReader[map[string]string](reader, 200, "200 OK", headers, requestTime)
	if err != nil {
		t.Fatalf("ParseResponseFromReader() error = %v", err)
	}

	if typedResp.StatusCode != 200 {
		t.Errorf("ParseResponseFromReader() StatusCode = %v, want 200", typedResp.StatusCode)
	}
	if typedResp.Body["key"] != "value" {
		t.Errorf("ParseResponseFromReader() Body = %v, want map with key='value'", typedResp.Body)
	}
}

func TestParseResponseFromReader_WithBytes(t *testing.T) {
	testData := []byte("test data")
	reader := io.NopCloser(strings.NewReader(string(testData)))

	headers := map[string]string{"Content-Type": "application/octet-stream"}
	requestTime := time.Now()

	typedResp, err := ParseResponseFromReader[[]byte](reader, 200, "200 OK", headers, requestTime)
	if err != nil {
		t.Fatalf("ParseResponseFromReader() error = %v", err)
	}

	if string(typedResp.Body) != string(testData) {
		t.Errorf("ParseResponseFromReader() Body = %v, want %v", typedResp.Body, testData)
	}
}

func TestParseResponse_Duration(t *testing.T) {
	testData := map[string]string{"key": "value"}
	jsonData, _ := json.Marshal(testData)
	body := io.NopCloser(strings.NewReader(string(jsonData)))

	resp := &Response{
		StatusCode: 200,
		Status:     "200 OK",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       body,
	}

	requestTime := time.Now().Add(-100 * time.Millisecond)
	typedResp, err := ParseResponse[map[string]string](resp, requestTime)
	if err != nil {
		t.Fatalf("ParseResponse() error = %v", err)
	}

	if typedResp.Duration < 100*time.Millisecond {
		t.Errorf("ParseResponse() Duration = %v, should be at least 100ms", typedResp.Duration)
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}
