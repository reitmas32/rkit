// Package kerrors provides a custom error type with structured information
// including error codes, metadata, and error chaining capabilities.
// This package extends Go's standard error handling with additional context
// for better error tracking and debugging.
package kerrors

// KError represents a structured error with a message, error code, optional metadata,
// and an optional underlying cause error. It implements the error interface
// and supports error unwrapping for error chain inspection.
type KError struct {
	// Message is the human-readable error message describing what went wrong.
	Message string

	// Code is a numeric error code that can be used for programmatic error handling
	// and error categorization.
	Code int

	// Metadata is an optional map containing additional context about the error,
	// such as request IDs, timestamps, or other relevant debugging information.
	Metadata map[string]any

	// Cause is the underlying error that caused this error, if any.
	// This enables error wrapping and chain inspection using errors.Unwrap().
	Cause error
}

// NewKError creates a new KError with the given message, code, and metadata.
// The returned error will not have a cause error set.
func NewKError(
	message string,
	code int,
	metadata map[string]any,
) *KError {
	return &KError{Message: message, Code: code, Metadata: metadata}
}

// NewKErrorWithCause creates a new KError with the given message, code, metadata,
// and an underlying cause error. This is useful for wrapping existing errors
// with additional context while preserving the original error chain.
func NewKErrorWithCause(
	message string,
	code int,
	metadata map[string]any,
	cause error,
) *KError {
	return &KError{Message: message, Code: code, Metadata: metadata, Cause: cause}
}

// Error returns the error message, implementing the error interface.
// This method allows KError to be used anywhere a standard error is expected.
func (e *KError) Error() string {
	return e.Message
}

// Unwrap returns the underlying cause error, if any.
// This method enables error chain inspection using errors.Unwrap() and errors.Is().
// Returns nil if there is no underlying cause.
func (e *KError) Unwrap() error {
	return e.Cause
}

func (e *KError) WithCause(cause error) *KError {
	return &KError{Message: e.Message, Code: e.Code, Metadata: e.Metadata, Cause: cause}
}
