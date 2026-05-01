package rabbit

import (
	"encoding/json"
	"time"

	"github.com/reitmas32/rkit/core/eventbus"
	"github.com/rabbitmq/amqp091-go"
)

// message wraps amqp091.Delivery to implement eventbus.Message interface.
type message struct {
	delivery *amqp091.Delivery
	event    eventbus.Event
}

// newMessage creates a new message wrapper from an AMQP delivery.
// It deserializes the event from the delivery body.
func newMessage(delivery *amqp091.Delivery, event eventbus.Event) *message {
	return &message{
		delivery: delivery,
		event:    event,
	}
}

// Event returns the event associated with this message.
func (m *message) Event() eventbus.Event {
	return m.event
}

// Ack acknowledges the message, indicating successful processing.
func (m *message) Ack() error {
	return m.delivery.Ack(false)
}

// Nack negatively acknowledges the message, indicating failed processing.
func (m *message) Nack(requeue bool) error {
	return m.delivery.Nack(false, requeue)
}

// Reject rejects the message and optionally requeues it.
func (m *message) Reject(requeue bool) error {
	return m.delivery.Reject(requeue)
}

// DeliveryTag returns a unique identifier for this message delivery.
func (m *message) DeliveryTag() uint64 {
	return m.delivery.DeliveryTag
}

// Timestamp returns when the message was originally published.
func (m *message) Timestamp() time.Time {
	if m.delivery.Timestamp.IsZero() {
		return time.Now() // Fallback if timestamp is not set
	}
	return m.delivery.Timestamp
}

// Headers returns any headers associated with the message.
func (m *message) Headers() map[string]interface{} {
	if m.delivery.Headers == nil {
		return nil
	}

	// Convert amqp.Table to map[string]interface{}
	headers := make(map[string]interface{}, len(m.delivery.Headers))
	for k, v := range m.delivery.Headers {
		headers[k] = v
	}
	return headers
}

// deserializeEvent deserializes an event from AMQP delivery body.
// It expects the body to be JSON and tries to unmarshal it into the provided event type.
// This is a helper function that can be used by the EventBus implementation.
func deserializeEvent(delivery *amqp091.Delivery, eventFactory func() eventbus.Event) (eventbus.Event, error) {
	event := eventFactory()

	if err := json.Unmarshal(delivery.Body, event); err != nil {
		return nil, err
	}

	return event, nil
}
