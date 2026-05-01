// Package result provides a Result type that represents either a value or an error.
// This is similar to Rust's Result type or functional programming's Either type,
// allowing for explicit error handling without relying on Go's traditional error return pattern.
package result

import (
	"reflect"

	"github.com/reitmas32/rkit/core/kerrors"
)

// Result is a generic type that represents either a successful value of type T or an error.
// It provides a functional approach to error handling, allowing you to chain operations
// and handle errors in a more declarative way.
type Result[T any] struct {
	value  T
	_error *kerrors.KError
}

// NewResult creates a new Result with both a value and an error.
// If err is nil, the Result is considered successful (Ok).
// If err is not nil, the Result is considered failed (Err).
func NewResult[T any](value T, err *kerrors.KError) Result[T] {
	return Result[T]{value: value, _error: err}
}

// NewOkResult creates a new successful Result with the given value.
// The Result will have no error and IsOk() will return true.
func NewOkResult[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// NewErrResult creates a new failed Result with the given error.
// The Result will have a zero value for T and IsOk() will return false.
func NewErrResult[T any](err *kerrors.KError) Result[T] {
	return Result[T]{_error: err}
}

// Ok creates a new successful Result with the given value.
// This is a convenience function equivalent to NewOkResult.
// The Result will have no error and IsOk() will return true.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

func Empty[T any]() Result[T] {
	return Result[T]{}
}

// Err creates a new failed Result with the given error.
// This is a convenience function equivalent to NewErrResult.
// The Result will have a zero value for T and IsOk() will return false.
func Err[T any](err *kerrors.KError) Result[T] {
	return Result[T]{_error: err}
}

// Value returns the value contained in the Result.
// Note: This method returns the value regardless of whether the Result is Ok or Err.
// Use IsOk() to check if the Result is successful before accessing the value.
func (r Result[T]) Value() T {
	return r.value
}

// Error returns the error contained in the Result.
// If the Result is successful (Ok), this will return nil.
// If the Result is failed (Err), this will return the error that was set.
func (r Result[T]) Error() error {
	if r._error == nil {
		return nil
	}
	return r._error
}

// IsOk returns true if the Result represents a successful value (no error).
// Returns false if the Result represents an error.
func (r Result[T]) IsOk() bool {
	return r._error == nil
}

// IsEmpty returns true if the Result is empty (no value and no error).
// Returns false if the Result has a value or an error.
func (r Result[T]) IsEmpty() bool {
	return reflect.ValueOf(r.value).IsZero() && r._error == nil
}

func (r Result[T]) ToKError() *kerrors.KError {
	if r.IsOk() {
		return nil
	}
	return r._error
}
