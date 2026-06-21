// Package http provides a concrete HTTP Client built on net/http with retries,
// request/response logging and customctx integration. It implements the
// contracts defined in github.com/reitmas32/rkit/core/http.
//
// Create a client with NewClient(DefaultConfig()) and issue requests with Get,
// Post, Put, Patch, Delete, Head and Options, or use the generic typed helpers
// in core/http for automatic JSON (de)serialization. NewCustomContextFromGin
// bridges a *gin.Context into a *customctx.CustomContext.
//
//	import infrahttp "github.com/reitmas32/rkit/infrastructure/http"
//
//	client := infrahttp.NewClient(infrahttp.DefaultConfig())
package http
