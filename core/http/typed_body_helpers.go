package http

import (
	"time"

	"github.com/reitmas32/rkit/core/customctx"
)

// PostTyped performs a POST request and parses the response into a TypedResponse.
// TRequest is the type of the request body (must implement RequestBody or be JSON-serializable).
// TResponse is the type of the response body.
// If body implements RequestBody, it will be used directly. Otherwise, it will be JSON-marshaled.
// If body is nil, no body will be sent.
func PostTyped[TRequest any, TResponse any](client Client, ctx *customctx.CustomContext, url string, body TRequest, opts ...RequestOption) (*TypedResponse[TResponse], error) {
	requestTime := time.Now()

	bodyReader, err := convertBodyToReader(body)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(ctx, url, bodyReader, opts...)
	if err != nil {
		return nil, err
	}
	return ParseResponse[TResponse](resp, requestTime)
}

// PutTyped performs a PUT request and parses the response into a TypedResponse.
// TRequest is the type of the request body (must implement RequestBody or be JSON-serializable).
// TResponse is the type of the response body.
// If body implements RequestBody, it will be used directly. Otherwise, it will be JSON-marshaled.
// If body is nil, no body will be sent.
func PutTyped[TRequest any, TResponse any](client Client, ctx *customctx.CustomContext, url string, body TRequest, opts ...RequestOption) (*TypedResponse[TResponse], error) {
	requestTime := time.Now()

	bodyReader, err := convertBodyToReader(body)
	if err != nil {
		return nil, err
	}

	resp, err := client.Put(ctx, url, bodyReader, opts...)
	if err != nil {
		return nil, err
	}
	return ParseResponse[TResponse](resp, requestTime)
}

// PatchTyped performs a PATCH request and parses the response into a TypedResponse.
// TRequest is the type of the request body (must implement RequestBody or be JSON-serializable).
// TResponse is the type of the response body.
// If body implements RequestBody, it will be used directly. Otherwise, it will be JSON-marshaled.
// If body is nil, no body will be sent.
func PatchTyped[TRequest any, TResponse any](client Client, ctx *customctx.CustomContext, url string, body TRequest, opts ...RequestOption) (*TypedResponse[TResponse], error) {
	requestTime := time.Now()

	bodyReader, err := convertBodyToReader(body)
	if err != nil {
		return nil, err
	}

	resp, err := client.Patch(ctx, url, bodyReader, opts...)
	if err != nil {
		return nil, err
	}
	return ParseResponse[TResponse](resp, requestTime)
}
