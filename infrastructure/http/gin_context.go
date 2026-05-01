package http

import (
	"context"
	"net/http"
	"reflect"

	"github.com/reitmas32/rkit/core/customctx"
)

// extractGinContextValues extracts values from a gin.Context's Keys map
// This is a helper function to extract values from gin.Context when available
func extractGinContextValues(ginCtx interface{}) map[string]any {
	values := make(map[string]any)

	ctxType := reflect.TypeOf(ginCtx)
	if ctxType == nil || ctxType.Kind() != reflect.Ptr {
		return values
	}

	ctxValue := reflect.ValueOf(ginCtx)
	if !ctxValue.IsValid() || ctxValue.IsNil() {
		return values
	}

	elem := ctxValue.Elem()
	if !elem.IsValid() {
		return values
	}

	// Look for the Keys field in gin.Context
	keysField := elem.FieldByName("Keys")
	if !keysField.IsValid() {
		return values
	}

	// Keys is a map[interface{}]interface{}
	if keysField.Kind() == reflect.Map {
		for _, key := range keysField.MapKeys() {
			val := keysField.MapIndex(key)
			if key.IsValid() && key.CanInterface() && val.IsValid() && val.CanInterface() {
				// Only store string keys
				if keyStr, ok := key.Interface().(string); ok {
					values[keyStr] = val.Interface()
				}
			}
		}
	}

	return values
}

// NewCustomContextFromGin creates a new CustomContext from a gin.Context.
// It extracts values from both gin.Context.Keys and the underlying context.Context.
// This is the recommended way to create a CustomContext from a gin.Context.
func NewCustomContextFromGin(ginCtx interface{}) *customctx.CustomContext {
	values := make(map[string]any)

	// Extract values from gin.Context.Keys
	ginValues := extractGinContextValues(ginCtx)
	for k, v := range ginValues {
		values[k] = v
	}

	// Try to get the underlying context.Context from gin.Context
	var parent context.Context
	ctxValue := reflect.ValueOf(ginCtx)
	if ctxValue.IsValid() && !ctxValue.IsNil() {
		elem := ctxValue.Elem()
		if elem.IsValid() {
			requestField := elem.FieldByName("Request")
			if requestField.IsValid() && requestField.CanInterface() {
				if req, ok := requestField.Interface().(interface{ Context() context.Context }); ok {
					parent = req.Context()
				}
			}
		}
	}

	// If we have a parent context, extract its values too
	if parent != nil {
		extractedValues := customctx.ExtractContextValues(parent)
		for k, v := range extractedValues {
			// Don't override gin.Context values
			if _, exists := values[k]; !exists {
				values[k] = v
			}
		}
	}

	// Use the parent context if available, otherwise use context.Background()
	if parent == nil {
		parent = context.Background()
	}

	return customctx.NewWithValues(parent, values)
}

// GetBaseURL constructs the base URL from a gin.Context.
// It determines the scheme (http/https) from:
// 1. X-Forwarded-Proto header (if present, typically set by reverse proxies)
// 2. Request.TLS (if the request is using TLS)
// 3. Defaults to "http" if neither is available
//
// The returned URL includes the scheme, host, and path (without query parameters).
// This is useful for constructing pagination links or other absolute URLs.
func GetBaseURL(ginCtx interface{}) string {
	// Use reflection to access gin.Context fields
	ctxValue := reflect.ValueOf(ginCtx)
	if !ctxValue.IsValid() || ctxValue.IsNil() {
		return ""
	}

	elem := ctxValue.Elem()
	if !elem.IsValid() {
		return ""
	}

	// Get the Request field
	requestField := elem.FieldByName("Request")
	if !requestField.IsValid() || !requestField.CanInterface() {
		return ""
	}

	// Try to get the request as *http.Request
	var req *http.Request
	if reqPtr, ok := requestField.Interface().(*http.Request); ok {
		req = reqPtr
	} else {
		return ""
	}

	// Determine scheme
	scheme := "http"
	if proto := req.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	} else if req.TLS != nil {
		scheme = "https"
	}

	// Get host and path
	host := req.Host
	path := req.URL.Path

	return scheme + "://" + host + path
}
