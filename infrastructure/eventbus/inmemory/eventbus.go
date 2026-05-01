package inmemory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/eventbus"
	"github.com/reitmas32/rkit/core/result"
)

// EventBus is an in-memory implementation of eventbus.EventBus.
// It stores events in memory and delivers them to subscribers.
// This implementation is useful for testing and single-process applications.
type EventBus struct {
	// subscribers stores channels for each event type.
	// Key: event name, Value: slice of channels
	subscribers map[string][]chan eventbus.Message

	// publishedEvents stores all published events for debugging/testing.
	publishedEvents []eventbus.Event

	// delayedEvents stores events that should be published with delay.
	delayedEvents []delayedEvent

	// nextDeliveryTag is used to generate unique delivery tags.
	nextDeliveryTag uint64

	mu sync.RWMutex

	// closed indicates if the event bus has been closed.
	closed bool

	// done is used to signal goroutines to stop.
	done chan struct{}
	wg   sync.WaitGroup
}

// delayedEvent represents an event that should be published after a delay.
type delayedEvent struct {
	event   eventbus.Event
	delay   time.Duration
	publish time.Time
}

// NewEventBus creates a new in-memory event bus.
func NewEventBus() *EventBus {
	eb := &EventBus{
		subscribers:     make(map[string][]chan eventbus.Message),
		publishedEvents: make([]eventbus.Event, 0),
		delayedEvents:   make([]delayedEvent, 0),
		nextDeliveryTag: 1,
		done:            make(chan struct{}),
	}

	// Start goroutine to handle delayed events
	eb.wg.Add(1)
	go eb.processDelayedEvents()

	return eb
}

// Publish publishes an event to all subscribers.
func (eb *EventBus) Publish(ctx *customctx.CustomContext, event eventbus.Event) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.closed {
		return fmt.Errorf("event bus is closed")
	}

	// Store published event
	eb.publishedEvents = append(eb.publishedEvents, event)

	// Get subscribers for this event type
	subscribers := eb.subscribers[event.Name()]

	// Create message for each subscriber
	for _, subChan := range subscribers {
		deliveryTag := eb.nextDeliveryTag
		eb.nextDeliveryTag++

		// Convert Metadata (map[string]string) to map[string]interface{}
		metadata := make(map[string]interface{})
		for k, v := range event.Metadata() {
			metadata[k] = v
		}
		msg := newMessage(event, deliveryTag, metadata)

		// Send message to subscriber (non-blocking)
		select {
		case subChan <- msg:
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Skip if channel is full (subscriber is not ready)
			// In a real implementation, you might want to buffer or handle this differently
		}
	}

	return nil
}

// PublishWithDelay publishes an event after the specified delay.
func (eb *EventBus) PublishWithDelay(ctx *customctx.CustomContext, event eventbus.Event, delay time.Duration) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.closed {
		return fmt.Errorf("event bus is closed")
	}

	eb.delayedEvents = append(eb.delayedEvents, delayedEvent{
		event:   event,
		delay:   delay,
		publish: time.Now().Add(delay),
	})

	return nil
}

// Consume creates a channel that receives messages for the specified event type.
func (eb *EventBus) Consume(ctx *customctx.CustomContext, event eventbus.Event) result.Result[eventbus.DeliveryChannel] {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.closed {
		return result.NewErrResult[eventbus.DeliveryChannel](
			eventbus.ErrEventFailedToSubscribe,
		)
	}

	eventName := event.Name()

	// Create channel for this subscriber
	msgChan := make(chan eventbus.Message, 100) // Buffered channel

	// Add subscriber
	eb.subscribers[eventName] = append(eb.subscribers[eventName], msgChan)

	// Clean up when context is done
	go func() {
		<-ctx.Done()
		eb.mu.Lock()
		defer eb.mu.Unlock()

		// Remove this channel from subscribers
		subscribers := eb.subscribers[eventName]
		for i, sub := range subscribers {
			if sub == msgChan {
				eb.subscribers[eventName] = append(subscribers[:i], subscribers[i+1:]...)
				close(msgChan)
				break
			}
		}
	}()

	return result.Ok(eventbus.DeliveryChannel(msgChan))
}

// Close closes the event bus and cleans up resources.
func (eb *EventBus) Close() error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.closed {
		return nil
	}

	eb.closed = true
	close(eb.done)

	// Close all subscriber channels
	for _, subscribers := range eb.subscribers {
		for _, subChan := range subscribers {
			close(subChan)
		}
	}
	eb.subscribers = make(map[string][]chan eventbus.Message)

	// Wait for goroutines to finish
	eb.mu.Unlock()
	eb.wg.Wait()
	eb.mu.Lock()

	return nil
}

// processDelayedEvents processes delayed events in a background goroutine.
func (eb *EventBus) processDelayedEvents() {
	defer eb.wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-eb.done:
			return
		case <-ticker.C:
			eb.mu.Lock()
			now := time.Now()
			readyEvents := make([]eventbus.Event, 0)
			remainingDelayed := make([]delayedEvent, 0)

			for _, delayed := range eb.delayedEvents {
				if now.After(delayed.publish) || now.Equal(delayed.publish) {
					readyEvents = append(readyEvents, delayed.event)
				} else {
					remainingDelayed = append(remainingDelayed, delayed)
				}
			}

			eb.delayedEvents = remainingDelayed
			eb.mu.Unlock()

			// Publish ready events
			for _, event := range readyEvents {
				_ = eb.Publish(customctx.New(context.Background()), event)
			}
		}
	}
}
