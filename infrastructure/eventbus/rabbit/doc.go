// Package rabbit implements the core/eventbus contracts on top of RabbitMQ
// (github.com/rabbitmq/amqp091-go). It supports publishing and consuming events
// across services, including delayed/scheduled delivery.
//
// See the runnable programs under examples/infrastructure/eventbus/rabbit for
// publisher, consumer and delayed-delivery patterns.
//
//	import "github.com/reitmas32/rkit/infrastructure/eventbus/rabbit"
package rabbit
