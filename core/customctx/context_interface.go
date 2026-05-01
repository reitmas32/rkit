package customctx

import (
	"context"
	"time"
)

// --- context.Context interface ---

// Deadline returns the time when work done on behalf of this context
// should be canceled. It delegates to the parent context's Deadline method.
func (c *CustomContext) Deadline() (time.Time, bool) { return c.parent.Deadline() }

// Done returns a channel that's closed when work done on behalf of this context
// should be canceled. It delegates to the parent context's Done method.
func (c *CustomContext) Done() <-chan struct{} { return c.parent.Done() }

// Err returns nil if Done is not closed. If Done is closed, Err returns
// the error that explains why. It delegates to the parent context's Err method.
func (c *CustomContext) Err() error { return c.parent.Err() }

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. It first checks the CustomContext's
// internal values map (only if key is a string), then delegates to the parent
// context's Value method.
//
// IMPORTANT: This should only be used for request-scoped technical metadata
// (e.g., request ID, trace ID, correlation ID). Never store business data here.
// Business data should be passed as function parameters.
//
// For string keys, it checks CustomContext's internal map first.
// For non-string keys or if not found in CustomContext, it delegates to parent.
func (c *CustomContext) Value(key interface{}) interface{} {
	// Only check CustomContext values if key is a string
	if strKey, ok := key.(string); ok && c.values != nil {
		c.mu.RLock()
		if val, exists := c.values[strKey]; exists {
			c.mu.RUnlock()
			return val
		}
		c.mu.RUnlock()
	}
	return c.parent.Value(key)
}

// --- Parent context access ---

// Context returns the parent context.Context that this CustomContext wraps.
// This allows access to the underlying context for operations that require
// the standard context interface.
func (c *CustomContext) Context() context.Context { return c.parent }
