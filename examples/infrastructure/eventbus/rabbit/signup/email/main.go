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
	fmt.Println("=== RabbitMQ Event Consumer (Email Service) ===")

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

	// Create event template for subscription
	eventTemplate := &UserSignupEvent{
		EventName:       "user.signup",
		EventVersion:    "1.0",
		EventOccurredAt: time.Now(),
		EventPayload:    UserSignupPayload{},
		EventMetadata:   eventbus.Metadata{},
	}

	// Subscribe with custom queue name (each consumer has its own queue)
	queueName := "email-service.user.signup"
	fmt.Printf("Subscribing to user.signup events (queue: %s)...\n", queueName)
	consumeResult := eb.ConsumeWithQueue(ctx, eventTemplate, queueName)
	if !consumeResult.IsOk() {
		fmt.Printf("Error subscribing: %v\n", consumeResult.Error())
		return
	}

	deliveryChan := consumeResult.Value()
	fmt.Println("✓ Subscribed successfully")
	fmt.Println("\nWaiting for messages to send welcome emails... (Press Ctrl+C to stop)")
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
			event := msg.Event()
			if signupEvent, ok := event.(*UserSignupEvent); ok {
				payload := signupEvent.Payload().(UserSignupPayload)
				fmt.Printf("--- Email Service Message #%d ---\n", messageCount)
				fmt.Printf("Event: %s\n", event.Name())
				fmt.Printf("UserID: %s\n", payload.UserID)
				fmt.Printf("Email: %s\n", payload.Email)
				fmt.Printf("Username: %s\n", payload.Username)
				fmt.Println()
				fmt.Printf("  → Sending welcome email to: %s\n", payload.Email)
				fmt.Println("  ✓ Welcome email sent!")
			}

			// Acknowledge message
			if err := msg.Ack(); err != nil {
				fmt.Printf("Error acknowledging message: %v\n", err)
			}

			fmt.Println()

		case <-sigChan:
			fmt.Println("\n✓ Shutting down gracefully...")
			fmt.Printf("Total welcome emails sent: %d\n", messageCount)
			return
		}
	}
}
