# HTTP Client - Infrastructure

Implementación concreta del cliente HTTP usando la biblioteca estándar `net/http` de Go. Proporciona una implementación completa y lista para producción del `core/http.Client` interface.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Configuración](#configuración)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Logging](#logging)
- [Retries](#retries)
- [Ejemplos](#ejemplos)
- [Casos de Uso](#casos-de-uso)

## ✨ Características

- **Implementación completa**: Implementa todas las interfaces de `core/http`
- **Configuración flexible**: Config struct con múltiples opciones
- **Base URL**: Soporte para URLs base con resolución automática
- **Headers por defecto**: Headers que se añaden a todas las peticiones
- **Timeouts configurables**: Timeout global y por request
- **Logging integrado**: Soporte para logging opcional de requests/responses
- **Retries automáticos**: Mecanismo de retry configurable con backoff exponencial
- **Thread-safe**: Seguro para uso concurrente

## 📦 Instalación

```bash
go get github.com/foundathyon/base/infrastructure/http
```

## ⚙️ Configuración

### Config Struct

```go
type Config struct {
    // Timeout es el timeout por defecto para requests (en segundos)
    Timeout int

    // BaseURL es una URL base opcional que se antepondrá a todas las requests
    BaseURL string

    // DefaultHeaders son headers que se añadirán a todas las requests
    DefaultHeaders map[string]string

    // Transport es el HTTP transport a usar (opcional, usa http.DefaultTransport si nil)
    Transport http.RoundTripper

    // CheckRedirect especifica la política para manejar redirects
    CheckRedirect func(req *http.Request, via []*http.Request) error

    // Logger es un logger opcional para logging de requests/responses
    Logger logger.ILogger

    // DisableLogging desactiva explícitamente el logging incluso si Logger está configurado
    DisableLogging bool

    // MaxRetries es el número máximo de intentos de retry para requests fallidas
    // Si es 0, no se realizarán retries (por defecto)
    MaxRetries int

    // RetryDelay es el delay entre intentos de retry en milisegundos
    // Si es 0, se usará un delay por defecto de 100ms
    RetryDelay int

    // RetryableStatusCodes es una lista de códigos de estado HTTP que deberían activar un retry
    // Por defecto: 429, 500, 502, 503, 504
    RetryableStatusCodes []int

    // RetryableMethods especifica qué métodos HTTP deberían ser reintentados
    // Por defecto, solo métodos idempotentes (GET, HEAD, OPTIONS, DELETE)
    RetryableMethods []string
}
```

### DefaultConfig

```go
config := http.DefaultConfig()
// Timeout: 30 segundos
// BaseURL: ""
// DefaultHeaders: map vacío
// Logger: nil
// MaxRetries: 0 (sin retries)
// RetryDelay: 100ms
// RetryableStatusCodes: [429, 500, 502, 503, 504]
// RetryableMethods: ["GET", "HEAD", "OPTIONS", "DELETE"]
```

### Configuración Personalizada

```go
config := http.DefaultConfig()
config.Timeout = 60                    // 60 segundos
config.BaseURL = "https://api.example.com"
config.DefaultHeaders = map[string]string{
    "Authorization": "Bearer token",
    "X-API-Key":     "key",
}
config.MaxRetries = 3
config.RetryDelay = 200                // 200ms
config.RetryableStatusCodes = []int{500, 502, 503}
```

## 🚀 Uso Básico

### Cliente Simple

```go
import (
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/infrastructure/http"
)

ctx := customctx.New(context.Background())
client := http.NewClient(http.DefaultConfig())

resp, err := client.Get(ctx, "https://api.example.com/data")
if err != nil {
    log.Fatal(err)
}
defer resp.Close()
```

### Cliente con Base URL

```go
config := http.DefaultConfig()
config.BaseURL = "https://api.example.com"
client := http.NewClient(config)

// Esta URL se resolverá como: https://api.example.com/v1/users
resp, err := client.Get(ctx, "/v1/users")
```

### Cliente con Headers por Defecto

```go
config := http.DefaultConfig()
config.DefaultHeaders = map[string]string{
    "Authorization": "Bearer my-token",
    "X-Client-Version": "1.0.0",
}
client := http.NewClient(config)

// Todos los requests incluirán estos headers automáticamente
resp, err := client.Get(ctx, "https://api.example.com/data")
```

## 📚 API

### NewClient

Crea una nueva instancia del cliente HTTP:

```go
func NewClient(config Config) *Client
```

### Métodos HTTP

El cliente implementa todos los métodos de la interfaz `core/http.Client`:

- `Do(ctx, req) (*Response, error)`
- `Get(ctx, url, opts...) (*Response, error)`
- `Post(ctx, url, body, opts...) (*Response, error)`
- `Put(ctx, url, body, opts...) (*Response, error)`
- `Patch(ctx, url, body, opts...) (*Response, error)`
- `Delete(ctx, url, opts...) (*Response, error)`
- `Head(ctx, url, opts...) (*Response, error)`
- `Options(ctx, url, opts...) (*Response, error)`

## 📝 Logging

El cliente soporta logging opcional de todas las requests y responses.

### Habilitar Logging

```go
import (
    "github.com/foundathyon/base/core/logger"
    "github.com/foundathyon/base/infrastructure/http"
)

config := http.DefaultConfig()
config.Logger = logger.NewSimpleLogger("debug")
client := http.NewClient(config)
```

### Desactivar Logging

```go
config := http.DefaultConfig()
config.Logger = logger.NewSimpleLogger("debug")
config.DisableLogging = true // Desactiva logs incluso con logger configurado
client := http.NewClient(config)
```

### Niveles de Log

El cliente loguea con diferentes niveles según el resultado:

- **Debug**: Inicio de request
- **Info**: Request exitoso (status 2xx)
- **Warn**: Error del cliente (status 4xx)
- **Error**: Error del servidor (status 5xx) o errores de red

### Ejemplo de Logs

```
HTTP request started: method=GET url=https://api.example.com/users timestamp=2026-01-12 10:00:00
HTTP request completed successfully: method=GET url=https://api.example.com/users status_code=200 duration=150ms
```

## 🔄 Retries

El cliente incluye un mecanismo de retry configurable para manejar errores temporales.

### Configuración de Retries

```go
config := http.DefaultConfig()
config.MaxRetries = 3                    // Máximo 3 retries (4 intentos totales)
config.RetryDelay = 200                  // 200ms delay base
config.RetryableStatusCodes = []int{     // Códigos que activan retry
    429, // Too Many Requests
    500, // Internal Server Error
    502, // Bad Gateway
    503, // Service Unavailable
    504, // Gateway Timeout
}
config.RetryableMethods = []string{      // Métodos que se pueden reintentar
    "GET", "HEAD", "OPTIONS", "DELETE",
}
client := http.NewClient(config)
```

### Comportamiento de Retries

- **Errores de red**: Timeouts, errores de conexión, etc. se reintentan automáticamente
- **Status codes retryables**: Los códigos configurados en `RetryableStatusCodes` activan retry
- **Backoff exponencial**: El delay entre retries aumenta exponencialmente (máx. 5 segundos)
- **Métodos idempotentes**: Por defecto solo se reintentan métodos seguros (GET, HEAD, OPTIONS, DELETE)
- **Logging de retries**: Cada intento de retry se loguea con el número de intento y razón

### Ejemplo de Retry

```go
config := http.DefaultConfig()
config.MaxRetries = 3
config.RetryDelay = 100
config.Logger = logger.NewSimpleLogger("debug")
client := http.NewClient(config)

// Si esta request falla con un error de red o status 500,
// se reintentará hasta 3 veces con delays crecientes
resp, err := client.Get(ctx, "https://api.example.com/data")
```

### Logs de Retry

```
HTTP request started: method=GET url=https://api.example.com/data timestamp=...
HTTP request retry: method=GET url=https://api.example.com/data attempt=1/4 delay=100ms reason=network error: timeout
HTTP request retry: method=GET url=https://api.example.com/data attempt=2/4 delay=200ms reason=network error: timeout
HTTP request completed successfully: method=GET url=https://api.example.com/data status_code=200 duration=450ms
```

## 💡 Ejemplos

### Ejemplo 1: Cliente Básico

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/infrastructure/http"
)

func main() {
    ctx := customctx.New(context.Background())
    client := http.NewClient(http.DefaultConfig())

    resp, err := client.Get(ctx, "https://jsonplaceholder.typicode.com/posts/1")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Close()

    body, _ := resp.ReadBodyString()
    fmt.Println(body)
}
```

### Ejemplo 2: Cliente con Logging

```go
package main

import (
    "context"
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/core/logger"
    "github.com/foundathyon/base/infrastructure/http"
)

func main() {
    ctx := customctx.New(context.Background())

    config := http.DefaultConfig()
    config.Logger = logger.NewSimpleLogger("debug")
    client := http.NewClient(config)

    resp, err := client.Get(ctx, "https://api.example.com/data")
    // Los logs se mostrarán automáticamente
}
```

### Ejemplo 3: Cliente con Retries

```go
package main

import (
    "context"
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/core/logger"
    "github.com/foundathyon/base/infrastructure/http"
)

func main() {
    ctx := customctx.New(context.Background())

    config := http.DefaultConfig()
    config.MaxRetries = 3
    config.RetryDelay = 200
    config.Logger = logger.NewSimpleLogger("debug")
    client := http.NewClient(config)

    resp, err := client.Get(ctx, "https://api.example.com/data")
    // Se reintentará automáticamente si falla
}
```

### Ejemplo 4: Cliente para API Externa

```go
type APIClient struct {
    client  *http.Client
    baseURL string
}

func NewAPIClient(baseURL, apiKey string) *APIClient {
    config := http.DefaultConfig()
    config.BaseURL = baseURL
    config.DefaultHeaders = map[string]string{
        "Authorization": fmt.Sprintf("Bearer %s", apiKey),
        "Content-Type":  "application/json",
    }
    config.Timeout = 30
    config.MaxRetries = 3

    return &APIClient{
        client:  http.NewClient(config),
        baseURL: baseURL,
    }
}

func (c *APIClient) GetUser(ctx *customctx.CustomContext, id int) (*User, error) {
    resp, err := c.client.Get(ctx, fmt.Sprintf("/users/%d", id))
    // ...
}
```

## 🎯 Casos de Uso

### Integración con APIs REST Externas

Ideal para integrar con APIs REST externas con soporte para:
- Autenticación mediante headers por defecto
- Retries automáticos para manejar errores temporales
- Logging para debugging y monitoreo

### Microservicios

Útil para comunicación entre microservicios:
- Base URL para servicios internos
- Timeouts configurables
- Retries para resiliencia

### Clientes de API

Perfecto para crear clientes de API estructurados:
- Configuración centralizada
- Headers por defecto para autenticación
- Logging integrado

## 🔗 Referencias

- [Core HTTP Documentation](../core/http/README.md) - Interfaces y tipos base
- [Ejemplos](../../../examples/infrastructure/http/) - Ejemplos completos de uso
