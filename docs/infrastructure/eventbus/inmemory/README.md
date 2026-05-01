# Event Bus - In-Memory Implementation

Implementación en memoria del event bus, ideal para testing, desarrollo y aplicaciones simples.

## Descripción

La implementación in-memory almacena eventos en memoria y los distribuye a consumidores suscritos. Es thread-safe y soporta múltiples consumidores (patrón broadcast).

**Características:**
- ✅ Thread-safe usando `sync.RWMutex`
- ✅ Soporte para múltiples consumidores del mismo evento
- ✅ Delay mediante goroutines
- ✅ Auto-delete de queues cuando no hay consumidores
- ✅ Sin persistencia (eventos se pierden al reiniciar)

**Limitaciones:**
- ❌ No persiste eventos (se pierden al reiniciar)
- ❌ No funciona entre procesos/máquinas
- ❌ Delay no es preciso para tiempos largos

## Instalación

```go
import "github.com/foundathyon/base/infrastructure/eventbus/inmemory"
```

## Uso Básico

### Crear Event Bus

```go
import (
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/core/eventbus"
    "github.com/foundathyon/base/infrastructure/eventbus/inmemory"
)

// Crear event bus
eb := inmemory.NewEventBus()

// Defer cleanup
defer eb.Close()
```

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

## Publicar con Delay

```go
delay := 5 * time.Second
err := eb.PublishWithDelay(ctx, event, delay)
```

**Nota:** El delay se implementa mediante goroutines, por lo que no es preciso para tiempos muy largos.

## Múltiples Consumidores (Broadcast)

La implementación in-memory soporta múltiples consumidores para el mismo evento:

```go
// Consumer 1
consumeResult1 := eb.Consume(ctx, eventTemplate)
deliveryChan1 := consumeResult1.Value()

// Consumer 2
consumeResult2 := eb.Consume(ctx, eventTemplate)
deliveryChan2 := consumeResult2.Value()

// Cuando se publica un evento, ambos consumidores lo reciben
eb.Publish(ctx, event)
// → deliveryChan1 recibe el mensaje
// → deliveryChan2 recibe el mensaje
```

## API

### NewEventBus

```go
func NewEventBus() *EventBus
```

Crea una nueva instancia del event bus in-memory.

### Publish

```go
func (eb *EventBus) Publish(ctx *customctx.CustomContext, event eventbus.Event) error
```

Publica un evento inmediatamente a todos los consumidores suscritos.

### PublishWithDelay

```go
func (eb *EventBus) PublishWithDelay(ctx *customctx.CustomContext, event eventbus.Event, delay time.Duration) error
```

Publica un evento que se entregará después del delay especificado.

### Consume

```go
func (eb *EventBus) Consume(ctx *customctx.CustomContext, event eventbus.Event) result.Result[eventbus.DeliveryChannel]
```

Crea un consumidor para el tipo de evento especificado. Retorna un canal de mensajes.

**Parámetros:**
- `ctx`: Contexto para cancelación graceful
- `event`: Template del evento (usa `Name()` para identificar el tipo)

**Retorna:**
- `result.Result[DeliveryChannel]`: Canal de mensajes o error

### Close

```go
func (eb *EventBus) Close() error
```

Cierra el event bus y libera recursos. Cierra todos los canales de consumidores.

## Ejemplos

### Ejemplo Completo

Ver `examples/infrastructure/eventbus/inmemory/main.go` para un ejemplo completo que demuestra:

- Publicación básica de eventos
- Consumo de eventos
- Múltiples consumidores
- Delay en publicación
- Manejo de acknowledgment

## Casos de Uso

### 1. Testing

Ideal para tests unitarios e integración:

```go
func TestUserService(t *testing.T) {
    eb := inmemory.NewEventBus()
    defer eb.Close()
    
    service := NewUserService(eb)
    // ... tests
}
```

### 2. Desarrollo Local

Perfecto para desarrollo sin necesidad de infraestructura externa:

```go
func main() {
    eb := inmemory.NewEventBus()
    defer eb.Close()
    
    // Desarrollar y probar sin RabbitMQ
}
```

### 3. Aplicaciones Monolíticas Simples

Para aplicaciones simples que no requieren persistencia:

```go
func main() {
    eb := inmemory.NewEventBus()
    defer eb.Close()
    
    // Lógica de la aplicación
}
```

## Limitaciones y Consideraciones

### No Persiste Eventos

Los eventos se pierden si la aplicación se reinicia. Para persistencia, usa la implementación RabbitMQ.

### No Funciona Entre Procesos

Solo funciona dentro del mismo proceso. Para comunicación entre procesos/máquinas, usa RabbitMQ.

### Delay No Preciso

El delay se implementa con goroutines y no es preciso para tiempos muy largos (horas/días). Para delays largos, usa RabbitMQ con DLX.

### Memoria

Todos los eventos se mantienen en memoria hasta que son consumidos. Para grandes volúmenes, considera RabbitMQ.

## Thread Safety

La implementación es thread-safe usando `sync.RWMutex`. Puedes usar el mismo `EventBus` desde múltiples goroutines de manera segura.

## Comparación con RabbitMQ

| Característica | In-Memory | RabbitMQ |
|----------------|-----------|----------|
| Persistencia | ❌ | ✅ |
| Entre procesos | ❌ | ✅ |
| Alta disponibilidad | ❌ | ✅ |
| Delay preciso | ⚠️ | ✅ |
| Facilidad de uso | ✅ | ⚠️ |
| Testing | ✅ | ⚠️ |
| Producción | ❌ | ✅ |

## Referencias

- Ver documentación del [eventbus core](../core/eventbus/README.md)
- Ver [ejemplo completo](../../../examples/infrastructure/eventbus/inmemory/main.go)
