package http

// WithExpectedStatusCode sets the expected status code for the response.
// When set, IsSuccess() will return true only if the StatusCode matches this value.
func WithExpectedStatusCode[T any](code int) func(*TypedResponse[T]) {
	return func(resp *TypedResponse[T]) {
		resp.ExpectedStatusCode = &code
	}
}

// WithSuccessStatusCodeRange sets the range of status codes considered successful.
// When set, IsSuccess() will return true if StatusCode is within [Min, Max] (inclusive).
// If both ExpectedStatusCode and SuccessStatusCodeRange are set, ExpectedStatusCode takes precedence.
func WithSuccessStatusCodeRange[T any](min, max int) func(*TypedResponse[T]) {
	return func(resp *TypedResponse[T]) {
		resp.SuccessStatusCodeRange = &StatusCodeRange{
			Min: min,
			Max: max,
		}
	}
}

// SetExpectedStatusCode sets the expected status code for the response.
// This is a convenience method that can be called on the response after it's created.
func (r *TypedResponse[T]) SetExpectedStatusCode(code int) {
	r.ExpectedStatusCode = &code
}

// SetSuccessStatusCodeRange sets the range of status codes considered successful.
// This is a convenience method that can be called on the response after it's created.
func (r *TypedResponse[T]) SetSuccessStatusCodeRange(min, max int) {
	r.SuccessStatusCodeRange = &StatusCodeRange{
		Min: min,
		Max: max,
	}
}
