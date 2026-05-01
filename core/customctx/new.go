package customctx

import "context"

// New creates a new CustomContext that wraps the parent context.
// The returned CustomContext implements context.Context and delegates all
// standard context operations to the parent context while maintaining
// its own error collection.
// If the parent context is also a CustomContext, its values are copied to the new context.
// If the parent is a standard context, all values are extracted using reflection.
func New(parent context.Context) *CustomContext {
	values := make(map[string]any)

	// If parent is a CustomContext, copy its values
	if parentCC, ok := parent.(*CustomContext); ok {
		parentCC.mu.RLock()
		if parentCC.values != nil {
			for k, v := range parentCC.values {
				values[k] = v
			}
		}
		parentCC.mu.RUnlock()
	} else {
		// For standard contexts, extract all values using reflection
		extractedValues := extractContextValues(parent)
		for k, v := range extractedValues {
			values[k] = v
		}
	}

	return &CustomContext{parent: parent, values: values}
}

// NewWithValues creates a new CustomContext with the given parent context and pre-populated values.
// This is useful for infrastructure adapters that need to create a CustomContext with values
// extracted from framework-specific contexts (e.g., gin.Context).
func NewWithValues(parent context.Context, values map[string]any) *CustomContext {
	if values == nil {
		values = make(map[string]any)
	}
	return &CustomContext{parent: parent, values: values}
}

// NewCustomContext is an alias for New for compatibility.
// Deprecated: use New instead.
func NewCustomContext(parent context.Context) *CustomContext {
	return New(parent)
}
