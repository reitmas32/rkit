package http

import (
	"bytes"
	"encoding/json"
	"io"
)

// RequestBody is an interface for types that can be serialized to a request body.
// Types implementing this interface can be passed directly to Post, Put, Patch methods.
type RequestBody interface {
	// ToReader converts the request body to an io.Reader.
	// This allows the HTTP client to read the body for the request.
	ToReader() (io.Reader, error)
}

// ResponseBody is an interface for types that can be deserialized from a response body.
// This is mainly for documentation purposes, as any type can be used with TypedResponse.
type ResponseBody interface {
	// Any type can be a ResponseBody
	// The actual deserialization is handled by ParseResponse using JSON unmarshaling
}

// jsonRequestBody is a helper type that wraps any JSON-serializable value.
type jsonRequestBody struct {
	value any
}

// ToReader implements RequestBody by marshaling the value to JSON.
func (j *jsonRequestBody) ToReader() (io.Reader, error) {
	data, err := json.Marshal(j.value)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

// NewJSONRequestBody creates a RequestBody from any JSON-serializable value.
func NewJSONRequestBody(value any) RequestBody {
	return &jsonRequestBody{value: value}
}

// readerRequestBody is a helper type that wraps an io.Reader.
type readerRequestBody struct {
	reader io.Reader
}

// ToReader implements RequestBody by returning the wrapped reader.
func (r *readerRequestBody) ToReader() (io.Reader, error) {
	return r.reader, nil
}

// NewReaderRequestBody creates a RequestBody from an io.Reader.
func NewReaderRequestBody(reader io.Reader) RequestBody {
	return &readerRequestBody{reader: reader}
}
