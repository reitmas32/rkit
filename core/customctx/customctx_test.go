package customctx

import (
	"context"
	"testing"
	"time"

	"github.com/reitmas32/rkit/core/kerrors"
)

func TestNew(t *testing.T) {
	parent := context.Background()
	ctx := New(parent)

	if ctx == nil {
		t.Fatal("Expected CustomContext to be non-nil")
	}

	if ctx.Context() != parent {
		t.Error("Expected Context() to return parent context")
	}
}

func TestNewCustomContext(t *testing.T) {
	parent := context.Background()
	ctx := NewCustomContext(parent)

	if ctx == nil {
		t.Fatal("Expected CustomContext to be non-nil")
	}
}

func TestContextInterface(t *testing.T) {
	parent := context.WithValue(context.Background(), "key", "value")
	ctx := New(parent)

	// Test Deadline
	deadline, ok := ctx.Deadline()
	parentDeadline, parentOk := parent.Deadline()
	if ok != parentOk {
		t.Error("Expected Deadline ok to match parent")
	}
	if deadline != parentDeadline {
		t.Error("Expected Deadline to match parent")
	}

	// Test Done
	if ctx.Done() != parent.Done() {
		t.Error("Expected Done() to return parent's Done channel")
	}

	// Test Err
	if ctx.Err() != parent.Err() {
		t.Error("Expected Err() to match parent")
	}

	// Test Value
	if ctx.Value("key") != "value" {
		t.Error("Expected Value() to return parent's value")
	}
}

func TestAddError(t *testing.T) {
	ctx := New(context.Background())

	err := kerrors.NewKError("test error", 500, nil)
	result := ctx.AddError(err)

	if result != err {
		t.Error("Expected AddError to return the error")
	}

	if !ctx.HasErrors() {
		t.Error("Expected HasErrors() to return true after adding error")
	}

	errors := ctx.Errors()
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}

	if errors[0].Error != err {
		t.Error("Expected error to match added error")
	}

	if errors[0].CallIn == "" {
		t.Error("Expected CallIn to be set")
	}
}

func TestAddErrorNil(t *testing.T) {
	ctx := New(context.Background())

	result := ctx.AddError(nil)

	if result != nil {
		t.Error("Expected AddError(nil) to return nil")
	}

	if ctx.HasErrors() {
		t.Error("Expected HasErrors() to return false after adding nil error")
	}
}

func TestHasErrors(t *testing.T) {
	ctx := New(context.Background())

	if ctx.HasErrors() {
		t.Error("Expected HasErrors() to return false initially")
	}

	ctx.AddError(kerrors.NewKError("error 1", 500, nil))
	if !ctx.HasErrors() {
		t.Error("Expected HasErrors() to return true after adding error")
	}
}

func TestErrors(t *testing.T) {
	ctx := New(context.Background())

	err1 := kerrors.NewKError("error 1", 500, nil)
	err2 := kerrors.NewKError("error 2", 404, nil)

	ctx.AddError(err1)
	ctx.AddError(err2)

	errors := ctx.Errors()
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}

	if errors[0].Error != err1 {
		t.Error("Expected first error to match err1")
	}

	if errors[1].Error != err2 {
		t.Error("Expected second error to match err2")
	}
}

func TestFirstError(t *testing.T) {
	ctx := New(context.Background())

	err1 := kerrors.NewKError("first", 500, nil)
	err2 := kerrors.NewKError("second", 404, nil)

	ctx.AddError(err1)
	ctx.AddError(err2)

	first := ctx.FirstError()
	if first.Error != err1 {
		t.Error("Expected FirstError() to return first error")
	}
}

func TestLastError(t *testing.T) {
	ctx := New(context.Background())

	err1 := kerrors.NewKError("first", 500, nil)
	err2 := kerrors.NewKError("second", 404, nil)

	ctx.AddError(err1)
	ctx.AddError(err2)

	last := ctx.LastError()
	if last.Error != err2 {
		t.Error("Expected LastError() to return last error")
	}
}

func TestFirstErrorPanic(t *testing.T) {
	ctx := New(context.Background())

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected FirstError() to panic when no errors")
		}
	}()

	ctx.FirstError()
}

func TestLastErrorPanic(t *testing.T) {
	ctx := New(context.Background())

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected LastError() to panic when no errors")
		}
	}()

	ctx.LastError()
}

func TestClear(t *testing.T) {
	ctx := New(context.Background())

	ctx.AddError(kerrors.NewKError("error", 500, nil))
	ctx.Clear()

	if ctx.HasErrors() {
		t.Error("Expected HasErrors() to return false after Clear()")
	}

	errors := ctx.Errors()
	if len(errors) != 0 {
		t.Errorf("Expected 0 errors after Clear(), got %d", len(errors))
	}
}

func TestConcurrentAccess(t *testing.T) {
	ctx := New(context.Background())

	// Test concurrent AddError
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			err := kerrors.NewKError("error", id, nil)
			ctx.AddError(err)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	if !ctx.HasErrors() {
		t.Error("Expected HasErrors() to return true after concurrent adds")
	}

	errors := ctx.Errors()
	if len(errors) != 10 {
		t.Errorf("Expected 10 errors, got %d", len(errors))
	}
}

func TestWithTimeout(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ctx := New(parent)

	// Test that timeout is propagated
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Expected Deadline to be set from parent")
	}

	if deadline.IsZero() {
		t.Error("Expected Deadline to be non-zero")
	}

	// Test Done channel
	select {
	case <-ctx.Done():
		t.Error("Expected Done channel to not be closed immediately")
	case <-time.After(50 * time.Millisecond):
		// Expected
	}

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Expected Done channel to be closed after timeout")
	}
}

func TestNewError(t *testing.T) {
	ctx := New(context.Background())

	err := kerrors.NewKError("test", 500, nil)
	result := ctx.NewError(err)

	if result != err {
		t.Error("Expected NewError to return the error")
	}

	if !ctx.HasErrors() {
		t.Error("Expected HasErrors() to return true after NewError")
	}
}

// --- Value storage tests ---

func TestWithValue(t *testing.T) {
	ctx := New(context.Background())

	// Test basic value storage
	ctx1 := ctx.WithValue("request_id", "req-123")

	if !ctx1.HasValue("request_id") {
		t.Error("Expected HasValue() to return true after WithValue")
	}

	val := ctx1.GetValue("request_id")
	if val != "req-123" {
		t.Errorf("Expected GetValue() to return 'req-123', got %v", val)
	}

	// Test immutability - original context should not have the value
	if ctx.HasValue("request_id") {
		t.Error("Expected original context to not have the value")
	}

	if ctx.GetValue("request_id") != nil {
		t.Error("Expected original context GetValue() to return nil")
	}
}

func TestWithValueMultiple(t *testing.T) {
	ctx := New(context.Background())

	ctx1 := ctx.WithValue("request_id", "req-123")
	ctx2 := ctx1.WithValue("trace_id", "trace-456")
	ctx3 := ctx2.WithValue("user_id", 789)

	// All values should be accessible
	if ctx3.GetValue("request_id") != "req-123" {
		t.Error("Expected first value to be preserved")
	}

	if ctx3.GetValue("trace_id") != "trace-456" {
		t.Error("Expected second value to be accessible")
	}

	if ctx3.GetValue("user_id") != 789 {
		t.Error("Expected third value to be accessible")
	}

	// Previous contexts should still have their values
	if ctx1.GetValue("request_id") != "req-123" {
		t.Error("Expected ctx1 to still have its value")
	}

	if ctx1.GetValue("trace_id") != nil {
		t.Error("Expected ctx1 to not have values added later")
	}
}

func TestValueMethodChecksCustomContextFirst(t *testing.T) {
	// Add value to parent context (using string key)
	parent := context.WithValue(context.Background(), "request_id", "parent-value")
	ctxWithParent := New(parent)

	// Add value to CustomContext with same key
	ctxWithValue := ctxWithParent.WithValue("request_id", "custom-value")

	// Value() should return CustomContext value first (not parent value) when key is string
	val := ctxWithValue.Value("request_id")
	if val != "custom-value" {
		t.Errorf("Expected Value() to return CustomContext value 'custom-value', got %v", val)
	}
}

func TestValueMethodFallsBackToParent(t *testing.T) {
	parent := context.WithValue(context.Background(), "request_id", "parent-value")
	ctx := New(parent)

	// Value() should check CustomContext first (for string keys), then parent
	val := ctx.Value("request_id")
	if val != "parent-value" {
		t.Errorf("Expected Value() to return parent value 'parent-value', got %v", val)
	}

	// If not in CustomContext or parent, should return nil
	if ctx.Value("unknown_key") != nil {
		t.Error("Expected Value() to return nil for unknown key")
	}

	// Non-string keys should always delegate to parent
	type unknownKey struct{}
	parentWithStruct := context.WithValue(context.Background(), unknownKey{}, "struct-value")
	ctxWithStruct := New(parentWithStruct)
	valStruct := ctxWithStruct.Value(unknownKey{})
	if valStruct != "struct-value" {
		t.Errorf("Expected Value() to return parent value for non-string key, got %v", valStruct)
	}
}

func TestGetValue(t *testing.T) {
	ctx := New(context.Background())

	// Initially should return nil
	if ctx.GetValue("request_id") != nil {
		t.Error("Expected GetValue() to return nil initially")
	}

	// After WithValue, should return the value
	ctx1 := ctx.WithValue("request_id", "req-123")
	if ctx1.GetValue("request_id") != "req-123" {
		t.Error("Expected GetValue() to return the value")
	}

	// GetValue should not check parent
	parent := context.WithValue(context.Background(), "request_id", "parent-value")
	ctxWithParent := New(parent)
	if ctxWithParent.GetValue("request_id") != nil {
		t.Error("Expected GetValue() to not check parent context")
	}
}

func TestHasValue(t *testing.T) {
	ctx := New(context.Background())

	// Initially should return false
	if ctx.HasValue("request_id") {
		t.Error("Expected HasValue() to return false initially")
	}

	// After WithValue, should return true
	ctx1 := ctx.WithValue("request_id", "req-123")
	if !ctx1.HasValue("request_id") {
		t.Error("Expected HasValue() to return true after WithValue")
	}

	// HasValue should not check parent
	parent := context.WithValue(context.Background(), "request_id", "parent-value")
	ctxWithParent := New(parent)
	if ctxWithParent.HasValue("request_id") {
		t.Error("Expected HasValue() to not check parent context")
	}

	// Different key should return false
	if ctx1.HasValue("other_key") {
		t.Error("Expected HasValue() to return false for different key")
	}
}

func TestWithValueOverwritesExisting(t *testing.T) {
	ctx := New(context.Background())

	ctx1 := ctx.WithValue("request_id", "req-123")
	ctx2 := ctx1.WithValue("request_id", "req-456")

	// New context should have the new value
	if ctx2.GetValue("request_id") != "req-456" {
		t.Error("Expected WithValue to overwrite existing value")
	}

	// Old context should still have the old value (immutability)
	if ctx1.GetValue("request_id") != "req-123" {
		t.Error("Expected old context to still have old value")
	}
}

func TestWithValueAndErrors(t *testing.T) {
	ctx := New(context.Background())

	err := kerrors.NewKError("error", 500, nil)
	ctx.AddError(err)

	ctx1 := ctx.WithValue("request_id", "req-123")

	// Both context should have the error
	if !ctx.HasErrors() {
		t.Error("Expected original context to have error")
	}

	if !ctx1.HasErrors() {
		t.Error("Expected new context to have error")
	}

	// New context should have the value
	if ctx1.GetValue("request_id") != "req-123" {
		t.Error("Expected new context to have value")
	}

	// Original context should not have the value
	if ctx.GetValue("request_id") != nil {
		t.Error("Expected original context to not have value")
	}
}

func TestConcurrentWithValue(t *testing.T) {
	ctx := New(context.Background())

	// Test concurrent WithValue calls
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			ctx.WithValue("request_id", id)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Original context should not have values (immutability)
	if ctx.HasValue("request_id") {
		t.Error("Expected original context to not have values from concurrent calls")
	}
}

func TestWithValueNilValue(t *testing.T) {
	ctx := New(context.Background())

	// Should allow nil values
	ctx1 := ctx.WithValue("request_id", nil)

	if !ctx1.HasValue("request_id") {
		t.Error("Expected HasValue() to return true even for nil value")
	}

	if ctx1.GetValue("request_id") != nil {
		t.Error("Expected GetValue() to return nil")
	}
}
