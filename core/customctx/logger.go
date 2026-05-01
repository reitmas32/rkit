package customctx

import "github.com/reitmas32/rkit/core/logger"

// --- Logger management ---

// WithLogger returns a new CustomContext with the given logger.
// The returned context is immutable - the original context is not modified.
func (c *CustomContext) WithLogger(logger logger.ILogger) *CustomContext {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Copy errors slice
	newErrors := make([]WrapError, len(c.errors))
	copy(newErrors, c.errors)

	// Copy values map
	var newValues map[string]any
	if c.values != nil {
		newValues = make(map[string]any)
		for k, v := range c.values {
			newValues[k] = v
		}
	}

	return &CustomContext{
		parent: c.parent,
		errors: newErrors,
		values: newValues,
		logger: logger,
	}
}

// Logger returns the logger associated with this context.
// If no logger is set, it returns a default simple logger.
func (c *CustomContext) Logger() logger.ILogger {
	if c.logger == nil {
		return logger.NewSimpleLogger(logger.LevelDebug.String())
	}
	return c.logger
}
