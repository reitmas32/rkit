// Package dtos provides helpers to bind and validate request DTOs from a
// *gin.Context into typed values, integrating with customctx error accumulation
// and the result.Result type.
//
// Key helpers: GetDTO and GetDTOWithResponse (destructive bind + validate),
// GetDTONonDestructive (returns a result.Result without writing a response), and
// GetAuthToken / GetAuthTokenWithEarlyResponse for extracting bearer tokens.
//
//	import "github.com/reitmas32/rkit/infrastructure/dtos"
package dtos
