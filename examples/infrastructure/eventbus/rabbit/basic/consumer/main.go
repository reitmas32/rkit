package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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
	fmt.Println("=== RabbitMQ Event Consumer ===")

	// Configuration
	amqpURL := "amqps://admin:supersecret@rabbit.konectus.tech:443/"

	config := rabbit.DefaultConfig(amqpURL)
	config.ExchangeName = "events"
	config.ExchangeType = "topic"
	config.QueuePrefix = "consumer" // Optional: prefix for queue names

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
	fmt.Printf("  Queue Prefix: %s\n", config.QueuePrefix)
	fmt.Println()

	ctx := customctx.New(context.Background())

	// Create event instance for subscription
	eventTemplate := &UserCreatedEvent{
		EventName:       "user.created",
		EventVersion:    "1.0",
		EventOccurredAt: time.Now(),
		EventPayload:    UserPayload{},
		EventMetadata:   eventbus.Metadata{},
	}

	// Subscribe to events
	fmt.Println("Subscribing to user.created events...")
	consumeResult := eb.Consume(ctx, eventTemplate)
	if !consumeResult.IsOk() {
		fmt.Printf("Error subscribing: %v\n", consumeResult.Error())
		return
	}

	deliveryChan := consumeResult.Value()
	fmt.Println("✓ Subscribed successfully")
	fmt.Println("\nWaiting for messages... (Press Ctrl+C to stop)")
	fmt.Println()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	messageCount := 0

	// Consume messages
	for {
		select {
		case msg, ok := <-deliveryChan:
			if !ok {
				fmt.Println("Channel closed")
				return
			}

			messageCount++
			fmt.Printf("--- Message #%d ---\n", messageCount)

			event := msg.Event()
			if userEvent, ok := event.(*UserCreatedEvent); ok {
				payload := userEvent.Payload().(UserPayload)
				fmt.Printf("Event: %s\n", event.Name())
				fmt.Printf("Version: %s\n", event.Version())
				fmt.Printf("UserID: %s\n", payload.UserID)
				fmt.Printf("Username: %s\n", payload.Username)
				fmt.Printf("Email: %s\n", payload.Email)
				fmt.Printf("Delivery Tag: %d\n", msg.DeliveryTag())
				fmt.Printf("Timestamp: %s\n", msg.Timestamp().Format(time.RFC3339))
				fmt.Printf("Headers: %v\n", msg.Headers())
				fmt.Printf("Metadata: %v\n", event.Metadata())
			} else {
				fmt.Printf("Event: %s (unknown type)\n", event.Name())
			}

			// Acknowledge message
			if err := msg.Ack(); err != nil {
				fmt.Printf("Error acknowledging message: %v\n", err)
			} else {
				fmt.Println("✓ Message acknowledged")
			}

			fmt.Println()

		case <-sigChan:
			fmt.Println("\n✓ Shutting down gracefully...")
			fmt.Printf("Total messages processed: %d\n", messageCount)
			return
		}
	}
}
