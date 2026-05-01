package result

import (
	"net/http"
	"testing"

	"github.com/reitmas32/rkit/core/kerrors"
)

func TestEmpty(t *testing.T) {
	r := Empty[string]()
	// Empty Result has zero value and no error, so IsEmpty should be true
	// But IsOk() returns true when _error is nil, which is the case for Empty
	if !r.IsEmpty() {
		t.Error("Empty() should return an empty Result")
	}
	// Empty Result has no error, so IsOk() returns true
	if !r.IsOk() {
		t.Error("Empty() Result should be Ok (no error)")
	}
	if r.Error() != nil {
		t.Error("Empty() Result should not have an error")
	}
}

func TestResult_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		result   Result[string]
		expected bool
	}{
		{"empty result", Empty[string](), true},
		{"ok result with value", Ok("test"), false},
		{"err result", Err[string](kerrors.NewKError("error", http.StatusInternalServerError, nil)), false},
		{"ok result with empty string", Ok(""), true}, // Empty string is zero value, so IsEmpty returns true
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.IsEmpty(); got != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestResult_ToKError(t *testing.T) {
	tests := []struct {
		name     string
		result   Result[string]
		expected *kerrors.KError
	}{
		{"ok result", Ok("test"), nil},
		{"err result", Err[string](kerrors.NewKError("error", http.StatusInternalServerError, nil)), kerrors.NewKError("error", http.StatusInternalServerError, nil)},
		{"empty result", Empty[string](), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.result.ToKError()
			if tt.expected == nil {
				if got != nil {
					t.Errorf("ToKError() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("ToKError() = nil, want %v", tt.expected)
				} else if got.Error() != tt.expected.Error() {
					t.Errorf("ToKError() = %v, want %v", got.Error(), tt.expected.Error())
				}
			}
		})
	}
}

func TestNewResult_Additional(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		err      *kerrors.KError
		wantOk   bool
		wantErr  bool
	}{
		{"result with value and no error", "test", nil, true, false},
		{"result with value and error", "test", kerrors.NewKError("error", http.StatusInternalServerError, nil), false, true},
		{"result with empty value and no error", "", nil, true, false},
		{"result with empty value and error", "", kerrors.NewKError("error", http.StatusInternalServerError, nil), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewResult(tt.value, tt.err)
			if r.IsOk() != tt.wantOk {
				t.Errorf("IsOk() = %v, want %v", r.IsOk(), tt.wantOk)
			}
			if (r.Error() != nil) != tt.wantErr {
				t.Errorf("Error() != nil = %v, want %v", r.Error() != nil, tt.wantErr)
			}
			if r.Value() != tt.value {
				t.Errorf("Value() = %v, want %v", r.Value(), tt.value)
			}
		})
	}
}

func TestResult_Value_WithError(t *testing.T) {
	err := kerrors.NewKError("test error", http.StatusInternalServerError, nil)
	r := Err[string](err)

	// Value should return zero value even when there's an error
	value := r.Value()
	if value != "" {
		t.Errorf("Value() = %v, want empty string", value)
	}
}

func TestResult_Error_WithOk(t *testing.T) {
	r := Ok("test")
	if r.Error() != nil {
		t.Errorf("Error() = %v, want nil", r.Error())
	}
}

func TestResult_Error_WithErr(t *testing.T) {
	err := kerrors.NewKError("test error", http.StatusInternalServerError, nil)
	r := Err[string](err)
	if r.Error() == nil {
		t.Error("Error() = nil, want error")
	}
	if r.Error().Error() != "test error" {
		t.Errorf("Error().Error() = %v, want 'test error'", r.Error().Error())
	}
}

func TestOk_Additional(t *testing.T) {
	r := Ok("test")
	if !r.IsOk() {
		t.Error("Ok() should return a successful Result")
	}
	if r.Error() != nil {
		t.Error("Ok() Result should not have an error")
	}
	if r.Value() != "test" {
		t.Errorf("Ok() Value() = %v, want 'test'", r.Value())
	}
}

func TestErr_Additional(t *testing.T) {
	err := kerrors.NewKError("test error", http.StatusInternalServerError, nil)
	r := Err[string](err)
	if r.IsOk() {
		t.Error("Err() should return a failed Result")
	}
	if r.Error() == nil {
		t.Error("Err() Result should have an error")
	}
}

func TestResult_WithIntType(t *testing.T) {
	r := Ok(42)
	if !r.IsOk() {
		t.Error("Result with int should be Ok")
	}
	if r.Value() != 42 {
		t.Errorf("Value() = %v, want 42", r.Value())
	}
}

func TestResult_WithStructType(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	r := Ok(TestStruct{Name: "test", Age: 30})
	if !r.IsOk() {
		t.Error("Result with struct should be Ok")
	}
	if r.Value().Name != "test" {
		t.Errorf("Value().Name = %v, want 'test'", r.Value().Name)
	}
	if r.Value().Age != 30 {
		t.Errorf("Value().Age = %v, want 30", r.Value().Age)
	}
}

func TestResult_WithSliceType(t *testing.T) {
	r := Ok([]string{"a", "b", "c"})
	if !r.IsOk() {
		t.Error("Result with slice should be Ok")
	}
	if len(r.Value()) != 3 {
		t.Errorf("Value() length = %v, want 3", len(r.Value()))
	}
}

func TestResult_WithMapType(t *testing.T) {
	r := Ok(map[string]int{"a": 1, "b": 2})
	if !r.IsOk() {
		t.Error("Result with map should be Ok")
	}
	if r.Value()["a"] != 1 {
		t.Errorf("Value()['a'] = %v, want 1", r.Value()["a"])
	}
}
