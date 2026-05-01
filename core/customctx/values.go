package customctx

// --- Value storage (for technical metadata only) ---

// WithValue returns a new CustomContext with the given key-value pair added.
// The returned context is immutable - the original context is not modified.
//
// IMPORTANT: This should only be used for request-scoped technical metadata:
//   - Request IDs, trace IDs, correlation IDs
//   - User IDs (for technical purposes like logging/tracing, not business logic)
//   - Service names, operation names
//   - Technical flags or configuration
//
// DO NOT use this for business data:
//   - Domain entities (users, orders, products, etc.)
//   - Business state or context
//   - Application-specific domain values
//
// Business data should be passed as explicit function parameters, not stored in context.
//
// The key must be a string. Use descriptive key names to avoid collisions, for example:
//
//	ctx := customctx.New(context.Background())
//	ctx = ctx.WithValue("request_id", "req-123")
//	ctx = ctx.WithValue("trace_id", "trace-456")
//	ctx = ctx.WithValue("user_id", "user-789")
//
// This method is safe for concurrent use and returns a new context that shares
// the parent context, errors, and logger, but has its own values map.
func (c *CustomContext) WithValue(key string, val any) *CustomContext {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create new context with copied values
	newValues := make(map[string]any)
	if c.values != nil {
		for k, v := range c.values {
			newValues[k] = v
		}
	}
	newValues[key] = val

	// Copy errors slice
	newErrors := make([]WrapError, len(c.errors))
	copy(newErrors, c.errors)

	return &CustomContext{
		parent: c.parent,
		errors: newErrors,
		values: newValues,
		logger: c.logger,
	}
}

// GetValue returns the value associated with the given string key in this CustomContext only,
// without checking the parent context. Returns nil if the key is not found.
//
// This is a convenience method that is equivalent to checking if the key exists
// in the CustomContext's internal values map. To check both CustomContext and
// parent context, use Value() instead.
//
// IMPORTANT: See WithValue() for guidelines on what types of values should be stored.
//
// This method is thread-safe.
func (c *CustomContext) GetValue(key string) any {
	if c.values == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.values[key]
}

// HasValue returns true if the given string key exists in this CustomContext's values map,
// without checking the parent context.
//
// This is a convenience method for checking key existence. To check both CustomContext
// and parent context, use Value() and check for nil instead.
//
// This method is thread-safe.
func (c *CustomContext) HasValue(key string) bool {
	if c.values == nil {
		return false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.values[key]
	return exists
}
