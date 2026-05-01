package main

import (
	"context"
	"fmt"
	"time"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/eventbus"
	"github.com/reitmas32/rkit/infrastructure/eventbus/rabbit"
)

// UserSignupEvent represents a user signup event.
type UserSignupEvent struct {
	EventName       string            `json:"name"`
	EventVersion    string            `json:"version"`
	EventOccurredAt time.Time         `json:"occurred_at"`
	EventPayload    UserSignupPayload `json:"payload"`
	EventMetadata   eventbus.Metadata `json:"metadata"`
}

type UserSignupPayload struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (e *UserSignupEvent) Name() string                { return e.EventName }
func (e *UserSignupEvent) Version() string             { return e.EventVersion }
func (e *UserSignupEvent) OccurredAt() time.Time       { return e.EventOccurredAt }
func (e *UserSignupEvent) Payload() any                { return e.EventPayload }
func (e *UserSignupEvent) Metadata() eventbus.Metadata { return e.EventMetadata }

func main() {
	fmt.Println("=== RabbitMQ Event Publisher (User Signup) ===")

	// Configuration
	amqpURL := "amqp://guest:guest@localhost:5672/"

	config := rabbit.DefaultConfig(amqpURL)
	config.ExchangeName = "events"
	config.ExchangeType = "topic"

	// Event factory
	eventFactory := func(eventName string) eventbus.Event {
		switch eventName {
		case "user.signup":
			return &UserSignupEvent{}
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

	// Publish user signup events
	fmt.Println("Publishing user.signup events...")
	fmt.Println("These events will be received by ALL queues bound to 'user.signup' routing key")
	fmt.Println()

	events := []struct {
		userID   string
		username string
		email    string
	}{
		{"123", "john_doe", "john@example.com"},
		{"456", "jane_smith", "jane@example.com"},
		{"789", "bob_wilson", "bob@example.com"},
	}

	for i, e := range events {
		event := &UserSignupEvent{
			EventName:       "user.signup",
			EventVersion:    "1.0",
			EventOccurredAt: time.Now(),
			EventPayload: UserSignupPayload{
				UserID:   e.userID,
				Username: e.username,
				Email:    e.email,
			},
			EventMetadata: eventbus.Metadata{
				"source":   "api",
				"trace_id": fmt.Sprintf("trace-%d", i+1),
			},
		}

		if err := eb.Publish(ctx, event); err != nil {
			fmt.Printf("Error publishing event: %v\n", err)
			return
		}

		fmt.Printf("✓ Published: user.signup (UserID: %s, Username: %s, Email: %s)\n",
			e.userID, e.username, e.email)

		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n✓ All events published successfully!")
	fmt.Println("\nBoth email-service and push-service queues should receive all events.")
}
