package http

import (
	"io"

	"github.com/reitmas32/rkit/core/customctx"
	corehttp "github.com/reitmas32/rkit/core/http"
)

// Get performs a GET request.
func (c *Client) Get(ctx *customctx.CustomContext, url string, opts ...corehttp.RequestOption) (*corehttp.Response, error) {
	req := corehttp.NewRequest("GET", url, opts...)
	return c.Do(ctx, req)
}

// Post performs a POST request with the given body.
func (c *Client) Post(ctx *customctx.CustomContext, url string, body io.Reader, opts ...corehttp.RequestOption) (*corehttp.Response, error) {
	req := corehttp.NewRequest("POST", url, opts...)
	req.Body = body
	return c.Do(ctx, req)
}

// Put performs a PUT request with the given body.
func (c *Client) Put(ctx *customctx.CustomContext, url string, body io.Reader, opts ...corehttp.RequestOption) (*corehttp.Response, error) {
	req := corehttp.NewRequest("PUT", url, opts...)
	req.Body = body
	return c.Do(ctx, req)
}

// Patch performs a PATCH request with the given body.
func (c *Client) Patch(ctx *customctx.CustomContext, url string, body io.Reader, opts ...corehttp.RequestOption) (*corehttp.Response, error) {
	req := corehttp.NewRequest("PATCH", url, opts...)
	req.Body = body
	return c.Do(ctx, req)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx *customctx.CustomContext, url string, opts ...corehttp.RequestOption) (*corehttp.Response, error) {
	req := corehttp.NewRequest("DELETE", url, opts...)
	return c.Do(ctx, req)
}

// Head performs a HEAD request.
func (c *Client) Head(ctx *customctx.CustomContext, url string, opts ...corehttp.RequestOption) (*corehttp.Response, error) {
	req := corehttp.NewRequest("HEAD", url, opts...)
	return c.Do(ctx, req)
}

// Options performs an OPTIONS request.
func (c *Client) Options(ctx *customctx.CustomContext, url string, opts ...corehttp.RequestOption) (*corehttp.Response, error) {
	req := corehttp.NewRequest("OPTIONS", url, opts...)
	return c.Do(ctx, req)
}
