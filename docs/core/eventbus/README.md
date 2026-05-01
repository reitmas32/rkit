# Event Bus

El paquete `eventbus` proporciona una abstracción para un bus de eventos asíncrono que permite publicar y consumir eventos de manera desacoplada en la aplicación.

## Descripción General

El event bus sigue el patrón **Publisher-Subscriber (Pub/Sub)**, permitiendo:

- **Publicación de eventos**: Los productores publican eventos sin conocer quién los consume
- **Consumo de eventos**: Los consumidores se suscriben a tipos de eventos sin conocer quién los publica
- **Desacoplamiento**: Los componentes de la aplicación se comunican de manera asíncrona a través de eventos
- **Escalabilidad**: Múltiples consumidores pueden procesar el mismo evento (patrón broadcast)

## Arquitectura

El event bus está diseñado siguiendo los principios de **Clean Architecture**:

```
┌─────────────────────────────────────────┐
│         Core (Interfaces)               │
│  Event, Publisher, Consumer, Message    │
└─────────────────────────────────────────┘
                   ▲
                   │
        ┌──────────┴──────────┐
        │                     │
┌──────────────────┐  ┌──────────────────┐
│  Infrastructure  │  │  Infrastructure  │
│   In-Memory      │  │     RabbitMQ     │
└──────────────────┘  └──────────────────┘
```

### Capas

1. **Core**: Define las interfaces y tipos base (`Event`, `Publisher`, `Consumer`, `Message`)
2. **Infrastructure**: Implementaciones concretas (in-memory, RabbitMQ, etc.)

## Componentes Principales

### Event

La interfaz `Event` representa un evento que puede ser publicado y consumido:

```go
type Event interface {
    Name() string              // Nombre único del evento (ej: "user.created")
    Version() string           // Versión del evento (ej: "1.0")
    OccurredAt() time.Time     // Timestamp de cuándo ocurrió el evento
    Payload() any              // Datos del evento (estructura específica)
    Metadata() Metadata        // Metadatos adicionales (map[string]string)
}
```

**Características:**
- **Inmutabilidad**: Los eventos son inmutables una vez creados
- **Versionado**: Cada evento tiene una versión para manejar cambios de esquema
- **Metadatos**: Permite agregar información adicional (trace_id, source, etc.)

### Message

La interfaz `Message` encapsula un evento recibido con información de entrega:

```go
type Message interface {
    Event() Event              // El evento recibido
    Ack() error                // Confirmar procesamiento exitoso
    Nack(requeue bool) error   // Rechazar con opción de reencolar
    Reject(requeue bool) error // Rechazar (alias de Nack)
    DeliveryTag() uint64       // Identificador único de entrega
    Timestamp() time.Time      // Timestamp del mensaje
    Headers() map[string]any   // Headers adicionales
}
```

**Métodos de Acknowledgment:**
- **`Ack()`**: Confirma que el mensaje fue procesado exitosamente
- **`Nack(requeue bool)`**: Rechaza el mensaje; si `requeue=true`, se reencola
- **`Reject(requeue bool)`**: Alias de `Nack`

### Publisher

La interfaz `Publisher` define cómo publicar eventos:

```go
type Publisher interface {
    Publish(ctx *customctx.CustomContext, event Event) error
    PublishWithDelay(ctx *customctx.CustomContext, event Event, delay time.Duration) error
}
```

**Métodos:**
- **`Publish()`**: Publica un evento inmediatamente
- **`PublishWithDelay()`**: Publica un evento que se entregará después del delay especificado

### Consumer

La interfaz `Consumer` define cómo consumir eventos:

```go
type Consumer interface {
    Consume(ctx *customctx.CustomContext, event Event) result.Result[DeliveryChannel]
}
```

**Retorna:**
- **`DeliveryChannel`**: Un canal (`<-chan Message`) por el cual se reciben mensajes
- **`result.Result`**: Maneja errores de manera explícita usando el tipo `Result`

### EventBus

La interfaz `EventBus` combina `Publisher` y `Consumer`:

```go
type EventBus interface {
    Publisher
    Consumer
    Close() error  // Cerrar conexiones y liberar recursos
}
```

## Tipos y Utilidades

### Metadata

```go
type Metadata map[string]string
```

Mapa de metadatos adicionales asociados con un evento (ej: `trace_id`, `source`, `user_id`).

### DeliveryChannel

```go
type DeliveryChannel <-chan Message
```

Canal de mensajes recibidos por un consumidor.

## Flujo de Uso

### Publicar un Evento

```go
ctx := customctx.New(context.Background())
event := &UserCreatedEvent{
    EventName: "user.created",
    EventVersion: "1.0",
    EventOccurredAt: time.Now(),
    EventPayload: UserPayload{...},
    EventMetadata: eventbus.Metadata{
        "trace_id": "abc-123",
        "source": "api",
    },
}

err := eventBus.Publish(ctx, event)
```

### Publicar con Delay

```go
delay := 5 * time.Second
err := eventBus.PublishWithDelay(ctx, event, delay)
```

### Consumir Eventos

```go
ctx := customctx.New(context.Background())
eventTemplate := &UserCreatedEvent{...}

consumeResult := eventBus.Consume(ctx, eventTemplate)
if !consumeResult.IsOk() {
    return consumeResult.Error()
}

deliveryChan := consumeResult.Value()
for msg := range deliveryChan {
    event := msg.Event()
    // Procesar evento
    if err := msg.Ack(); err != nil {
        // Manejar error
    }
}
```

## Patrón de Broadcast

El event bus soporta el patrón **broadcast**, donde múltiples consumidores pueden recibir el mismo evento:

```
Publisher ──> Exchange ──> Queue 1 ──> Consumer 1
                          └─> Queue 2 ──> Consumer 2
                          └─> Queue 3 ──> Consumer 3
```

Cada consumidor tiene su propia cola, y todas las colas están vinculadas al mismo exchange con la misma routing key (event name).

## Implementaciones Disponibles

### In-Memory

Ideal para:
- Testing y desarrollo
- Aplicaciones monolíticas simples
- Prototipado rápido

**Características:**
- Implementación en memoria (no persiste eventos)
- Thread-safe
- Soporte para delay mediante goroutines
- Broadcast a múltiples consumidores

**Ubicación:** `infrastructure/eventbus/inmemory`

### RabbitMQ

Ideal para:
- Producción
- Microservicios distribuidos
- Alta disponibilidad y escalabilidad
- Integración con sistemas externos

**Características:**
- Persistencia de mensajes
- Alta disponibilidad
- Delay usando DLX (Dead Letter Exchange) con TTL
- Broadcast mediante exchange + múltiples queues
- Soporte para múltiples colas por evento

**Ubicación:** `infrastructure/eventbus/rabbit`

## Mejores Prácticas

### 1. Inmutabilidad de Eventos

Los eventos deben ser inmutables una vez creados. Evita modificar un evento después de publicarlo.

### 2. Versionado de Eventos

Usa versiones semánticas para tus eventos (`"1.0"`, `"2.0"`, etc.) para manejar cambios de esquema de manera compatible.

### 3. Acknowledgment

Siempre haz `Ack()` después de procesar exitosamente un mensaje. Si hay error, usa `Nack(requeue: true)` para reintentar.

```go
for msg := range deliveryChan {
    event := msg.Event()
    
    if err := processEvent(event); err != nil {
        // Reencolar para reintento
        msg.Nack(true)
        continue
    }
    
    // Confirmar procesamiento exitoso
    msg.Ack()
}
```

### 4. Manejo de Errores

Usa el tipo `Result` para manejar errores de manera explícita:

```go
consumeResult := eventBus.Consume(ctx, eventTemplate)
if !consumeResult.IsOk() {
    return consumeResult.Error()
}
deliveryChan := consumeResult.Value()
```

### 5. Context Cancellation

Usa `customctx.CustomContext` para cancelar consumidores de manera graceful:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
customCtx := customctx.New(ctx)

consumeResult := eventBus.Consume(customCtx, eventTemplate)
```

### 6. Metadatos

Usa metadatos para información técnica (trace_id, source) pero evita poner datos de negocio. Los datos de negocio deben ir en el `Payload`.

```go
// ✅ Correcto
eventbus.Metadata{
    "trace_id": "abc-123",
    "source": "api",
}

// ❌ Incorrecto
eventbus.Metadata{
    "user_email": "user@example.com",  // Debe ir en Payload
}
```

### 7. Delay para Tareas Programadas

Usa `PublishWithDelay` para tareas programadas (ej: envío de emails después de 1 hora, recordatorios).

### 8. Múltiples Consumidores

Aprovecha el patrón broadcast para tener múltiples servicios procesando el mismo evento (ej: email service + push notification service).

## Ejemplos

Ver la carpeta `examples/infrastructure/eventbus/` para ejemplos completos de uso:

- **In-Memory**: Ejemplo básico de publish/consume
- **RabbitMQ**:
  - `basic/publisher` y `basic/consumer`: Ejemplos básicos
  - `basic/delayed_publisher` y `basic/delayed_consumer`: Ejemplos con delay
  - `signup/publisher`, `signup/email`, `signup/push`: Ejemplo de broadcast

## Referencias

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Event-Driven Architecture](https://martinfowler.com/articles/201701-event-driven.html)
- [Publisher-Subscriber Pattern](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern)
