package http

import "github.com/reitmas32/rkit/core/kerrors"

var (
	// ErrInvalidRequest is returned when a request is invalid
	ErrInvalidRequest = kerrors.NewKError("invalid HTTP request", 400, nil)

	// ErrRequestTimeout is returned when a request times out
	ErrRequestTimeout = kerrors.NewKError("HTTP request timeout", 408, nil)

	// ErrRequestFailed is returned when a request fails
	ErrRequestFailed = kerrors.NewKError("HTTP request failed", 500, nil)

	// ErrInvalidResponse is returned when a response is invalid
	ErrInvalidResponse = kerrors.NewKError("invalid HTTP response", 500, nil)
)
