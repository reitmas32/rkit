package kerrors_test

import (
	"errors"
	"fmt"

	"github.com/reitmas32/rkit/core/kerrors"
)

// ExampleNewKError builds a structured error carrying a numeric code and
// contextual metadata.
func ExampleNewKError() {
	err := kerrors.NewKError("user not found", 404, map[string]any{"user_id": "u-123"})
	fmt.Println(err.Code, err.Error())
	// Output: 404 user not found
}

// ExampleNewKErrorWithCause wraps an underlying error so that errors.Is and
// errors.Unwrap keep working across the chain.
func ExampleNewKErrorWithCause() {
	cause := errors.New("connection timeout")
	err := kerrors.NewKErrorWithCause("failed to save user", 500, nil, cause)
	fmt.Println(errors.Is(err, cause))
	// Output: true
}
