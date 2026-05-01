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
	fmt.Println("=== RabbitMQ Event Consumer (Delayed Messages) ===")

	// Configuration
	amqpURL := "amqp://guest:guest@localhost:5672/"

	config := rabbit.DefaultConfig(amqpURL)
	config.ExchangeName = "events"
	config.ExchangeType = "topic"
	config.QueuePrefix = "delayed-consumer" // Optional: prefix for queue names

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
	fmt.Printf("  Queue Prefix: %s\n", config.QueuePrefix)
	fmt.Println()
	fmt.Println("This consumer will process events that were published with delay.")
	fmt.Println("You should see events arriving at their scheduled delivery times.")
	fmt.Println()

	ctx := customctx.New(context.Background())

	// Create event instance for subscription
	eventTemplate := &OrderCreatedEvent{
		EventName:       "order.created",
		EventVersion:    "1.0",
		EventOccurredAt: time.Now(),
		EventPayload:    OrderPayload{},
		EventMetadata:   eventbus.Metadata{},
	}

	// Subscribe to events
	fmt.Println("Subscribing to order.created events...")
	consumeResult := eb.Consume(ctx, eventTemplate)
	if !consumeResult.IsOk() {
		fmt.Printf("Error subscribing: %v\n", consumeResult.Error())
		return
	}

	deliveryChan := consumeResult.Value()
	fmt.Println("✓ Subscribed successfully")
	fmt.Println("\nWaiting for messages... (Press Ctrl+C to stop)")
	fmt.Println()
	fmt.Println("Note: Events published with delay will arrive at their scheduled times.")
	fmt.Println("      You should see immediate events first, then delayed ones later.")
	fmt.Println()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	messageCount := 0
	startTime := time.Now()

	// Consume messages
	for {
		select {
		case msg, ok := <-deliveryChan:
			if !ok {
				fmt.Println("Channel closed")
				return
			}

			messageCount++
			receiveTime := time.Now()
			elapsedTime := receiveTime.Sub(startTime)

			fmt.Printf("--- Message #%d (Received at +%v) ---\n", messageCount, elapsedTime.Round(time.Second))

			event := msg.Event()
			if orderEvent, ok := event.(*OrderCreatedEvent); ok {
				payload := orderEvent.Payload().(OrderPayload)

				fmt.Printf("Event: %s\n", event.Name())
				fmt.Printf("Version: %s\n", event.Version())
				fmt.Printf("OrderID: %s\n", payload.OrderID)
				fmt.Printf("UserID: %s\n", payload.UserID)
				fmt.Printf("Amount: $%.2f\n", payload.Amount)
				fmt.Printf("Items: %v\n", payload.Items)
				fmt.Printf("Created At: %s\n", payload.CreatedAt.Format(time.RFC3339))
				fmt.Printf("Received At: %s\n", receiveTime.Format(time.RFC3339))
				fmt.Printf("Delivery Tag: %d\n", msg.DeliveryTag())
				fmt.Printf("Timestamp: %s\n", msg.Timestamp().Format(time.RFC3339))

				// Check if there's delay information in headers
				headers := msg.Headers()
				if delayMs, ok := headers["x-delay"]; ok {
					fmt.Printf("Delay (from header): %v ms\n", delayMs)
				}

				fmt.Printf("Metadata: %v\n", event.Metadata())
				fmt.Println()

				// Process the order (simulate business logic)
				fmt.Printf("  → Processing order: %s\n", payload.OrderID)
				fmt.Printf("  → Order amount: $%.2f\n", payload.Amount)
				fmt.Printf("  → Order items: %d\n", len(payload.Items))
				fmt.Println("  ✓ Order processed successfully!")
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
			fmt.Printf("Total elapsed time: %v\n", time.Since(startTime).Round(time.Second))
			return
		}
	}
}
