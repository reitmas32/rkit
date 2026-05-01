package http

import (
	"io"

	"github.com/reitmas32/rkit/core/customctx"
)

// Client defines the interface for making HTTP requests.
// It abstracts away the specific HTTP client implementation,
// allowing for different implementations (standard library, http.Client, etc.)
// and easier testing with mocks.
type Client interface {
	// Do performs an HTTP request with the given context and request options.
	// It returns a Response or an error.
	Do(ctx *customctx.CustomContext, req *Request) (*Response, error)

	// Get performs a GET request.
	Get(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)

	// Post performs a POST request with the given body.
	Post(ctx *customctx.CustomContext, url string, body io.Reader, opts ...RequestOption) (*Response, error)

	// Put performs a PUT request with the given body.
	Put(ctx *customctx.CustomContext, url string, body io.Reader, opts ...RequestOption) (*Response, error)

	// Patch performs a PATCH request with the given body.
	Patch(ctx *customctx.CustomContext, url string, body io.Reader, opts ...RequestOption) (*Response, error)

	// Delete performs a DELETE request.
	Delete(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)

	// Head performs a HEAD request.
	Head(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)

	// Options performs an OPTIONS request.
	Options(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)
}
