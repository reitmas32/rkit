# rkit

Biblioteca general de Go con utilidades base para construir microservicios y aplicaciones backend. Proporciona abstracciones de dominio, implementaciones de infraestructura y patrones comunes listos para usar.

```
go get github.com/reitmas32/rkit
```

> Requiere Go 1.25+

---

## Contenido

- [Paquetes del núcleo](#paquetes-del-núcleo)
  - [kerrors](#kerrors--errores-estructurados)
  - [result](#result--manejo-funcional-de-errores)
  - [customctx](#customctx--contexto-personalizado)
  - [logger](#logger)
  - [eventbus](#eventbus--bus-de-eventos)
  - [http](#http--cliente-http)
  - [types](#types)
- [Infraestructura](#infraestructura)
  - [infrastructure/http](#infrastructurehttp)
  - [infrastructure/eventbus](#infrastructureeventbus)
- [Persistencia](#persistencia)
  - [Contratos](#contratos)
  - [Criteria y Filtros](#criteria-y-filtros)
  - [Paginación](#paginación)
  - [persistence/inmemory](#persistenceinmemory)
  - [persistence/postgres](#persistencepostgres)
- [Observabilidad](#observabilidad)

---

## Paquetes del núcleo

### `kerrors` — Errores estructurados

Extiende el manejo de errores de Go con código, metadata y encadenamiento.

```go
import "github.com/reitmas32/rkit/core/kerrors"

// Error simple
err := kerrors.NewKError("recurso no encontrado", 404, nil)

// Con metadata de contexto
err := kerrors.NewKError("validación fallida", 400, map[string]any{
    "field":      "email",
    "request_id": "req-abc-123",
})

// Con causa subyacente
dbErr := errors.New("connection timeout")
err := kerrors.NewKErrorWithCause("error al guardar usuario", 500, nil, dbErr)

// Encadenamiento compatible con errors.Is / errors.Unwrap
errors.Is(err, dbErr) // true
```

---

### `result` — Manejo funcional de errores

Tipo genérico inspirado en el `Result<T, E>` de Rust. Evita el patrón `value, err` en cascada.

```go
import (
    "github.com/reitmas32/rkit/core/kerrors"
    "github.com/reitmas32/rkit/core/result"
)

func divide(a, b int) result.Result[float64] {
    if b == 0 {
        return result.Err[float64](kerrors.NewKError("división por cero", 400, nil))
    }
    return result.Ok(float64(a) / float64(b))
}

r := divide(10, 2)
if r.IsOk() {
    fmt.Println(r.Value()) // 5.0
} else {
    fmt.Println(r.Error())
}
```

| Función | Descripción |
|---------|-------------|
| `result.Ok(value)` | Result exitoso |
| `result.Err[T](kerr)` | Result con error |
| `result.Empty[T]()` | Result sin valor ni error |
| `.IsOk()` | Verifica si es exitoso |
| `.IsEmpty()` | Verifica si está vacío |
| `.Value()` | Retorna el valor |
| `.Error()` | Retorna el error como `error` |
| `.ToKError()` | Retorna el error como `*KError` |

---

### `customctx` — Contexto personalizado

Implementación de `context.Context` que acumula errores estructurados a lo largo de un flujo, con tracking del call site y soporte de logger.

```go
import (
    "context"
    "github.com/reitmas32/rkit/core/customctx"
    "github.com/reitmas32/rkit/core/kerrors"
)

ctx := customctx.New(context.Background())

// Acumular errores sin detener la ejecución
ctx.AddError(kerrors.NewKError("validación fallida", 400, nil))

// Consultar errores acumulados
if ctx.HasErrors() {
    errs := ctx.GetErrors()
    // ...
}

// Almacenar valores con clave string
ctx.SetValue("request_id", "req-abc-123")
reqID := customctx.ExtractValue[string](ctx, "request_id")

// Compatible con context.WithTimeout, etc.
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
```

---

### `logger`

Contrato de logger y una implementación simple integrable en cualquier componente.

```go
import "github.com/reitmas32/rkit/core/logger"

log := logger.NewSimpleLogger()
log.Info("servidor iniciado", map[string]any{"port": 8080})
log.Error("fallo inesperado", map[string]any{"err": err})
```

Niveles disponibles: `Debug`, `Info`, `Warn`, `Error`.

---

### `eventbus` — Bus de eventos

Contratos para publicar y consumir eventos de forma desacoplada.

```go
import "github.com/reitmas32/rkit/core/eventbus"

// Implementar Publisher
type MyPublisher struct{}
func (p *MyPublisher) Publish(ctx context.Context, event eventbus.Event) error { ... }

// Implementar Consumer
type MyConsumer struct{}
func (c *MyConsumer) Consume(ctx context.Context, handler eventbus.HandlerFunc) error { ... }
```

Implementaciones incluidas en `infrastructure/eventbus`.

---

### `http` — Cliente HTTP

Contratos y helpers para hacer peticiones HTTP con respuestas tipadas.

```go
import (
    "github.com/reitmas32/rkit/core/customctx"
    corehttp "github.com/reitmas32/rkit/core/http"
    infrahttp "github.com/reitmas32/rkit/infrastructure/http"
)

ctx := customctx.New(context.Background())
client := infrahttp.NewClient(infrahttp.DefaultConfig())

type Character struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

resp, err := corehttp.GetTyped[Character](
    client,
    ctx,
    "https://rickandmortyapi.com/api/character/1",
    corehttp.WithHeader("Accept", "application/json"),
)

if err == nil && resp.IsSuccess() {
    fmt.Println(resp.Body.Name)
}
```

---

### `types`

```go
import "github.com/reitmas32/rkit/core/types"

id := types.NewUUIDx()         // genera nuevo UUID
id, err := types.ParseUUIDx(s) // parsea desde string
```

---

## Infraestructura

### `infrastructure/http`

Cliente HTTP construido sobre `net/http` con soporte para reintentos, logging y contexto personalizado.

```go
import infrahttp "github.com/reitmas32/rkit/infrastructure/http"

config := infrahttp.DefaultConfig()
config.Timeout = 10 * time.Second
config.MaxRetries = 3

client := infrahttp.NewClient(config)
```

### `infrastructure/eventbus`

**In-memory** (útil para tests y desarrollo local):

```go
import "github.com/reitmas32/rkit/infrastructure/eventbus/inmemory"

bus := inmemory.NewEventBus()
```

**RabbitMQ**:

```go
import "github.com/reitmas32/rkit/infrastructure/eventbus/rabbit"

bus, err := rabbit.NewEventBus(rabbit.Config{
    URL:      "amqp://guest:guest@localhost:5672/",
    Exchange: "my.exchange",
})
```

---

## Persistencia

### Contratos

```go
import "github.com/reitmas32/rkit/persistence/contracts"

// IEntity — entidad de dominio con ID
// IModel  — modelo de base de datos
```

### Criteria y Filtros

Sistema de filtros componibles para construir queries sin acoplar al ORM.

```go
import (
    "github.com/reitmas32/rkit/persistence/criteria"
)

c := criteria.New().
    AddFilter(criteria.NewFilter("status", criteria.EQ, "active")).
    AddFilter(criteria.NewFilter("age", criteria.GT, 18))
```

### Paginación

```go
import "github.com/reitmas32/rkit/persistence/pagination"

req := pagination.NewPageRequest(1, 20,
    pagination.NewSort("created_at", pagination.DESC),
)

// El repositorio retorna:
// PageResult[T] con Items, Total, Page, PageSize, TotalPages
```

### `persistence/inmemory`

Repositorio genérico en memoria, ideal para tests y prototipado.

```go
import "github.com/reitmas32/rkit/persistence/inmemory"

repo := inmemory.NewRepository[MyEntity, MyModel]()

repo.Save(ctx, entity)
r := repo.GetByID(ctx, id)
repo.DeleteByID(ctx, id)
```

### `persistence/postgres`

Repositorio genérico sobre GORM para PostgreSQL.

```go
import (
    "github.com/reitmas32/rkit/persistence/postgres"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

repo := &postgres.PostgresRepository[MyEntity, MyModel]{
    Connection: db,
}

repo.Save(ctx, entity)
r := repo.GetByID(ctx, id)
```

---

## Observabilidad

### Loguru (Logrus + Loki)

Logger con formato estructurado, hook para Grafana Loki y campos HTTP/WebSocket predefinidos.

```go
import (
    "github.com/reitmas32/rkit/observability/logger/loguru"
    "github.com/reitmas32/rkit/observability/logger/loguru/fields"
)

log := loguru.New(loguru.Config{
    Level:  "info",
    Format: "json",
})

log.WithFields(fields.HTTP(r)).Info("request received")
log.WithFields(fields.WS(conn)).Info("websocket connected")
```

**Hook para Loki** (Grafana):

```go
import "github.com/reitmas32/rkit/observability/logger/loguru/hooks"

log.AddHook(hooks.NewLoki(hooks.LokiConfig{
    URL:    "http://localhost:3100",
    Labels: map[string]string{"app": "my-service"},
}))
```

Para levantar un servidor Loki mock local, ver `mock/loki/`.

---

## Estructura del proyecto

```
rkit/
├── core/
│   ├── customctx/     # Contexto con acumulación de errores
│   ├── eventbus/      # Contratos de bus de eventos
│   ├── http/          # Contratos y tipos HTTP
│   ├── kerrors/       # Errores estructurados
│   ├── logger/        # Contrato de logger
│   ├── result/        # Tipo Result genérico
│   └── types/         # UUIDx y otros tipos base
├── infrastructure/
│   ├── dtos/          # DTOs compartidos
│   ├── eventbus/      # In-memory y RabbitMQ
│   └── http/          # Cliente HTTP con reintentos
├── observability/
│   └── logger/loguru/ # Logger estructurado + Loki
├── persistence/
│   ├── contracts/     # IEntity, IModel
│   ├── criteria/      # Filtros componibles
│   ├── inmemory/      # Repositorio en memoria
│   ├── models/        # Entidad base con mutations
│   ├── pagination/    # PageRequest / PageResult
│   └── postgres/      # Repositorio PostgreSQL (GORM)
├── examples/          # Ejemplos ejecutables por paquete
├── docs/              # Documentación detallada por paquete
└── mock/loki/         # Servidor Loki mock para desarrollo
```

---

## Licencia

MIT
