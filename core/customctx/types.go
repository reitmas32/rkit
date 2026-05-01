// Package customctx provides a custom context implementation that extends
// Go's standard context.Context with the ability to accumulate structured errors
// throughout the execution flow. This is useful for collecting multiple errors
// that occur during processing without stopping execution, and for tracking
// where each error was registered in the call stack.
package customctx

import (
	"context"
	"sync"

	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/logger"
)

// WrapError associates a structured error with the location where it was registered.
// It contains both the KError and the call site information (function name and line number)
// where the error was added to the context.
type WrapError struct {
	// Error is the structured error that was registered.
	Error *kerrors.KError `json:"error"`

	// CallIn is the location where the error was registered, in the format "functionName:lineNumber".
	CallIn string `json:"call_in"`
}

// CustomContext is a context.Context implementation that accumulates structured errors.
// It wraps a parent context.Context and maintains all the standard context behavior
// while adding the ability to collect multiple errors during execution.
// All methods are safe for concurrent use.
type CustomContext struct {
	parent context.Context
	mu     sync.RWMutex
	errors []WrapError
	values map[string]any // Technical metadata only, not business data

	// Logger is the logger that was used to register the error.
	logger logger.ILogger
}
