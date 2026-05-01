package http

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestJSONRequestBody_ToReader(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"string value", "test", false},
		{"int value", 42, false},
		{"map value", map[string]string{"key": "value"}, false},
		{"struct value", struct{ Name string }{"test"}, false},
		{"nil value", nil, false},
		{"slice value", []int{1, 2, 3}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := NewJSONRequestBody(tt.value)
			reader, err := body.ToReader()
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonRequestBody.ToReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && reader == nil {
				t.Error("jsonRequestBody.ToReader() returned nil reader")
			}
			if reader != nil {
				data, _ := io.ReadAll(reader)
				if len(data) == 0 && tt.value != nil {
					t.Error("jsonRequestBody.ToReader() returned empty data")
				}
			}
		})
	}
}

func TestJSONRequestBody_ToReader_Readable(t *testing.T) {
	testValue := map[string]string{"key": "value", "test": "data"}
	body := NewJSONRequestBody(testValue)
	reader, err := body.ToReader()
	if err != nil {
		t.Fatalf("ToReader() error = %v", err)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if len(data) == 0 {
		t.Error("ToReader() returned empty data")
	}

	// Verify it's valid JSON
	if !strings.Contains(string(data), "key") {
		t.Error("ToReader() data doesn't contain expected content")
	}
}

func TestReaderRequestBody_ToReader(t *testing.T) {
	tests := []struct {
		name    string
		reader  io.Reader
		wantErr bool
	}{
		{"bytes reader", bytes.NewReader([]byte("test")), false},
		{"strings reader", strings.NewReader("test"), false},
		{"nil reader", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := NewReaderRequestBody(tt.reader)
			reader, err := body.ToReader()
			if (err != nil) != tt.wantErr {
				t.Errorf("readerRequestBody.ToReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && body != nil {
				if tt.reader == nil && reader != nil {
					t.Error("readerRequestBody.ToReader() should return nil for nil reader")
				}
				if tt.reader != nil && reader == nil {
					t.Error("readerRequestBody.ToReader() returned nil reader for non-nil input")
				}
			}
		})
	}
}

func TestReaderRequestBody_ToReader_Readable(t *testing.T) {
	testData := []byte("test data")
	originalReader := bytes.NewReader(testData)
	body := NewReaderRequestBody(originalReader)
	reader, err := body.ToReader()
	if err != nil {
		t.Fatalf("ToReader() error = %v", err)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if !bytes.Equal(data, testData) {
		t.Errorf("ToReader() data = %v, want %v", data, testData)
	}
}

func TestNewJSONRequestBody(t *testing.T) {
	value := map[string]int{"test": 42}
	body := NewJSONRequestBody(value)

	if body == nil {
		t.Fatal("NewJSONRequestBody() returned nil")
	}

	reader, err := body.ToReader()
	if err != nil {
		t.Fatalf("ToReader() error = %v", err)
	}

	data, _ := io.ReadAll(reader)
	if len(data) == 0 {
		t.Error("NewJSONRequestBody() created body with no data")
	}
}

func TestNewReaderRequestBody(t *testing.T) {
	testReader := strings.NewReader("test")
	body := NewReaderRequestBody(testReader)

	if body == nil {
		t.Fatal("NewReaderRequestBody() returned nil")
	}

	reader, err := body.ToReader()
	if err != nil {
		t.Fatalf("ToReader() error = %v", err)
	}

	if reader != testReader {
		t.Error("NewReaderRequestBody() should return the same reader")
	}
}

func TestJSONRequestBody_MultipleReads(t *testing.T) {
	value := map[string]string{"key": "value"}
	body := NewJSONRequestBody(value)

	// First read
	reader1, _ := body.ToReader()
	data1, _ := io.ReadAll(reader1)

	// Second read should work (new reader)
	reader2, _ := body.ToReader()
	data2, _ := io.ReadAll(reader2)

	if !bytes.Equal(data1, data2) {
		t.Error("Multiple reads should return same data")
	}
}
