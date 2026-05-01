package inmemory

import (
	"sync"
	"time"

	"github.com/reitmas32/rkit/core/eventbus"
)

// message implements the eventbus.Message interface for in-memory event bus.
type message struct {
	event       eventbus.Event
	deliveryTag uint64
	timestamp   time.Time
	headers     map[string]interface{}
	acked       bool
	nacked      bool
	rejected    bool
	mu          sync.RWMutex
}

// newMessage creates a new in-memory message.
func newMessage(event eventbus.Event, deliveryTag uint64, headers map[string]interface{}) *message {
	return &message{
		event:       event,
		deliveryTag: deliveryTag,
		timestamp:   event.OccurredAt(),
		headers:     headers,
	}
}

// Event returns the event associated with this message.
func (m *message) Event() eventbus.Event {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.event
}

// Ack acknowledges the message, indicating successful processing.
// For in-memory implementation, this is a no-op but marks the message as acknowledged.
func (m *message) Ack() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.acked || m.nacked || m.rejected {
		return nil // Already processed
	}

	m.acked = true
	return nil
}

// Nack negatively acknowledges the message.
// For in-memory implementation, this marks the message as nacked.
// If requeue is true, the message can be redelivered (handled by the event bus).
func (m *message) Nack(requeue bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.acked || m.nacked || m.rejected {
		return nil // Already processed
	}

	m.nacked = true
	return nil
}

// Reject rejects the message.
// For in-memory implementation, this marks the message as rejected.
// If requeue is true, the message can be redelivered.
func (m *message) Reject(requeue bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.acked || m.nacked || m.rejected {
		return nil // Already processed
	}

	m.rejected = true
	return nil
}

// DeliveryTag returns a unique identifier for this message delivery.
func (m *message) DeliveryTag() uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.deliveryTag
}

// Timestamp returns when the message was originally published.
func (m *message) Timestamp() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.timestamp
}

// Headers returns any headers associated with the message.
func (m *message) Headers() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent external modification
	if m.headers == nil {
		return nil
	}

	headers := make(map[string]interface{}, len(m.headers))
	for k, v := range m.headers {
		headers[k] = v
	}
	return headers
}

// isProcessed returns true if the message has been acked, nacked, or rejected.
func (m *message) isProcessed() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.acked || m.nacked || m.rejected
}
