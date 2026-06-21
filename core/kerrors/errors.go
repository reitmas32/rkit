// Package kerrors provides a custom error type with structured information
// including error codes, an error category (Kind), metadata, and cause chaining.
//
// The design separates two audiences:
//
//   - Error() returns only the human-readable Message and is therefore safe to
//     surface to API clients. Keep secrets, SQL, and internal detail out of it.
//   - Detail() renders the full context (kind, code, metadata and the unwrapped
//     cause chain) and is meant for logs, so it is always clear *why* a failure
//     happened.
package kerrors

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// KError represents a structured error with a message, an error category, a
// numeric code, optional metadata, and an optional underlying cause.
// It implements the error interface and supports unwrapping.
type KError struct {
	// Message is the human-readable, client-safe message describing what went wrong.
	Message string

	// Kind is a coarse category (validation, not_found, internal, ...) that lets
	// callers branch on the class of error without parsing the numeric Code.
	// It is optional; the zero value is KindUnknown.
	Kind Kind

	// Code is a numeric error code, typically aligned with HTTP status codes.
	Code int

	// Metadata carries additional, structured context about the error
	// (request IDs, offending field names, limits, etc.). It is intended for
	// logs/telemetry and may be exposed to clients only in non-production builds.
	Metadata map[string]any

	// Cause is the underlying error that triggered this one, if any.
	// It enables errors.Is / errors.Unwrap and is rendered by Detail().
	Cause error
}

// NewKError creates a new KError with the given message, code, and metadata.
// The Kind is inferred from the code (see KindFromCode).
func NewKError(message string, code int, metadata map[string]any) *KError {
	return &KError{Message: message, Kind: KindFromCode(code), Code: code, Metadata: metadata}
}

// NewKErrorWithCause creates a new KError that wraps an underlying cause,
// preserving the original error chain for errors.Is / errors.Unwrap.
func NewKErrorWithCause(message string, code int, metadata map[string]any, cause error) *KError {
	return &KError{Message: message, Kind: KindFromCode(code), Code: code, Metadata: metadata, Cause: cause}
}

// Error returns the client-safe Message, implementing the error interface.
func (e *KError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Message
}

// Unwrap returns the underlying cause, enabling errors.Is / errors.Unwrap.
func (e *KError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// Is reports whether target is a *KError with the same Code and Kind, so
// callers can match sentinel errors regardless of metadata or cause.
func (e *KError) Is(target error) bool {
	var t *KError
	if !errors.As(target, &t) || e == nil || t == nil {
		return false
	}
	return e.Code == t.Code && e.Kind == t.Kind
}

// clone returns a shallow copy with a copied metadata map so chained mutations
// never alter a shared sentinel error.
func (e *KError) clone() *KError {
	cp := &KError{Message: e.Message, Kind: e.Kind, Code: e.Code, Cause: e.Cause}
	if e.Metadata != nil {
		cp.Metadata = make(map[string]any, len(e.Metadata))
		for k, v := range e.Metadata {
			cp.Metadata[k] = v
		}
	}
	return cp
}

// WithCause returns a copy of the error with the given cause attached.
// The original (often a package-level sentinel) is left untouched.
func (e *KError) WithCause(cause error) *KError {
	c := e.clone()
	c.Cause = cause
	return c
}

// WithKind returns a copy of the error with its category overridden.
func (e *KError) WithKind(kind Kind) *KError {
	c := e.clone()
	c.Kind = kind
	return c
}

// WithMetadata returns a copy of the error with an extra metadata entry added.
// It is safe to call on a shared sentinel error: the sentinel is never mutated.
func (e *KError) WithMetadata(key string, value any) *KError {
	c := e.clone()
	if c.Metadata == nil {
		c.Metadata = make(map[string]any, 1)
	}
	c.Metadata[key] = value
	return c
}

// WithMetadataMap returns a copy of the error with all entries of m merged in.
func (e *KError) WithMetadataMap(m map[string]any) *KError {
	c := e.clone()
	if len(m) == 0 {
		return c
	}
	if c.Metadata == nil {
		c.Metadata = make(map[string]any, len(m))
	}
	for k, v := range m {
		c.Metadata[k] = v
	}
	return c
}

// Detail renders the full, log-oriented representation of the error:
// "[kind=... code=...] message | meta={...} | cause: <chain>".
// Unlike Error(), it is meant for logs and may include internal context, so do
// not return it to API clients.
func (e *KError) Detail() string {
	if e == nil {
		return "<nil>"
	}
	var b strings.Builder
	kind := e.Kind
	if kind == "" {
		kind = KindUnknown
	}
	fmt.Fprintf(&b, "[kind=%s code=%d] %s", kind, e.Code, e.Message)

	if len(e.Metadata) > 0 {
		keys := make([]string, 0, len(e.Metadata))
		for k := range e.Metadata {
			keys = append(keys, k)
		}
		sort.Strings(keys) // deterministic output for logs/tests
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s=%v", k, e.Metadata[k]))
		}
		fmt.Fprintf(&b, " | meta={%s}", strings.Join(parts, " "))
	}

	if e.Cause != nil {
		fmt.Fprintf(&b, " | cause: %s", e.Cause.Error())
	}
	return b.String()
}
