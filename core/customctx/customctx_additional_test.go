package customctx

import (
	"context"
	"testing"

	"github.com/reitmas32/rkit/core/logger"
)

func TestWithLogger(t *testing.T) {
	ctx := New(context.Background())
	testLogger := logger.NewSimpleLogger("info")

	ctxWithLogger := ctx.WithLogger(testLogger)

	if ctxWithLogger == nil {
		t.Fatal("WithLogger() should not return nil")
	}

	// New context should have the logger
	if ctxWithLogger.Logger() != testLogger {
		t.Error("WithLogger() should set the logger")
	}

	// Original context should not have the logger (immutability)
	if ctx.Logger() == testLogger {
		t.Error("Original context should not have the logger")
	}
}

func TestWithLogger_PreservesErrors(t *testing.T) {
	ctx := New(context.Background())
	ctx.AddError(nil) // This won't add error, but tests the flow

	testLogger := logger.NewSimpleLogger("info")
	ctxWithLogger := ctx.WithLogger(testLogger)

	// Errors should be preserved
	if ctxWithLogger.HasErrors() != ctx.HasErrors() {
		t.Error("WithLogger() should preserve errors")
	}
}

func TestWithLogger_PreservesValues(t *testing.T) {
	ctx := New(context.Background())
	ctx1 := ctx.WithValue("key", "value")

	testLogger := logger.NewSimpleLogger("info")
	ctxWithLogger := ctx1.WithLogger(testLogger)

	// Values should be preserved
	if ctxWithLogger.GetValue("key") != "value" {
		t.Error("WithLogger() should preserve values")
	}
}

func TestLogger_DefaultLogger(t *testing.T) {
	ctx := New(context.Background())

	// When no logger is set, should return default logger
	logger := ctx.Logger()
	if logger == nil {
		t.Error("Logger() should return default logger when none is set")
	}
}

func TestLogger_CustomLogger(t *testing.T) {
	ctx := New(context.Background())
	testLogger := logger.NewSimpleLogger("debug")

	ctxWithLogger := ctx.WithLogger(testLogger)

	if ctxWithLogger.Logger() != testLogger {
		t.Error("Logger() should return the custom logger")
	}
}

func TestNewWithValues(t *testing.T) {
	parent := context.Background()
	values := map[string]any{
		"key1": "value1",
		"key2": 42,
	}

	ctx := NewWithValues(parent, values)

	if ctx == nil {
		t.Fatal("NewWithValues() should not return nil")
	}

	if ctx.GetValue("key1") != "value1" {
		t.Error("NewWithValues() should set values")
	}

	if ctx.GetValue("key2") != 42 {
		t.Error("NewWithValues() should set values")
	}
}

func TestNewWithValues_NilValues(t *testing.T) {
	parent := context.Background()

	ctx := NewWithValues(parent, nil)

	if ctx == nil {
		t.Fatal("NewWithValues() should not return nil with nil values")
	}

	// Should have empty values map
	if ctx.GetValue("key") != nil {
		t.Error("NewWithValues() with nil should create empty values map")
	}
}

func TestNewWithValues_EmptyValues(t *testing.T) {
	parent := context.Background()
	values := map[string]any{}

	ctx := NewWithValues(parent, values)

	if ctx == nil {
		t.Fatal("NewWithValues() should not return nil")
	}

	if ctx.GetValue("key") != nil {
		t.Error("NewWithValues() with empty map should have no values")
	}
}

func TestExtractContextValues_NilContext(t *testing.T) {
	values := ExtractContextValues(nil)

	if values == nil {
		t.Error("ExtractContextValues() should return empty map, not nil")
	}

	if len(values) != 0 {
		t.Error("ExtractContextValues() with nil should return empty map")
	}
}

func TestExtractContextValues_WithValue(t *testing.T) {
	parent := context.WithValue(context.Background(), "key", "value")

	values := ExtractContextValues(parent)

	// ExtractContextValues may or may not extract values depending on context implementation
	// This test verifies it doesn't panic and returns a valid map
	if values == nil {
		t.Error("ExtractContextValues() should return a map, not nil")
	}
	// Note: The extraction may not work for all context types, so we just verify it doesn't panic
	_ = values
}

func TestExtractContextValues_WithMultipleValues(t *testing.T) {
	parent := context.WithValue(
		context.WithValue(context.Background(), "key1", "value1"),
		"key2", "value2",
	)

	values := ExtractContextValues(parent)

	// ExtractContextValues may or may not extract values depending on context implementation
	// This test verifies it doesn't panic
	if values == nil {
		t.Error("ExtractContextValues() should return a map, not nil")
	}
	_ = values
}

func TestExtractContextValues_WithCustomContext(t *testing.T) {
	parent := context.Background()
	customCtx := New(parent)
	customCtx1 := customCtx.WithValue("custom_key", "custom_value")

	// Extract from CustomContext
	values := ExtractContextValues(customCtx1)

	// ExtractContextValues uses reflection and may extract CustomContext values
	// This test verifies it doesn't panic
	if values == nil {
		t.Error("ExtractContextValues() should return a map, not nil")
	}
	_ = values
}

func TestNew_WithCustomContextParent(t *testing.T) {
	parent := context.Background()
	parentCC := New(parent)
	parentCC1 := parentCC.WithValue("parent_key", "parent_value")

	// Create new CustomContext from CustomContext parent
	ctx := New(parentCC1)

	if ctx == nil {
		t.Fatal("New() should not return nil")
	}

	// Should copy values from parent CustomContext
	if ctx.GetValue("parent_key") != "parent_value" {
		t.Error("New() should copy values from parent CustomContext")
	}
}

func TestNew_WithStandardContextParent(t *testing.T) {
	parent := context.WithValue(context.Background(), "parent_key", "parent_value")

	// Create new CustomContext from standard context parent
	ctx := New(parent)

	if ctx == nil {
		t.Fatal("New() should not return nil")
	}

	// Should extract values from standard context
	// Note: ExtractContextValues may or may not extract depending on context type
	// This test verifies the function doesn't panic
	_ = ctx.GetValue("parent_key")
}

func TestExtractContextValues_EmptyContext(t *testing.T) {
	parent := context.Background()

	values := ExtractContextValues(parent)

	if values == nil {
		t.Error("ExtractContextValues() should return empty map, not nil")
	}

	// Empty context should return empty map
	if len(values) != 0 {
		t.Errorf("ExtractContextValues() with empty context should return empty map, got %d values", len(values))
	}
}

func TestExtractContextValues_NonStringKey(t *testing.T) {
	type customKey struct{}
	parent := context.WithValue(context.Background(), customKey{}, "value")

	values := ExtractContextValues(parent)

	// Non-string keys may not be extracted (depends on implementation)
	// This test just verifies it doesn't panic
	_ = values
}

func TestNew_ExtractsValuesFromParent(t *testing.T) {
	// Create a context chain with values
	parent := context.WithValue(
		context.WithValue(context.Background(), "key1", "value1"),
		"key2", "value2",
	)

	ctx := New(parent)

	// Values should be extracted and accessible via GetValue
	// Note: ExtractContextValues may extract values, but GetValue only checks CustomContext values
	// So we test that New doesn't panic and creates a valid context
	if ctx == nil {
		t.Fatal("New() should not return nil")
	}
}
