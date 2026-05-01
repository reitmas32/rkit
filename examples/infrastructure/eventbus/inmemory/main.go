package main

import (
	"context"
	"fmt"
	"time"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/eventbus"
	"github.com/reitmas32/rkit/infrastructure/eventbus/inmemory"
)

// UserCreatedEvent represents a user created event.
type UserCreatedEvent struct {
	name       string
	version    string
	occurredAt time.Time
	payload    UserPayload
	metadata   eventbus.Metadata
}

type UserPayload struct {
	UserID   string
	Username string
	Email    string
}

func (e *UserCreatedEvent) Name() string                { return e.name }
func (e *UserCreatedEvent) Version() string             { return e.version }
func (e *UserCreatedEvent) OccurredAt() time.Time       { return e.occurredAt }
func (e *UserCreatedEvent) Payload() any                { return e.payload }
func (e *UserCreatedEvent) Metadata() eventbus.Metadata { return e.metadata }

// OrderPlacedEvent represents an order placed event.
type OrderPlacedEvent struct {
	name       string
	version    string
	occurredAt time.Time
	payload    OrderPayload
	metadata   eventbus.Metadata
}

type OrderPayload struct {
	OrderID string
	UserID  string
	Amount  float64
	Items   []string
}

func (e *OrderPlacedEvent) Name() string                { return e.name }
func (e *OrderPlacedEvent) Version() string             { return e.version }
func (e *OrderPlacedEvent) OccurredAt() time.Time       { return e.occurredAt }
func (e *OrderPlacedEvent) Payload() any                { return e.payload }
func (e *OrderPlacedEvent) Metadata() eventbus.Metadata { return e.metadata }

func main() {
	fmt.Println("=== EventBus In-Memory Example ===")

	// Create event bus
	eb := inmemory.NewEventBus()
	defer eb.Close()

	ctx := context.Background()

	// Example 1: Basic publish and consume
	fmt.Println("1. Basic Publish and Consume")
	exampleBasicPublishConsume(ctx, eb)
	fmt.Println()

	// Example 2: Multiple subscribers
	fmt.Println("2. Multiple Subscribers")
	exampleMultipleSubscribers(ctx, eb)
	fmt.Println()

	// Example 3: Publish with delay
	fmt.Println("3. Publish with Delay")
	examplePublishWithDelay(ctx, eb)
	fmt.Println()

	// Example 4: Message acknowledgment
	fmt.Println("4. Message Acknowledgment")
	exampleMessageAcknowledgment(ctx, eb)
	fmt.Println()
}

func exampleBasicPublishConsume(ctx context.Context, eb eventbus.EventBus) {
	// Create event
	event := &UserCreatedEvent{
		name:       "user.created",
		version:    "1.0",
		occurredAt: time.Now(),
		payload: UserPayload{
			UserID:   "123",
			Username: "john_doe",
			Email:    "john@example.com",
		},
		metadata: eventbus.Metadata{
			"source":   "api",
			"trace_id": "abc-123",
		},
	}

	cc := customctx.NewCustomContext(ctx)

	// Subscribe to events
	consumeResult := eb.Consume(cc, event)
	if !consumeResult.IsOk() {
		fmt.Printf("Error subscribing: %v\n", consumeResult.Error())
		return
	}

	deliveryChan := consumeResult.Value()

	// Publish event
	if err := eb.Publish(cc, event); err != nil {
		fmt.Printf("Error publishing: %v\n", err)
		return
	}

	fmt.Println("  Event published: user.created")

	// Receive message
	select {
	case msg := <-deliveryChan:
		if msg == nil {
			fmt.Println("  Channel closed")
			return
		}
		receivedEvent := msg.Event()
		if userEvent, ok := receivedEvent.(*UserCreatedEvent); ok {
			payload := userEvent.Payload().(UserPayload)
			fmt.Printf("  Event received: %s (UserID: %s, Username: %s)\n",
				receivedEvent.Name(), payload.UserID, payload.Username)
			fmt.Printf("  Delivery Tag: %d\n", msg.DeliveryTag())
			fmt.Printf("  Timestamp: %s\n", msg.Timestamp().Format(time.RFC3339))
			fmt.Printf("  Headers: %v\n", msg.Headers())

			// Acknowledge message
			if err := msg.Ack(); err != nil {
				fmt.Printf("  Error acknowledging: %v\n", err)
			} else {
				fmt.Println("  Message acknowledged")
			}
		}
	case <-time.After(2 * time.Second):
		fmt.Println("  Timeout waiting for message")
	}
}

func exampleMultipleSubscribers(ctx context.Context, eb eventbus.EventBus) {
	// Create event
	event := &OrderPlacedEvent{
		name:       "order.placed",
		version:    "1.0",
		occurredAt: time.Now(),
		payload: OrderPayload{
			OrderID: "order-456",
			UserID:  "123",
			Amount:  99.99,
			Items:   []string{"item1", "item2"},
		},
		metadata: eventbus.Metadata{
			"source": "web",
		},
	}

	cc := customctx.NewCustomContext(ctx)

	// Create multiple subscribers
	subscriber1 := eb.Consume(cc, event)
	subscriber2 := eb.Consume(cc, event)

	if !subscriber1.IsOk() || !subscriber2.IsOk() {
		fmt.Println("  Error creating subscribers")
		return
	}

	deliveryChan1 := subscriber1.Value()
	deliveryChan2 := subscriber2.Value()

	// Publish event
	if err := eb.Publish(cc, event); err != nil {
		fmt.Printf("  Error publishing: %v\n", err)
		return
	}

	fmt.Println("  Event published: order.placed")

	// Receive from first subscriber
	select {
	case msg := <-deliveryChan1:
		if msg != nil {
			fmt.Println("  Subscriber 1 received message")
			msg.Ack()
		}
	case <-time.After(1 * time.Second):
		fmt.Println("  Subscriber 1 timeout")
	}

	// Receive from second subscriber
	select {
	case msg := <-deliveryChan2:
		if msg != nil {
			fmt.Println("  Subscriber 2 received message")
			msg.Ack()
		}
	case <-time.After(1 * time.Second):
		fmt.Println("  Subscriber 2 timeout")
	}
}

func examplePublishWithDelay(ctx context.Context, eb eventbus.EventBus) {
	event := &UserCreatedEvent{
		name:       "user.created",
		version:    "1.0",
		occurredAt: time.Now(),
		payload: UserPayload{
			UserID:   "789",
			Username: "delayed_user",
			Email:    "delayed@example.com",
		},
		metadata: eventbus.Metadata{
			"source": "batch",
		},
	}

	cc := customctx.NewCustomContext(ctx)

	// Subscribe
	consumeResult := eb.Consume(cc, event)
	if !consumeResult.IsOk() {
		fmt.Println("  Error subscribing")
		return
	}

	deliveryChan := consumeResult.Value()

	// Publish with 1 second delay
	fmt.Println("  Publishing event with 1 second delay...")
	if err := eb.PublishWithDelay(cc, event, 1*time.Second); err != nil {
		fmt.Printf("  Error publishing: %v\n", err)
		return
	}

	fmt.Println("  Waiting for delayed message...")
	start := time.Now()

	select {
	case msg := <-deliveryChan:
		elapsed := time.Since(start)
		if msg != nil {
			fmt.Printf("  Event received after %v\n", elapsed)
			receivedEvent := msg.Event()
			if userEvent, ok := receivedEvent.(*UserCreatedEvent); ok {
				payload := userEvent.Payload().(UserPayload)
				fmt.Printf("  UserID: %s\n", payload.UserID)
			}
			msg.Ack()
		}
	case <-time.After(3 * time.Second):
		fmt.Println("  Timeout waiting for delayed message")
	}
}

func exampleMessageAcknowledgment(ctx context.Context, eb eventbus.EventBus) {
	event := &OrderPlacedEvent{
		name:       "order.placed",
		version:    "1.0",
		occurredAt: time.Now(),
		payload: OrderPayload{
			OrderID: "order-789",
			UserID:  "456",
			Amount:  149.99,
			Items:   []string{"item3"},
		},
		metadata: eventbus.Metadata{
			"source": "api",
		},
	}

	cc := customctx.NewCustomContext(ctx)

	// Subscribe
	consumeResult := eb.Consume(cc, event)
	if !consumeResult.IsOk() {
		fmt.Println("  Error subscribing")
		return
	}

	deliveryChan := consumeResult.Value()

	// Publish
	if err := eb.Publish(cc, event); err != nil {
		fmt.Printf("  Error publishing: %v\n", err)
		return
	}

	// Receive message
	select {
	case msg := <-deliveryChan:
		if msg == nil {
			return
		}

		fmt.Printf("  Message received (Tag: %d)\n", msg.DeliveryTag())

		// Try to ack multiple times (should only work once)
		if err := msg.Ack(); err != nil {
			fmt.Printf("  Error on first ack: %v\n", err)
		} else {
			fmt.Println("  First Ack: OK")
		}

		if err := msg.Ack(); err != nil {
			fmt.Printf("  Error on second ack: %v\n", err)
		} else {
			fmt.Println("  Second Ack: OK (no-op, already acknowledged)")
		}

		// Try to nack after ack (should be no-op)
		if err := msg.Nack(true); err != nil {
			fmt.Printf("  Error on nack after ack: %v\n", err)
		} else {
			fmt.Println("  Nack after Ack: OK (no-op, already acknowledged)")
		}

	case <-time.After(1 * time.Second):
		fmt.Println("  Timeout waiting for message")
	}
}
