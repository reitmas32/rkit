package result

import (
	"testing"

	"github.com/reitmas32/rkit/core/kerrors"
)

func TestNewResult(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		r := NewResult(42, nil)
		if !r.IsOk() {
			t.Error("Expected Result to be Ok when error is nil")
		}
		if r.Value() != 42 {
			t.Errorf("Expected value 42, got %d", r.Value())
		}
		err := r.Error()
		if err != nil {
			t.Errorf("Expected error to be nil, got %v", err)
		}
	})

	t.Run("with error", func(t *testing.T) {
		err := kerrors.NewKError("test error", 500, nil)
		r := NewResult(0, err)
		if r.IsOk() {
			t.Error("Expected Result to be Err when error is not nil")
		}
		if r.Error() != err {
			t.Error("Expected error to match the provided error")
		}
	})
}

func TestNewOkResult(t *testing.T) {
	r := NewOkResult("success")
	if !r.IsOk() {
		t.Error("Expected Result to be Ok")
	}
	if r.Value() != "success" {
		t.Errorf("Expected value 'success', got %s", r.Value())
	}
	err := r.Error()
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
}

func TestNewErrResult(t *testing.T) {
	err := kerrors.NewKError("error occurred", 404, nil)
	r := NewErrResult[string](err)
	if r.IsOk() {
		t.Error("Expected Result to be Err")
	}
	if r.Error() != err {
		t.Error("Expected error to match the provided error")
	}
	if r.Value() != "" {
		t.Error("Expected zero value for T")
	}
}

func TestOk(t *testing.T) {
	r := Ok(100)
	if !r.IsOk() {
		t.Error("Expected Result to be Ok")
	}
	if r.Value() != 100 {
		t.Errorf("Expected value 100, got %d", r.Value())
	}
}

func TestErr(t *testing.T) {
	err := kerrors.NewKError("failure", 500, nil)
	r := Err[int](err)
	if r.IsOk() {
		t.Error("Expected Result to be Err")
	}
	if r.Error() != err {
		t.Error("Expected error to match the provided error")
	}
}

func TestValue(t *testing.T) {
	t.Run("Ok Result", func(t *testing.T) {
		r := Ok("test")
		if r.Value() != "test" {
			t.Errorf("Expected 'test', got %s", r.Value())
		}
	})

	t.Run("Err Result", func(t *testing.T) {
		err := kerrors.NewKError("error", 500, nil)
		r := Err[string](err)
		if r.Value() != "" {
			t.Error("Expected zero value for Err Result")
		}
	})
}

func TestError(t *testing.T) {
	t.Run("Ok Result", func(t *testing.T) {
		r := Ok(42)
		err := r.Error()
		if err != nil {
			t.Errorf("Expected nil error for Ok Result, got %v", err)
		}
	})

	t.Run("Err Result", func(t *testing.T) {
		err := kerrors.NewKError("test error", 500, nil)
		r := Err[int](err)
		if r.Error() != err {
			t.Error("Expected error to match provided error")
		}
	})
}

func TestIsOk(t *testing.T) {
	t.Run("Ok Result", func(t *testing.T) {
		r := Ok(true)
		if !r.IsOk() {
			t.Error("Expected IsOk() to return true for Ok Result")
		}
	})

	t.Run("Err Result", func(t *testing.T) {
		err := kerrors.NewKError("error", 500, nil)
		r := Err[bool](err)
		if r.IsOk() {
			t.Error("Expected IsOk() to return false for Err Result")
		}
	})
}

func TestResultWithStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("Ok with struct", func(t *testing.T) {
		person := Person{Name: "John", Age: 30}
		r := Ok(person)
		if !r.IsOk() {
			t.Error("Expected Result to be Ok")
		}
		if r.Value().Name != "John" {
			t.Errorf("Expected Name 'John', got %s", r.Value().Name)
		}
		if r.Value().Age != 30 {
			t.Errorf("Expected Age 30, got %d", r.Value().Age)
		}
	})

	t.Run("Err with struct", func(t *testing.T) {
		err := kerrors.NewKError("person not found", 404, nil)
		r := Err[Person](err)
		if r.IsOk() {
			t.Error("Expected Result to be Err")
		}
		if r.Value() != (Person{}) {
			t.Error("Expected zero value for Person")
		}
	})
}
