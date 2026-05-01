package main

import (
	"context"
	"fmt"
	"time"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/eventbus"
	"github.com/reitmas32/rkit/infrastructure/eventbus/rabbit"
)

// UserCreatedEvent represents a user created event.
type UserCreatedEvent struct {
	EventName       string            `json:"name"`
	EventVersion    string            `json:"version"`
	EventOccurredAt time.Time         `json:"occurred_at"`
	EventPayload    UserPayload       `json:"payload"`
	EventMetadata   eventbus.Metadata `json:"metadata"`
}

type UserPayload struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (e *UserCreatedEvent) Name() string                { return e.EventName }
func (e *UserCreatedEvent) Version() string             { return e.EventVersion }
func (e *UserCreatedEvent) OccurredAt() time.Time       { return e.EventOccurredAt }
func (e *UserCreatedEvent) Payload() any                { return e.EventPayload }
func (e *UserCreatedEvent) Metadata() eventbus.Metadata { return e.EventMetadata }

func main() {
	fmt.Println("=== RabbitMQ Event Publisher ===")

	// Configuration
	amqpURL := "amqp://guest:guest@localhost:5672/"

	config := rabbit.DefaultConfig(amqpURL)
	config.ExchangeName = "events"
	config.ExchangeType = "topic"

	// Event factory - creates event instances based on event name
	eventFactory := func(eventName string) eventbus.Event {
		switch eventName {
		case "user.created":
			return &UserCreatedEvent{}
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

	ctx := customctx.New(context.Background())

	// Publish events
	fmt.Println("1. Publishing UserCreatedEvent...")

	event1 := &UserCreatedEvent{
		EventName:       "user.created",
		EventVersion:    "1.0",
		EventOccurredAt: time.Now(),
		EventPayload: UserPayload{
			UserID:   "123",
			Username: "john_doe",
			Email:    "john@example.com",
		},
		EventMetadata: eventbus.Metadata{
			"source":   "api",
			"trace_id": "abc-123",
		},
	}

	if err := eb.Publish(ctx, event1); err != nil {
		fmt.Printf("Error publishing event: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Published: user.created (UserID: %s, Username: %s)\n",
		event1.EventPayload.UserID, event1.EventPayload.Username)

	// Publish another event
	time.Sleep(500 * time.Millisecond)

	event2 := &UserCreatedEvent{
		EventName:       "user.created",
		EventVersion:    "1.0",
		EventOccurredAt: time.Now(),
		EventPayload: UserPayload{
			UserID:   "456",
			Username: "jane_smith",
			Email:    "jane@example.com",
		},
		EventMetadata: eventbus.Metadata{
			"source":   "api",
			"trace_id": "def-456",
		},
	}

	if err := eb.Publish(ctx, event2); err != nil {
		fmt.Printf("Error publishing event: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Published: user.created (UserID: %s, Username: %s)\n",
		event2.EventPayload.UserID, event2.EventPayload.Username)

	fmt.Println("\n✓ Events published successfully!")
}
