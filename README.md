# rkit

General-purpose Go library with base utilities for building microservices and
backend applications. It provides domain abstractions, infrastructure
implementations, and ready-to-use common patterns.

```bash
go get github.com/reitmas32/rkit
```

> Requires Go 1.25+

## Documentation

- **API reference (godoc):** https://pkg.go.dev/github.com/reitmas32/rkit — every
  package has a doc comment and runnable `Example` functions.
- **For AI coding agents:** [`llms.txt`](./llms.txt) (curated index) and
  [`llms-full.txt`](./llms-full.txt) (expanded content). See [`AGENTS.md`](./AGENTS.md)
  to work *on* this repo.
- **Per-package guides:** [`docs/`](./docs).

## Package map

| Import path | What it gives you |
|-------------|-------------------|
| `core/kerrors` | Structured errors: code, metadata, cause chaining |
| `core/result` | Generic `Result[T]` for functional error handling |
| `core/customctx` | `context.Context` that accumulates structured errors + logger |
| `core/logger` | `ILogger` contract + `SimpleLogger` + log levels |
| `core/eventbus` | Transport-agnostic event bus contracts |
| `core/http` | HTTP request/response contracts + generic typed helpers |
| `core/types` | `DomainID` (module/entity-encoded identifiers) |
| `infrastructure/http` | `net/http` client with retries, logging, customctx |
| `infrastructure/eventbus/inmemory` | In-process event bus (tests/dev) |
| `infrastructure/eventbus/rabbit` | RabbitMQ event bus (incl. delayed delivery) |
| `infrastructure/dtos` | Bind + validate gin request DTOs |
| `observability/logger/loguru` | Structured logrus logger + fields + Loki hook |
| `persistence/contracts` | `IEntity` / `IModel` + conversion helpers |
| `persistence/criteria` | Composable, ORM-agnostic query filters |
| `persistence/pagination` | `PageRequest` / `PageResult[T]` / `Pageable` |
| `persistence/models` | Base `Entity` + mutation notifications |
| `persistence/inmemory` | Generic in-memory repository |
| `persistence/postgres` | Generic GORM/PostgreSQL repository |
| `persistence/mongodb` | Generic MongoDB repository |

## Core packages

### `kerrors` — structured errors

```go
import "github.com/reitmas32/rkit/core/kerrors"

// Simple error with a code
err := kerrors.NewKError("resource not found", 404, nil)

// With contextual metadata
err = kerrors.NewKError("validation failed", 400, map[string]any{
    "field":      "email",
    "request_id": "req-abc-123",
})

// Wrapping an underlying cause (errors.Is / errors.Unwrap aware)
dbErr := errors.New("connection timeout")
err = kerrors.NewKErrorWithCause("failed to save user", 500, nil, dbErr)
errors.Is(err, dbErr) // true
```

### `result` — functional error handling

Generic type inspired by Rust's `Result<T, E>`; avoids the cascading
`value, err` pattern.

```go
import (
    "github.com/reitmas32/rkit/core/kerrors"
    "github.com/reitmas32/rkit/core/result"
)

func divide(a, b int) result.Result[float64] {
    if b == 0 {
        return result.Err[float64](kerrors.NewKError("division by zero", 400, nil))
    }
    return result.Ok(float64(a) / float64(b))
}

r := divide(10, 2)
if r.IsOk() {
    fmt.Println(r.Value()) // 5
}
```

| Function / method | Description |
|-------------------|-------------|
| `result.Ok(value)` | Successful result |
| `result.Err[T](kerr)` | Result carrying an error |
| `result.Empty[T]()` | Result with neither value nor error |
| `.IsOk()` / `.IsEmpty()` | State checks |
| `.Value()` | The value (zero value if not Ok) |
| `.Error()` | The error as a standard `error` |
| `.ToKError()` | The error as `*kerrors.KError` |

### `customctx` — error-accumulating context

Implements `context.Context` while collecting structured errors across a flow,
with call-site tracking and logger support.

```go
import (
    "context"
    "github.com/reitmas32/rkit/core/customctx"
    "github.com/reitmas32/rkit/core/kerrors"
)

cc := customctx.New(context.Background())

// Accumulate errors without stopping execution
cc.AddError(kerrors.NewKError("validation failed", 400, nil))

if cc.HasErrors() {
    for _, e := range cc.Errors() {
        _ = e
    }
}
```

### `logger`

```go
import "github.com/reitmas32/rkit/core/logger"

log := logger.NewSimpleLogger("info") // debug | info | warn | error
log.Info("server started on port %d", 8080)
```

### `eventbus`

Contracts to publish/consume events in a decoupled way; implementations live in
`infrastructure/eventbus`.

```go
import "github.com/reitmas32/rkit/core/eventbus"

type Publisher interface {
    Publish(ctx context.Context, event eventbus.Event) error
}
```

### `http` — typed client helpers

```go
import (
    "github.com/reitmas32/rkit/core/customctx"
    corehttp "github.com/reitmas32/rkit/core/http"
    infrahttp "github.com/reitmas32/rkit/infrastructure/http"
)

cc := customctx.New(context.Background())
client := infrahttp.NewClient(infrahttp.DefaultConfig())

resp, err := client.Get(cc, "https://example.com/health")
// Or use the generic typed helpers: corehttp.PostTyped[Req, Res](client, cc, url, body)
```

### `types`

```go
import "github.com/reitmas32/rkit/core/types"

id, err := types.NewDomainID("02", "01") // module 02, entity 01
parsed, err := types.ParseDomainID(id.String())
```

## Infrastructure

### `infrastructure/http`

```go
import infrahttp "github.com/reitmas32/rkit/infrastructure/http"

cfg := infrahttp.DefaultConfig()
cfg.Timeout = 10 * time.Second
client := infrahttp.NewClient(cfg)
```

### `infrastructure/eventbus`

```go
import "github.com/reitmas32/rkit/infrastructure/eventbus/inmemory"

bus := inmemory.NewEventBus() // in-process; great for tests
```

RabbitMQ: `rabbit.NewEventBus(cfg, eventFactory)` — see
`examples/infrastructure/eventbus/rabbit` for publisher/consumer/delayed flows.

## Persistence

### Criteria & filters

```go
import "github.com/reitmas32/rkit/persistence/criteria"

filters := criteria.NewFilters([]criteria.Filter{
    {Field: "status", Operator: criteria.OperatorEqual, Value: "active"},
    {Field: "age", Operator: criteria.OperatorGreaterThan, Value: 18},
})
```

### Pagination

```go
import "github.com/reitmas32/rkit/persistence/pagination"

req := pagination.NewPageRequest(0, 20) // page (0-indexed), size
// Repositories return PageResult[T]{Content, TotalElements, TotalPages, ...}
```

### Repositories

```go
import "github.com/reitmas32/rkit/persistence/postgres"

repo := &postgres.PostgresRepository[UserEntity, UserModel]{Connection: db}
repo.Save(cc, entity)
r := repo.GetByID(cc, id)
```

In-memory equivalent for tests: `inmemory.NewInMemoryMapRepository[E, M](onMutation)`.

## Observability — Loguru (Logrus + Loki)

Structured logger with caller info, predefined HTTP/WebSocket fields and a
batched Grafana Loki hook.

```go
import (
    "github.com/reitmas32/rkit/observability/logger/loguru"
    "github.com/reitmas32/rkit/observability/logger/loguru/fields"
    "github.com/reitmas32/rkit/observability/logger/loguru/hooks"
)

f := &fields.HTTPFileds{}
f.UpdateOne("method", "GET")
f.UpdateOne("path", "/api/v1/users")

log := loguru.NewLogger(f)

// Ship logs to Loki in batches of 50; flush the remainder at request end.
hook := hooks.NewLokiBufferedHook(
    "http://localhost:3100/loki/api/v1/push", 50,
    map[string]string{"app": "my-service", "environment": "production"},
)
log.AddHook(hook)
defer hook.Flush()

log.Info("request received")
```

Console format: `DATE_TIME | LEVEL | FILE.FUNCTION:LINE | [fields] | MESSAGE`.
A mock Loki server for local development lives in `mock/loki/`.

## Project structure

```
rkit/
├── core/            # customctx, eventbus, http, kerrors, logger, result, types
├── infrastructure/  # dtos, eventbus (inmemory + rabbit), http client
├── observability/   # logger/loguru (+ fields, hooks, Loki)
├── persistence/     # contracts, criteria, pagination, models, inmemory, postgres, mongodb
├── examples/        # runnable examples per package
├── docs/            # detailed per-package guides
└── mock/loki/       # mock Loki server for development
```

## License

MIT
