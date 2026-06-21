package kerrors

// Kind is a coarse classification of an error. It lets callers branch on the
// class of failure (and map it to a transport status) without hard-coding
// numeric codes.
type Kind string

const (
	// KindUnknown is the zero value, used when no category was set.
	KindUnknown Kind = "unknown"
	// KindValidation: the input failed validation (bad/illegal value). HTTP 400/422.
	KindValidation Kind = "validation"
	// KindUnauthorized: authentication is missing or invalid. HTTP 401.
	KindUnauthorized Kind = "unauthorized"
	// KindForbidden: authenticated but not allowed. HTTP 403.
	KindForbidden Kind = "forbidden"
	// KindNotFound: the requested resource does not exist. HTTP 404.
	KindNotFound Kind = "not_found"
	// KindConflict: the operation conflicts with current state (e.g. duplicate). HTTP 409.
	KindConflict Kind = "conflict"
	// KindRateLimited: the caller exceeded a rate limit. HTTP 429.
	KindRateLimited Kind = "rate_limited"
	// KindExternal: a downstream/external dependency failed. HTTP 502.
	KindExternal Kind = "external"
	// KindUnavailable: the service or a dependency is temporarily unavailable. HTTP 503.
	KindUnavailable Kind = "unavailable"
	// KindTimeout: an operation exceeded its deadline. HTTP 504.
	KindTimeout Kind = "timeout"
	// KindInternal: an unexpected internal failure. HTTP 500.
	KindInternal Kind = "internal"
)

// KindFromCode maps an HTTP-style status code to a Kind. It is used as the
// default classification when an error is built from a raw code.
func KindFromCode(code int) Kind {
	switch code {
	case 400, 422:
		return KindValidation
	case 401:
		return KindUnauthorized
	case 403:
		return KindForbidden
	case 404:
		return KindNotFound
	case 409:
		return KindConflict
	case 429:
		return KindRateLimited
	case 502:
		return KindExternal
	case 503:
		return KindUnavailable
	case 504:
		return KindTimeout
	}
	if code >= 500 {
		return KindInternal
	}
	if code >= 400 {
		return KindValidation
	}
	return KindUnknown
}

// The category constructors below set a sensible default Code and a Kind so the
// reason for a failure is unambiguous. Codes can still be overridden afterwards.

// NewValidation builds a validation error (code 422). Use the metadata to point
// at the offending field(s) and the rule that was violated.
func NewValidation(message string, metadata map[string]any) *KError {
	return &KError{Message: message, Kind: KindValidation, Code: 422, Metadata: metadata}
}

// NewNotFound builds a not-found error (code 404).
func NewNotFound(message string, metadata map[string]any) *KError {
	return &KError{Message: message, Kind: KindNotFound, Code: 404, Metadata: metadata}
}

// NewConflict builds a conflict error (code 409), e.g. a duplicate key.
func NewConflict(message string, metadata map[string]any) *KError {
	return &KError{Message: message, Kind: KindConflict, Code: 409, Metadata: metadata}
}

// NewUnauthorized builds an authentication error (code 401).
func NewUnauthorized(message string, metadata map[string]any) *KError {
	return &KError{Message: message, Kind: KindUnauthorized, Code: 401, Metadata: metadata}
}

// NewForbidden builds an authorization error (code 403).
func NewForbidden(message string, metadata map[string]any) *KError {
	return &KError{Message: message, Kind: KindForbidden, Code: 403, Metadata: metadata}
}

// NewInternal builds an internal error (code 500), optionally wrapping a cause.
func NewInternal(message string, cause error) *KError {
	return &KError{Message: message, Kind: KindInternal, Code: 500, Cause: cause}
}

// NewExternal builds a downstream/dependency error (code 502), optionally
// wrapping the cause returned by the dependency.
func NewExternal(message string, cause error) *KError {
	return &KError{Message: message, Kind: KindExternal, Code: 502, Cause: cause}
}
