package kerrors_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/reitmas32/rkit/core/kerrors"
)

func TestKindFromCode(t *testing.T) {
	cases := map[int]kerrors.Kind{
		400: kerrors.KindValidation,
		401: kerrors.KindUnauthorized,
		404: kerrors.KindNotFound,
		409: kerrors.KindConflict,
		500: kerrors.KindInternal,
		503: kerrors.KindUnavailable,
		200: kerrors.KindUnknown,
	}
	for code, want := range cases {
		if got := kerrors.KindFromCode(code); got != want {
			t.Errorf("KindFromCode(%d) = %q, want %q", code, got, want)
		}
	}
}

func TestKError_Detail(t *testing.T) {
	err := kerrors.NewKErrorWithCause(
		"database operation failed", 500,
		map[string]any{"table": "users", "op": "insert"},
		errors.New("connection refused"),
	)
	d := err.Detail()
	for _, want := range []string{"kind=internal", "code=500", "database operation failed", "op=insert", "table=users", "cause: connection refused"} {
		if !strings.Contains(d, want) {
			t.Errorf("Detail() = %q, missing %q", d, want)
		}
	}
	// Error() stays client-safe (message only).
	if err.Error() != "database operation failed" {
		t.Errorf("Error() = %q, want client-safe message", err.Error())
	}
}

func TestKError_WithMetadata_DoesNotMutateSentinel(t *testing.T) {
	sentinel := kerrors.NewValidation("invalid field", nil)
	derived := sentinel.WithMetadata("field", "email")

	if sentinel.Metadata != nil {
		t.Errorf("sentinel was mutated: %v", sentinel.Metadata)
	}
	if derived.Metadata["field"] != "email" {
		t.Errorf("derived metadata not set: %v", derived.Metadata)
	}
}

func TestKError_WithCause_PreservesKindAndIs(t *testing.T) {
	sentinel := kerrors.NewNotFound("item not found", nil)
	cause := errors.New("no rows")
	wrapped := sentinel.WithCause(cause)

	if wrapped.Kind != kerrors.KindNotFound {
		t.Errorf("WithCause dropped Kind: %q", wrapped.Kind)
	}
	if !errors.Is(wrapped, cause) {
		t.Error("errors.Is should find the wrapped cause")
	}
	// Is() matches by code+kind, so the wrapped error still matches the sentinel.
	if !errors.Is(wrapped, sentinel) {
		t.Error("wrapped error should match its sentinel by code+kind")
	}
}
