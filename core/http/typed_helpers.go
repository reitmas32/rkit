package http

import (
	"time"

	"github.com/reitmas32/rkit/core/customctx"
)

// DoTyped performs an HTTP request and parses the response into a TypedResponse.
func DoTyped[T any](client Client, ctx *customctx.CustomContext, req *Request) (*TypedResponse[T], error) {
	requestTime := time.Now()
	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	return ParseResponse[T](resp, requestTime)
}

// GetTyped performs a GET request and parses the response into a TypedResponse.
// requestOpts are options for the HTTP request (e.g., WithHeader).
func GetTyped[T any](client Client, ctx *customctx.CustomContext, url string, requestOpts ...RequestOption) (*TypedResponse[T], error) {
	requestTime := time.Now()
	resp, err := client.Get(ctx, url, requestOpts...)
	if err != nil {
		return nil, err
	}
	return ParseResponse[T](resp, requestTime)
}

// DeleteTyped performs a DELETE request and parses the response into a TypedResponse.
func DeleteTyped[T any](client Client, ctx *customctx.CustomContext, url string, opts ...RequestOption) (*TypedResponse[T], error) {
	requestTime := time.Now()
	resp, err := client.Delete(ctx, url, opts...)
	if err != nil {
		return nil, err
	}
	return ParseResponse[T](resp, requestTime)
}

// HeadTyped performs a HEAD request and parses the response into a TypedResponse.
func HeadTyped[T any](client Client, ctx *customctx.CustomContext, url string, opts ...RequestOption) (*TypedResponse[T], error) {
	requestTime := time.Now()
	resp, err := client.Head(ctx, url, opts...)
	if err != nil {
		return nil, err
	}
	return ParseResponse[T](resp, requestTime)
}

// OptionsTyped performs an OPTIONS request and parses the response into a TypedResponse.
func OptionsTyped[T any](client Client, ctx *customctx.CustomContext, url string, opts ...RequestOption) (*TypedResponse[T], error) {
	requestTime := time.Now()
	resp, err := client.Options(ctx, url, opts...)
	if err != nil {
		return nil, err
	}
	return ParseResponse[T](resp, requestTime)
}
