package main

import (
	"context"
	"fmt"
	"time"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/eventbus"
	"github.com/reitmas32/rkit/infrastructure/eventbus/rabbit"
)

// OrderCreatedEvent represents an order created event.
type OrderCreatedEvent struct {
	EventName       string            `json:"name"`
	EventVersion    string            `json:"version"`
	EventOccurredAt time.Time         `json:"occurred_at"`
	EventPayload    OrderPayload      `json:"payload"`
	EventMetadata   eventbus.Metadata `json:"metadata"`
}

type OrderPayload struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Items     []string  `json:"items"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *OrderCreatedEvent) Name() string                { return e.EventName }
func (e *OrderCreatedEvent) Version() string             { return e.EventVersion }
func (e *OrderCreatedEvent) OccurredAt() time.Time       { return e.EventOccurredAt }
func (e *OrderCreatedEvent) Payload() any                { return e.EventPayload }
func (e *OrderCreatedEvent) Metadata() eventbus.Metadata { return e.EventMetadata }

func main() {
	fmt.Println("=== RabbitMQ Event Publisher (Delayed Messages) ===")

	// Configuration
	amqpURL := "amqp://guest:guest@localhost:5672/"

	config := rabbit.DefaultConfig(amqpURL)
	config.ExchangeName = "events"
	config.ExchangeType = "topic"

	// Event factory - creates event instances based on event name
	eventFactory := func(eventName string) eventbus.Event {
		switch eventName {
		case "order.created":
			return &OrderCreatedEvent{}
		default:
			return nil
		}
	}

	// Create event bus
	eb, err := rabbit.NewEventBus(config, eventFactory)
	if err != nil {
		fmt.Printf("Error creating event bus: %v\n", err)
		return
	}
	defer eb.Close()

	fmt.Println("✓ Connected to RabbitMQ")
	fmt.Printf("  Exchange: %s (%s)\n", config.ExchangeName, config.ExchangeType)
	fmt.Println()
	fmt.Println("Note: This example requires the 'rabbitmq-delayed-message-exchange' plugin")
	fmt.Println("      to be installed and enabled in RabbitMQ.")
	fmt.Println()

	ctx := customctx.New(context.Background())

	// Publish immediate event
	fmt.Println("1. Publishing immediate event (order.created)...")

	event1 := &OrderCreatedEvent{
		EventName:       "order.created",
		EventVersion:    "1.0",
		EventOccurredAt: time.Now(),
		EventPayload: OrderPayload{
			OrderID:   "ORD-001",
			UserID:    "user-123",
			Amount:    99.99,
			Items:     []string{"Product A", "Product B"},
			CreatedAt: time.Now(),
		},
		EventMetadata: eventbus.Metadata{
			"source":   "api",
			"trace_id": "trace-immediate",
		},
	}

	if err := eb.Publish(ctx, event1); err != nil {
		fmt.Printf("Error publishing event: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Published immediately: order.created (OrderID: %s, Amount: $%.2f)\n",
		event1.EventPayload.OrderID, event1.EventPayload.Amount)
	fmt.Println()

	// Publish delayed event (5 seconds)
	fmt.Println("2. Publishing delayed event (5 seconds delay)...")

	event2 := &OrderCreatedEvent{
		EventName:       "order.created",
		EventVersion:    "1.0",
		EventOccurredAt: time.Now(),
		EventPayload: OrderPayload{
			OrderID:   "ORD-002",
			UserID:    "user-456",
			Amount:    149.99,
			Items:     []string{"Product C", "Product D"},
			CreatedAt: time.Now(),
		},
		EventMetadata: eventbus.Metadata{
			"source":   "api",
			"trace_id": "trace-delayed-5s",
		},
	}

	delay5s := 5 * time.Second
	if err := eb.PublishWithDelay(ctx, event2, delay5s); err != nil {
		fmt.Printf("Error publishing delayed event: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Published with 5s delay: order.created (OrderID: %s, Amount: $%.2f)\n",
		event2.EventPayload.OrderID, event2.EventPayload.Amount)
	fmt.Printf("  → This event will be delivered in %v\n", delay5s)
	fmt.Println()

	// Publish delayed event (10 seconds)
	fmt.Println("3. Publishing delayed event (10 seconds delay)...")

	event3 := &OrderCreatedEvent{
		EventName:       "order.created",
		EventVersion:    "1.0",
		EventOccurredAt: time.Now(),
		EventPayload: OrderPayload{
			OrderID:   "ORD-003",
			UserID:    "user-789",
			Amount:    199.99,
			Items:     []string{"Product E"},
			CreatedAt: time.Now(),
		},
		EventMetadata: eventbus.Metadata{
			"source":   "api",
			"trace_id": "trace-delayed-10s",
		},
	}

	delay10s := 10 * time.Second
	if err := eb.PublishWithDelay(ctx, event3, delay10s); err != nil {
		fmt.Printf("Error publishing delayed event: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Published with 10s delay: order.created (OrderID: %s, Amount: $%.2f)\n",
		event3.EventPayload.OrderID, event3.EventPayload.Amount)
	fmt.Printf("  → This event will be delivered in %v\n", delay10s)
	fmt.Println()

	// Publish delayed event (1 minute)
	fmt.Println("4. Publishing delayed event (1 minute delay)...")

	event4 := &OrderCreatedEvent{
		EventName:       "order.created",
		EventVersion:    "1.0",
		EventOccurredAt: time.Now(),
		EventPayload: OrderPayload{
			OrderID:   "ORD-004",
			UserID:    "user-321",
			Amount:    249.99,
			Items:     []string{"Product F", "Product G", "Product H"},
			CreatedAt: time.Now(),
		},
		EventMetadata: eventbus.Metadata{
			"source":   "api",
			"trace_id": "trace-delayed-1m",
		},
	}

	delay1m := 1 * time.Minute
	if err := eb.PublishWithDelay(ctx, event4, delay1m); err != nil {
		fmt.Printf("Error publishing delayed event: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Published with 1m delay: order.created (OrderID: %s, Amount: $%.2f)\n",
		event4.EventPayload.OrderID, event4.EventPayload.Amount)
	fmt.Printf("  → This event will be delivered in %v\n", delay1m)
	fmt.Println()

	fmt.Println("✓ All events published successfully!")
	fmt.Println("\nSummary:")
	fmt.Println("  - 1 immediate event (delivered immediately)")
	fmt.Println("  - 1 delayed event (5 seconds)")
	fmt.Println("  - 1 delayed event (10 seconds)")
	fmt.Println("  - 1 delayed event (1 minute)")
	fmt.Println("\nStart a consumer to see the events being delivered at their scheduled times.")
}
