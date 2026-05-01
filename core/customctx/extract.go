package customctx

import (
	"context"
	"reflect"
)

// ExtractContextValues uses reflection to extract all key-value pairs from a context chain.
// This is exported so it can be used by infrastructure adapters (e.g., Gin).
func ExtractContextValues(ctx context.Context) map[string]any {
	values := make(map[string]any)
	if ctx == nil {
		return values
	}

	// Recursively traverse the context chain
	extractFromContext(ctx, values, make(map[context.Context]bool))

	return values
}

// extractContextValues is an internal alias for ExtractContextValues for backward compatibility
func extractContextValues(ctx context.Context) map[string]any {
	return ExtractContextValues(ctx)
}

// extractFromContext recursively extracts values from a context and its parents
func extractFromContext(ctx context.Context, values map[string]any, visited map[context.Context]bool) {
	if ctx == nil {
		return
	}

	// Avoid infinite loops
	if visited[ctx] {
		return
	}
	visited[ctx] = true

	ctxType := reflect.TypeOf(ctx)
	if ctxType == nil {
		return
	}

	ctxValue := reflect.ValueOf(ctx)
	if !ctxValue.IsValid() {
		return
	}

	typeName := ctxType.String()

	// Skip empty contexts (they don't have values)
	if typeName == "*context.emptyCtx" || typeName == "context.emptyCtx" {
		return
	}

	// Handle valueCtx - standard Go context with a value
	if ctxType.Kind() == reflect.Ptr {
		elemType := ctxType.Elem()

		// Check if it's a valueCtx by looking for key and val fields
		if elemType.NumField() >= 2 {
			var keyField, valField, parentField *reflect.StructField

			for i := 0; i < elemType.NumField(); i++ {
				field := elemType.Field(i)
				switch field.Name {
				case "key":
					keyField = &field
				case "val":
					valField = &field
				case "Context":
					parentField = &field
				}
			}

			// If we found key and val, it's likely a valueCtx
			if keyField != nil && valField != nil {
				elem := ctxValue.Elem()
				if elem.IsValid() {
					keyValue := elem.FieldByName("key")
					valValue := elem.FieldByName("val")

					if keyValue.IsValid() && keyValue.CanInterface() &&
						valValue.IsValid() && valValue.CanInterface() {
						key := keyValue.Interface()
						val := valValue.Interface()

						// Only store string keys
						if keyStr, ok := key.(string); ok {
							values[keyStr] = val
						}
					}

					// Continue with parent context
					if parentField != nil {
						parentValue := elem.FieldByName("Context")
						if parentValue.IsValid() && parentValue.CanInterface() {
							if parentCtx, ok := parentValue.Interface().(context.Context); ok {
								extractFromContext(parentCtx, values, visited)
							}
						}
					}
				}
				return
			}
		}
	}

	// For other context types, try to find a parent context field
	if ctxType.Kind() == reflect.Ptr {
		elemType := ctxType.Elem()
		elem := ctxValue.Elem()
		if elem.IsValid() {
			for i := 0; i < elemType.NumField(); i++ {
				field := elemType.Field(i)
				if field.Name == "Context" || field.Name == "parent" {
					fieldValue := elem.Field(i)
					if fieldValue.IsValid() && fieldValue.CanInterface() {
						if parentCtx, ok := fieldValue.Interface().(context.Context); ok {
							extractFromContext(parentCtx, values, visited)
							return
						}
					}
				}
			}
		}
	}
}
