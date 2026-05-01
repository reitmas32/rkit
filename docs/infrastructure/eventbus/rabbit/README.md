# Event Bus - RabbitMQ Implementation

Implementación del event bus usando RabbitMQ, ideal para producción, microservicios y aplicaciones distribuidas.

## Descripción

La implementación RabbitMQ proporciona un event bus robusto, escalable y con alta disponibilidad usando RabbitMQ como broker de mensajes.

**Características:**
- ✅ Persistencia de mensajes
- ✅ Funciona entre procesos y máquinas
- ✅ Alta disponibilidad y escalabilidad
- ✅ Delay usando DLX (Dead Letter Exchange) con TTL
- ✅ Broadcast mediante exchange + múltiples queues
- ✅ Soporte para múltiples colas por evento (ConsumeWithQueue)
- ✅ Message acknowledgment (Ack/Nack/Reject)
- ✅ QoS (Quality of Service) configurable

## Requisitos

- RabbitMQ server instalado y corriendo
- Go 1.19+
- Dependencia: `github.com/rabbitmq/amqp091-go`

## Instalación

```bash
go get github.com/foundathyon/base/infrastructure/eventbus/rabbit
```

```go
import "github.com/foundathyon/base/infrastructure/eventbus/rabbit"
```

## Configuración

### Config

```go
type Config struct {
    URL          string  // AMQP connection URL (ej: "amqp://guest:guest@localhost:5672/")
    ExchangeName string  // Nombre del exchange (ej: "events")
    ExchangeType string  // Tipo de exchange (direct, topic, fanout, headers)
    QueuePrefix  string  // Prefijo opcional para nombres de queues
    Durable      bool    // Si queues y exchanges sobreviven reinicios
    AutoDelete   bool    // Si queues/exchanges se eliminan cuando no se usan
    PrefetchCount int    // Número de mensajes no confirmados por consumer
    PrefetchSize  int    // Tamaño del prefetch window en bytes (0 = ilimitado)
}
```

### DefaultConfig

```go
func DefaultConfig(url string) Config {
    return Config{
        URL:          url,
        ExchangeName: "events",
        ExchangeType: "topic",
        QueuePrefix:  "",
        Durable:      true,
        AutoDelete:   false,
        PrefetchCount: 10,
        PrefetchSize:  0,
    }
}
```

### Configuración Personalizada

```go
config := rabbit.DefaultConfig("amqp://guest:guest@localhost:5672/")
config.ExchangeName = "my-events"
config.ExchangeType = "topic"
config.QueuePrefix = "my-service"
config.PrefetchCount = 20
```

## Uso Básico

### Crear Event Bus

```go
import (
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/core/eventbus"
    "github.com/foundathyon/base/infrastructure/eventbus/rabbit"
)

// Configuración
config := rabbit.DefaultConfig("amqp://guest:guest@localhost:5672/")
config.ExchangeName = "events"
config.ExchangeType = "topic"

// Event factory: crea instancias de eventos basado en el nombre
eventFactory := func(eventName string) eventbus.Event {
    switch eventName {
    case "user.created":
        return &UserCreatedEvent{}
    case "user.updated":
        return &UserUpdatedEvent{}
    default:
        return nil
    }
}

// Crear event bus
eb, err := rabbit.NewEventBus(config, eventFactory)
if err != nil {
    log.Fatal(err)
}
defer eb.Close()
```

**Event Factory:** La función `eventFactory` es requerida para deserializar eventos. Debe retornar una nueva instancia del tipo de evento correspondiente al `eventName`.

### Definir Evento

```go
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

// Implementar interfaz Event
func (e *UserCreatedEvent) Name() string                { return e.EventName }
func (e *UserCreatedEvent) Version() string             { return e.EventVersion }
func (e *UserCreatedEvent) OccurredAt() time.Time       { return e.EventOccurredAt }
func (e *UserCreatedEvent) Payload() any                { return e.EventPayload }
func (e *UserCreatedEvent) Metadata() eventbus.Metadata { return e.EventMetadata }
```

**Importante:** Los campos del evento deben ser públicos (capitalizados) y tener tags JSON para la serialización/deserialización.

### Publicar Evento

```go
ctx := customctx.New(context.Background())

event := &UserCreatedEvent{
    EventName:       "user.created",
    EventVersion:    "1.0",
    EventOccurredAt: time.Now(),
    EventPayload: UserPayload{
        UserID:   "123",
        Username: "john_doe",
        Email:    "john@example.com",
    },
    EventMetadata: eventbus.Metadata{
        "trace_id": "abc-123",
        "source":   "api",
    },
}

if err := eb.Publish(ctx, event); err != nil {
    log.Fatal(err)
}
```

### Consumir Eventos

#### Consumo Básico (Consume)

```go
ctx := customctx.New(context.Background())

// Crear template del evento
eventTemplate := &UserCreatedEvent{
    EventName:       "user.created",
    EventVersion:    "1.0",
    EventOccurredAt: time.Now(),
    EventPayload:    UserPayload{},
    EventMetadata:   eventbus.Metadata{},
}

// Suscribirse
consumeResult := eb.Consume(ctx, eventTemplate)
if !consumeResult.IsOk() {
    log.Fatal(consumeResult.Error())
}

deliveryChan := consumeResult.Value()

// Procesar mensajes
for msg := range deliveryChan {
    event := msg.Event()
    if userEvent, ok := event.(*UserCreatedEvent); ok {
        payload := userEvent.Payload().(UserPayload)
        fmt.Printf("Received: UserID=%s, Email=%s\n", payload.UserID, payload.Email)
    }
    
    // Confirmar procesamiento
    if err := msg.Ack(); err != nil {
        log.Printf("Error acknowledging: %v", err)
    }
}
```

**Nota:** `Consume` crea una queue automáticamente usando el nombre del evento (o con prefijo si está configurado).

#### Consumo con Queue Personalizada (ConsumeWithQueue)

Para múltiples consumidores del mismo evento (patrón broadcast), usa `ConsumeWithQueue`:

```go
// Consumer 1: Email Service
queueName1 := "email-service.user.created"
consumeResult1 := eb.ConsumeWithQueue(ctx, eventTemplate, queueName1)
deliveryChan1 := consumeResult1.Value()

// Consumer 2: Push Notification Service
queueName2 := "push-service.user.created"
consumeResult2 := eb.ConsumeWithQueue(ctx, eventTemplate, queueName2)
deliveryChan2 := consumeResult2.Value()

// Ambos consumidores recibirán el mismo evento
eb.Publish(ctx, event)
// → email-service.user.created recibe el evento
// → push-service.user.created recibe el evento
```

**Ventajas:**
- Cada consumidor tiene su propia queue
- Todos los consumidores reciben el mismo evento (broadcast)
- Permite procesar eventos a diferentes velocidades

### Publicar con Delay

```go
delay := 5 * time.Second
err := eb.PublishWithDelay(ctx, event, delay)
```

**Cómo funciona:**
1. Se crea una queue temporal con TTL = delay
2. Se configura DLX (Dead Letter Exchange) al exchange principal
3. El mensaje se publica a la queue temporal
4. Cuando expira el TTL, RabbitMQ reenvía el mensaje al DLX con la routing key original
5. Los consumidores normales reciben el mensaje

**Limitaciones:**
- No requiere plugins adicionales (usa funcionalidad nativa de RabbitMQ)
- El delay máximo depende de la configuración de RabbitMQ (por defecto ~2^31 ms ≈ 24 días)

## API

### NewEventBus

```go
func NewEventBus(config Config, eventFactory func(eventName string) eventbus.Event) (*EventBus, error)
```

Crea una nueva instancia del event bus RabbitMQ.

**Parámetros:**
- `config`: Configuración del event bus
- `eventFactory`: Función que crea instancias de eventos basado en el nombre

**Retorna:**
- `*EventBus`: Instancia del event bus
- `error`: Error si falla la conexión o configuración

### Publish

```go
func (eb *EventBus) Publish(ctx *customctx.CustomContext, event eventbus.Event) error
```

Publica un evento inmediatamente al exchange usando el nombre del evento como routing key.

### PublishWithDelay

```go
func (eb *EventBus) PublishWithDelay(ctx *customctx.CustomContext, event eventbus.Event, delay time.Duration) error
```

Publica un evento que se entregará después del delay especificado usando DLX con TTL.

**Parámetros:**
- `ctx`: Contexto personalizado
- `event`: Evento a publicar
- `delay`: Duración del delay (si <= 0, usa Publish normal)

### Consume

```go
func (eb *EventBus) Consume(ctx *customctx.CustomContext, event eventbus.Event) result.Result[eventbus.DeliveryChannel]
```

Crea un consumidor usando una queue automática (basada en el nombre del evento).

### ConsumeWithQueue

```go
func (eb *EventBus) ConsumeWithQueue(ctx *customctx.CustomContext, event eventbus.Event, queueName string) result.Result[eventbus.DeliveryChannel]
```

Crea un consumidor usando una queue personalizada. Ideal para múltiples consumidores del mismo evento.

**Parámetros:**
- `ctx`: Contexto para cancelación graceful
- `event`: Template del evento (usa `Name()` para la routing key)
- `queueName`: Nombre de la queue (debe ser único por consumidor)

**Retorna:**
- `result.Result[DeliveryChannel]`: Canal de mensajes o error

### Close

```go
func (eb *EventBus) Close() error
```

Cierra las conexiones de RabbitMQ y libera recursos.

## Ejemplos

### Ejemplos Básicos

Ver `examples/infrastructure/eventbus/rabbit/basic/`:

- **`publisher/main.go`**: Publicación básica de eventos
- **`consumer/main.go`**: Consumo básico de eventos

### Ejemplos con Delay

Ver `examples/infrastructure/eventbus/rabbit/basic/`:

- **`delayed_publisher/main.go`**: Publicación con diferentes delays
- **`delayed_consumer/main.go`**: Consumo de eventos con delay

### Ejemplo de Broadcast

Ver `examples/infrastructure/eventbus/rabbit/signup/`:

- **`publisher/main.go`**: Publica eventos `user.signup`
- **`email/main.go`**: Consumer de email service (queue: `email-service.user.signup`)
- **`push/main.go`**: Consumer de push service (queue: `push-service.user.signup`)

Este ejemplo demuestra cómo múltiples servicios pueden consumir el mismo evento.

## Casos de Uso

### 1. Microservicios

Ideal para comunicación asíncrona entre microservicios:

```go
// Order Service
eb.Publish(ctx, &OrderCreatedEvent{...})

// Email Service (recibe el evento)
eb.ConsumeWithQueue(ctx, eventTemplate, "email-service.order.created")
```

### 2. Procesamiento Asíncrono

Procesa tareas pesadas de manera asíncrona:

```go
// Publicar tarea
eb.Publish(ctx, &ProcessImageEvent{ImageID: "123"})

// Worker processa en background
eb.ConsumeWithQueue(ctx, eventTemplate, "worker.process-image")
```

### 3. Tareas Programadas

Usa delay para tareas programadas:

```go
// Enviar recordatorio en 1 hora
delay := 1 * time.Hour
eb.PublishWithDelay(ctx, &SendReminderEvent{...}, delay)
```

### 4. Event Sourcing

Almacena eventos para auditoría y replay:

```go
// Todos los eventos se persisten en RabbitMQ
eb.Publish(ctx, event)
```

## Mejores Prácticas

### 1. Event Factory

Implementa una factory robusta que maneje todos los tipos de eventos:

```go
eventFactory := func(eventName string) eventbus.Event {
    switch eventName {
    case "user.created":
        return &UserCreatedEvent{}
    case "user.updated":
        return &UserUpdatedEvent{}
    case "order.created":
        return &OrderCreatedEvent{}
    default:
        log.Printf("Unknown event type: %s", eventName)
        return nil
    }
}
```

### 2. Acknowledgment

Siempre haz `Ack()` después de procesar exitosamente:

```go
for msg := range deliveryChan {
    if err := processEvent(msg.Event()); err != nil {
        // Reencolar para reintento
        msg.Nack(true)
        continue
    }
    
    // Confirmar procesamiento
    msg.Ack()
}
```

### 3. Error Handling

Maneja errores de deserialización:

```go
for msg := range deliveryChan {
    event := msg.Event()
    if event == nil {
        // Evento desconocido o error de deserialización
        msg.Nack(false) // No reencolar
        continue
    }
    
    // Procesar...
    msg.Ack()
}
```

### 4. Context Cancellation

Usa context para cancelación graceful:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

customCtx := customctx.New(ctx)
consumeResult := eb.Consume(customCtx, eventTemplate)
```

### 5. Queue Naming

Usa nombres descriptivos para queues:

```go
// ✅ Bueno
"email-service.user.created"
"push-service.user.signup"
"worker.process-image"

// ❌ Malo
"q1"
"consumer"
"queue"
```

### 6. Prefetch Count

Ajusta `PrefetchCount` según tu caso de uso:

```go
// Para procesamiento rápido
config.PrefetchCount = 100

// Para procesamiento lento (una tarea a la vez)
config.PrefetchCount = 1
```

### 7. Durable vs AutoDelete

- **Durable = true**: Las queues/exchanges sobreviven reinicios (producción)
- **AutoDelete = true**: Las queues se eliminan cuando no hay consumidores (desarrollo/testing)

## Configuración de RabbitMQ

### Requisitos Mínimos

- RabbitMQ 3.8+
- Exchange tipo `topic` (recomendado) o `direct`
- No se requieren plugins adicionales (delay usa DLX nativo)

### Recomendaciones de Producción

1. **Clusters**: Usa un cluster de RabbitMQ para alta disponibilidad
2. **Persistencia**: Usa `Durable = true` para queues y exchanges
3. **Replicación**: Configura mirroring de queues en el cluster
4. **Monitoring**: Monitorea queues, mensajes y conexiones
5. **Connection Pooling**: Reutiliza conexiones (el EventBus maneja esto automáticamente)

### Configuración de RabbitMQ Management

```bash
# Habilitar management plugin
rabbitmq-plugins enable rabbitmq_management

# Acceder a la interfaz web
# http://localhost:15672
# Usuario: guest / Contraseña: guest
```

## Troubleshooting

### Error: "connection refused"

- Verifica que RabbitMQ esté corriendo: `rabbitmq-server status`
- Verifica la URL de conexión
- Verifica firewall/red

### Mensajes no se entregan

- Verifica que el exchange exista
- Verifica que las queues estén vinculadas al exchange
- Verifica la routing key (debe ser el nombre del evento)

### Delay no funciona

- Verifica que las queues temporales se creen correctamente
- Verifica logs de RabbitMQ para errores de DLX
- Verifica que el exchange principal exista

### Eventos no se deserializan

- Verifica que el event factory retorne el tipo correcto
- Verifica que los campos del evento sean públicos
- Verifica que los tags JSON sean correctos
- Verifica logs para errores de unmarshal

## Comparación con In-Memory

| Característica | RabbitMQ | In-Memory |
|----------------|----------|-----------|
| Persistencia | ✅ | ❌ |
| Entre procesos | ✅ | ❌ |
| Alta disponibilidad | ✅ | ❌ |
| Delay preciso | ✅ | ⚠️ |
| Producción | ✅ | ❌ |
| Testing | ⚠️ | ✅ |
| Complejidad | ⚠️ | ✅ |

## Referencias

- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
- [AMQP 0-9-1 Model](https://www.rabbitmq.com/tutorials/amqp-concepts.html)
- [Dead Letter Exchanges](https://www.rabbitmq.com/dlx.html)
- Ver documentación del [eventbus core](../../core/eventbus/README.md)
