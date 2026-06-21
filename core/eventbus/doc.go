// Package eventbus defines the core, transport-agnostic contracts for publishing
// and consuming events in a decoupled way: Event, Message, Publisher, Consumer,
// EventBus, HandlerFunc and Metadata.
//
// This package contains interfaces only. Concrete transports live under
// github.com/reitmas32/rkit/infrastructure/eventbus (in-memory and RabbitMQ),
// so application code can depend on these contracts and swap the backend freely.
//
//	import "github.com/reitmas32/rkit/core/eventbus"
package eventbus
