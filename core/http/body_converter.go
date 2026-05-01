package http

import "io"

// convertBodyToReader converts a request body to an io.Reader.
// It handles the following cases:
//   - If body is nil (zero value), returns nil
//   - If body implements RequestBody, uses ToReader()
//   - If body is already an io.Reader, returns it directly
//   - Otherwise, JSON marshals the body automatically
func convertBodyToReader[TRequest any](body TRequest) (io.Reader, error) {
	// Check if body is nil (for zero values)
	var zeroTRequest TRequest
	if any(body) == any(zeroTRequest) {
		return nil, nil
	}

	// Check if body implements RequestBody interface
	if reqBody, ok := any(body).(RequestBody); ok {
		return reqBody.ToReader()
	}

	// Check if body is already an io.Reader
	if reader, ok := any(body).(io.Reader); ok {
		return reader, nil
	}

	// JSON marshal the body automatically
	jsonBody := NewJSONRequestBody(body)
	return jsonBody.ToReader()
}
