package customctx

import (
	"fmt"
	"runtime"

	"github.com/reitmas32/rkit/core/kerrors"
)

// --- Error management ---

// Errors returns all accumulated errors with their associated call sites.
// The returned slice contains all errors that have been registered via AddError,
// ordered by the sequence in which they were added.
// Note: The returned slice is the internal slice, so modifications to it may affect
// the context's internal state. Consider making a copy if you need to modify the slice.
func (c *CustomContext) Errors() []WrapError {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.errors
}

// HasErrors returns true if there is at least one registered error in the context.
// This method is safe for concurrent use and can be used to check if any errors
// were accumulated before accessing error details.
func (c *CustomContext) HasErrors() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.errors) > 0
}

// FirstError returns the first registered error (the earliest error in the collection).
// This method panics if there are no errors registered.
// Use HasErrors() to check for errors before calling this method.
func (c *CustomContext) FirstError() WrapError {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.errors) == 0 {
		panic("CustomContext: FirstError called on context with no errors")
	}
	return c.errors[0]
}

// LastError returns the last registered error (the most recent error in the collection).
// This method panics if there are no errors registered.
// Use HasErrors() to check for errors before calling this method.
func (c *CustomContext) LastError() WrapError {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.errors) == 0 {
		panic("CustomContext: LastError called on context with no errors")
	}
	return c.errors[len(c.errors)-1]
}

// AddError registers a structured error with the context, automatically capturing
// the caller information (function name and line number) where this method was called.
// If err is nil, this method does nothing and returns nil.
// The error is added to the internal collection and can be retrieved later using
// Errors(), FirstError(), or LastError().
// This method is safe for concurrent use.
func (c *CustomContext) AddError(err *kerrors.KError) *kerrors.KError {
	if err == nil {
		return nil
	}

	pc, _, line, ok := runtime.Caller(1)
	caller := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		caller = fmt.Sprintf("%s:%d", fn.Name(), line)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.errors = append(c.errors, WrapError{Error: err, CallIn: caller})
	return err
}

// NewError is an alias for AddError for compatibility.
// Deprecated: use AddError instead.
func (c *CustomContext) NewError(err *kerrors.KError) *kerrors.KError {
	return c.AddError(err)
}

// Clear removes all accumulated errors from the context, resetting the error collection.
// After calling Clear(), HasErrors() will return false and Errors() will return an empty slice.
// This method is safe for concurrent use.
func (c *CustomContext) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errors = nil
}
