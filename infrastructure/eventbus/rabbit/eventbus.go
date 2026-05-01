package rabbit

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/eventbus"
	"github.com/reitmas32/rkit/core/result"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

// Config holds the configuration for the RabbitMQ event bus.
type Config struct {
	// URL is the AMQP connection URL (e.g., "amqp://guest:guest@localhost:5672/")
	URL string

	// ExchangeName is the name of the exchange to use for publishing events.
	ExchangeName string

	// ExchangeType is the type of exchange (direct, topic, fanout, headers).
	ExchangeType string

	// QueuePrefix is an optional prefix for queue names.
	// If empty, queues will be named based on event names.
	QueuePrefix string

	// Durable indicates if queues and exchanges should survive broker restarts.
	Durable bool

	// AutoDelete indicates if queues/exchanges should be deleted when unused.
	AutoDelete bool

	// PrefetchCount sets the number of unacknowledged messages per consumer.
	PrefetchCount int

	// PrefetchSize sets the prefetch window in bytes (0 means unlimited).
	PrefetchSize int
}

// DefaultConfig returns a default configuration.
func DefaultConfig(url string) Config {
	return Config{
		URL:           url,
		ExchangeName:  "events",
		ExchangeType:  "topic",
		QueuePrefix:   "",
		Durable:       true,
		AutoDelete:    false,
		PrefetchCount: 10,
		PrefetchSize:  0,
	}
}

// EventBus is a RabbitMQ implementation of eventbus.EventBus.
type EventBus struct {
	config       Config
	conn         *amqp091.Connection
	channel      *amqp091.Channel
	eventFactory func(eventName string) eventbus.Event
	mu           sync.RWMutex
	closed       bool
}

// NewEventBus creates a new RabbitMQ event bus with the given configuration.
// eventFactory is a function that creates an event instance given an event name.
// It should return a new instance of the event type that matches the event name.
func NewEventBus(config Config, eventFactory func(eventName string) eventbus.Event) (*EventBus, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("AMQP URL is required")
	}
	if config.ExchangeName == "" {
		return nil, fmt.Errorf("exchange name is required")
	}
	if eventFactory == nil {
		return nil, fmt.Errorf("event factory is required")
	}

	conn, err := amqp091.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS
	if err := channel.Qos(config.PrefetchCount, config.PrefetchSize, false); err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	eb := &EventBus{
		config:       config,
		conn:         conn,
		channel:      channel,
		eventFactory: eventFactory,
	}

	// Declare exchange
	if err := eb.declareExchange(); err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return eb, nil
}

// declareExchange declares the exchange.
func (eb *EventBus) declareExchange() error {
	return eb.channel.ExchangeDeclare(
		eb.config.ExchangeName, // name
		eb.config.ExchangeType, // type
		eb.config.Durable,      // durable
		eb.config.AutoDelete,   // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
}

// Publish publishes an event to RabbitMQ.
func (eb *EventBus) Publish(ctx *customctx.CustomContext, event eventbus.Event) error {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.closed {
		return fmt.Errorf("event bus is closed")
	}

	// Serialize event
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Build routing key from event name
	routingKey := event.Name()

	// Build headers from metadata
	headers := amqp091.Table{}
	for k, v := range event.Metadata() {
		headers[k] = v
	}

	// Add event name and version to headers
	headers["x-event-name"] = event.Name()
	headers["x-event-version"] = event.Version()

	// Publish message
	err = eb.channel.Publish(
		eb.config.ExchangeName, // exchange
		routingKey,             // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			Headers:      headers,
			Timestamp:    event.OccurredAt(),
			DeliveryMode: amqp091.Persistent, // Make message persistent
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// PublishWithDelay publishes an event after the specified delay.
// PublishWithDelay publishes an event that will be delivered after the specified delay.
// It uses RabbitMQ's Dead Letter Exchange (DLX) pattern with TTL:
// 1. Creates a temporary queue with TTL = delay
// 2. Configures DLX to the main exchange with original routing key
// 3. Publishes to the temporary queue
// 4. When the message expires, it's re-routed to the main exchange with the original routing key
func (eb *EventBus) PublishWithDelay(ctx *customctx.CustomContext, event eventbus.Event, delay time.Duration) error {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.closed {
		return fmt.Errorf("event bus is closed")
	}

	if delay <= 0 {
		// If no delay, use regular publish
		return eb.Publish(ctx, event)
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	delayMs := delay.Milliseconds()
	routingKey := event.Name()

	// Create temporary queue name with delay suffix
	tempQueueName := fmt.Sprintf("%s.delay.%d", routingKey, delayMs)

	// Declare temporary queue with TTL and DLX
	_, err = eb.channel.QueueDeclare(
		tempQueueName,
		false, // durable: false (temporary queue)
		true,  // autoDelete: true (delete when unused)
		false, // exclusive: false
		false, // noWait: false
		amqp091.Table{
			"x-message-ttl":             delayMs,                // TTL in milliseconds
			"x-dead-letter-exchange":    eb.config.ExchangeName, // DLX to main exchange
			"x-dead-letter-routing-key": routingKey,             // Original routing key
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare temporary delay queue: %w", err)
	}

	// Build headers
	headers := amqp091.Table{}
	for k, v := range event.Metadata() {
		headers[k] = v
	}
	headers["x-event-name"] = event.Name()
	headers["x-event-version"] = event.Version()

	// Publish to temporary queue (message will expire and go to DLX)
	err = eb.channel.Publish(
		"",            // exchange: empty string = default exchange (direct routing by queue name)
		tempQueueName, // routing key: queue name (direct routing)
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			Headers:      headers,
			Timestamp:    event.OccurredAt(),
			DeliveryMode: amqp091.Persistent,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish delayed event: %w", err)
	}

	return nil
}

// Consume creates a channel that receives messages for the specified event type.
// It creates a queue using the event name. For multiple consumers of the same event,
// use ConsumeWithQueue to specify a unique queue name per consumer.
func (eb *EventBus) Consume(ctx *customctx.CustomContext, event eventbus.Event) result.Result[eventbus.DeliveryChannel] {
	eventName := event.Name()
	queueName := eb.queueName(eventName)
	return eb.consumeWithQueueName(ctx, event, queueName)
}

// ConsumeWithQueue creates a channel that receives messages for the specified event type
// using a custom queue name. This allows multiple consumers to subscribe to the same event
// type with different queue names (e.g., one for email service, one for push notifications).
// The queue will be bound to the exchange using the event name as the routing key.
func (eb *EventBus) ConsumeWithQueue(ctx *customctx.CustomContext, event eventbus.Event, queueName string) result.Result[eventbus.DeliveryChannel] {
	if queueName == "" {
		return result.NewErrResult[eventbus.DeliveryChannel](
			eventbus.ErrEventFailedToSubscribe,
		)
	}
	return eb.consumeWithQueueName(ctx, event, queueName)
}

// consumeWithQueueName is the internal implementation that creates a consumer for a specific queue.
func (eb *EventBus) consumeWithQueueName(ctx *customctx.CustomContext, event eventbus.Event, queueName string) result.Result[eventbus.DeliveryChannel] {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.closed {
		return result.NewErrResult[eventbus.DeliveryChannel](
			eventbus.ErrEventFailedToSubscribe,
		)
	}

	eventName := event.Name()

	// Declare queue
	queue, err := eb.channel.QueueDeclare(
		queueName,            // name
		eb.config.Durable,    // durable
		eb.config.AutoDelete, // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return result.NewErrResult[eventbus.DeliveryChannel](
			eventbus.ErrEventFailedToSubscribe,
		)
	}

	// Bind queue to exchange using event name as routing key (binding key)
	// This allows the same event to be routed to multiple queues
	bindingKey := eventName
	err = eb.channel.QueueBind(
		queue.Name,             // queue name
		bindingKey,             // routing key (binding key = event name)
		eb.config.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return result.NewErrResult[eventbus.DeliveryChannel](
			eventbus.ErrEventFailedToSubscribe,
		)
	}

	// Start consuming
	deliveries, err := eb.channel.Consume(
		queue.Name, // queue
		"",         // consumer (empty = auto-generated)
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return result.NewErrResult[eventbus.DeliveryChannel](
			eventbus.ErrEventFailedToSubscribe,
		)
	}

	// Convert AMQP deliveries to eventbus.Message
	msgChan := make(chan eventbus.Message, 100)

	go func() {
		defer close(msgChan)

		for delivery := range deliveries {
			// Extract event name from headers
			eventNameFromHeader, ok := delivery.Headers["x-event-name"].(string)
			if !ok {
				eventNameFromHeader = eventName // Fallback
			}

			// Create event instance using factory
			eventInstance := eb.eventFactory(eventNameFromHeader)
			if eventInstance == nil {
				// Skip if factory returns nil
				delivery.Nack(false, false)
				continue
			}

			// Deserialize event
			if err := json.Unmarshal(delivery.Body, eventInstance); err != nil {
				// Nack on deserialization error
				delivery.Nack(false, false)
				continue
			}

			// Create message wrapper
			msg := newMessage(&delivery, eventInstance)

			// Send to channel or handle context cancellation
			select {
			case msgChan <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()

	return result.Ok(eventbus.DeliveryChannel(msgChan))
}

// Close closes the RabbitMQ connection and channel.
func (eb *EventBus) Close() error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.closed {
		return nil
	}

	eb.closed = true

	var errs []error

	if eb.channel != nil {
		if err := eb.channel.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if eb.conn != nil {
		if err := eb.conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing event bus: %v", errs)
	}

	return nil
}

// queueName generates a queue name from an event name.
func (eb *EventBus) queueName(eventName string) string {
	if eb.config.QueuePrefix != "" {
		return fmt.Sprintf("%s.%s", eb.config.QueuePrefix, eventName)
	}
	return eventName
}
