package eventbus

import "time"

// Message represents a message delivery from a message broker or event bus.
// It abstracts away the specific implementation details (AMQP, in-memory, etc.)
// and provides a common interface for consuming events.
type Message interface {
	// Event returns the event associated with this message.
	Event() Event

	// Ack acknowledges the message, indicating successful processing.
	// This should be called after the message has been successfully processed.
	Ack() error

	// Nack negatively acknowledges the message, indicating failed processing.
	// It allows the message to be requeued or sent to a dead-letter queue.
	Nack(requeue bool) error

	// Reject rejects the message and optionally requeues it.
	// Similar to Nack but with different semantics in some implementations.
	Reject(requeue bool) error

	// DeliveryTag returns a unique identifier for this message delivery.
	// This can be used for tracking and debugging purposes.
	DeliveryTag() uint64

	// Timestamp returns when the message was originally published.
	Timestamp() time.Time

	// Headers returns any headers associated with the message.
	Headers() map[string]interface{}
}

// DeliveryChannel is a channel that receives Message instances.
// It's used by Consumer implementations to deliver messages to subscribers.
type DeliveryChannel <-chan Message
