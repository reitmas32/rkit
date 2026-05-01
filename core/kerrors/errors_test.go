package kerrors

import (
	"errors"
	"testing"
)

func TestNewKError(t *testing.T) {
	metadata := map[string]any{
		"request_id": "12345",
		"user_id":    42,
	}

	err := NewKError("test error", 500, metadata)

	if err == nil {
		t.Fatal("Expected error to be non-nil")
	}

	if err.Message != "test error" {
		t.Errorf("Expected message 'test error', got %s", err.Message)
	}

	if err.Code != 500 {
		t.Errorf("Expected code 500, got %d", err.Code)
	}

	if err.Metadata == nil {
		t.Error("Expected metadata to be non-nil")
	}

	if err.Metadata["request_id"] != "12345" {
		t.Errorf("Expected request_id '12345', got %v", err.Metadata["request_id"])
	}

	if err.Cause != nil {
		t.Error("Expected Cause to be nil")
	}
}

func TestNewKErrorWithCause(t *testing.T) {
	originalErr := errors.New("original error")
	metadata := map[string]any{"key": "value"}

	err := NewKErrorWithCause("wrapped error", 400, metadata, originalErr)

	if err == nil {
		t.Fatal("Expected error to be non-nil")
	}

	if err.Message != "wrapped error" {
		t.Errorf("Expected message 'wrapped error', got %s", err.Message)
	}

	if err.Cause != originalErr {
		t.Error("Expected Cause to match original error")
	}
}

func TestError(t *testing.T) {
	err := NewKError("test message", 500, nil)

	if err.Error() != "test message" {
		t.Errorf("Expected 'test message', got %s", err.Error())
	}
}

func TestUnwrap(t *testing.T) {
	t.Run("without cause", func(t *testing.T) {
		err := NewKError("error", 500, nil)
		if err.Unwrap() != nil {
			t.Error("Expected Unwrap() to return nil when no cause is set")
		}
	})

	t.Run("with cause", func(t *testing.T) {
		originalErr := errors.New("original")
		err := NewKErrorWithCause("wrapped", 500, nil, originalErr)

		unwrapped := err.Unwrap()
		if unwrapped != originalErr {
			t.Error("Expected Unwrap() to return the original error")
		}
	})
}

func TestErrorChaining(t *testing.T) {
	// Test error chain with errors.Unwrap
	originalErr := errors.New("root cause")
	level1 := NewKErrorWithCause("level 1", 500, nil, originalErr)
	level2 := NewKErrorWithCause("level 2", 400, nil, level1)

	// Test unwrapping
	unwrapped1 := level2.Unwrap()
	if unwrapped1 != level1 {
		t.Error("Expected level2.Unwrap() to return level1")
	}

	unwrapped2 := level1.Unwrap()
	if unwrapped2 != originalErr {
		t.Error("Expected level1.Unwrap() to return originalErr")
	}

	// Test errors.Is
	if !errors.Is(level2, originalErr) {
		t.Error("Expected errors.Is to find originalErr in chain")
	}
}

func TestKErrorAsError(t *testing.T) {
	// Test that KError implements the error interface
	var err error = NewKError("test", 500, nil)

	if err.Error() != "test" {
		t.Errorf("Expected 'test', got %s", err.Error())
	}
}

func TestKErrorWithNilMetadata(t *testing.T) {
	err := NewKError("error", 500, nil)

	if err.Metadata != nil {
		t.Error("Expected metadata to be nil when nil is passed")
	}
}

func TestKErrorWithEmptyMetadata(t *testing.T) {
	err := NewKError("error", 500, map[string]any{})

	if err.Metadata == nil {
		t.Error("Expected metadata to be non-nil empty map")
	}

	if len(err.Metadata) != 0 {
		t.Error("Expected metadata to be empty")
	}
}

