// Package inmemory provides an in-process implementation of the core/eventbus
// contracts. It dispatches published events to registered consumers within the
// same process, which makes it ideal for unit tests, local development and
// monolithic deployments that do not need an external broker.
//
//	import "github.com/reitmas32/rkit/infrastructure/eventbus/inmemory"
//
// For a distributed transport, use
// github.com/reitmas32/rkit/infrastructure/eventbus/rabbit instead.
package inmemory
