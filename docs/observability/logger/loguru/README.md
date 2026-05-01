# Logger Logrus

`loguru` proporciona una implementación de `logger.ILogger` basada en Logrus con soporte para campos estructurados, hooks personalizados y formateadores personalizados. Es ideal para aplicaciones en producción que requieren logging estructurado avanzado.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Configuración](#configuración)
- [Fields](#fields)
- [Hooks](#hooks)
- [Formatters](#formatters)
- [Ejemplos](#ejemplos)
- [Casos de Uso](#casos-de-uso)
- [Mejores Prácticas](#mejores-prácticas)

## ✨ Características

- **Logging estructurado**: Basado en Logrus con formato personalizado
- **Fields configurables**: Soporte para campos HTTP, WebSocket y personalizados
- **Hooks**: Integración con sistemas externos (Loki, Sentry, etc.)
- **Formateadores personalizados**: Formato específico para diferentes contextos (HTTP, WS)
- **Colores opcionales**: Soporte para salida con colores en desarrollo
- **Niveles configurables**: Control completo sobre los niveles de log
- **Caller tracking**: Captura automática de información del caller

## 📦 Instalación

```bash
go get github.com/foundathyon/base/observability/logger/loguru
```

## 🚀 Uso Básico

```go
import (
    "github.com/foundathyon/base/observability/logger/loguru"
    "github.com/sirupsen/logrus"
)

// Logger básico sin campos
logger := loguru.NewLogger(nil)
logger.Info("Mensaje informativo")
logger.Error("Mensaje de error")

// Logger con configuración personalizada
config := loguru.Config{
    Level:     logrus.InfoLevel,
    Colorable: true,
}
logger2 := loguru.NewLoggerWithConfig(nil, config)
logger2.Info("Mensaje con colores")
```

## 📚 API

### Logger

#### `Logger`

```go
type Logger struct {
    // Campos privados
}
```

`Logger` es una implementación de `logger.ILogger` basada en Logrus que proporciona logging estructurado con soporte para campos, hooks y formateadores personalizados.

**Características:**
- Implementa la interfaz `logger.ILogger`
- Soporte para campos estructurados (HTTP, WebSocket, personalizados)
- Configuración de nivel y colores
- Integración con hooks de Logrus

#### `NewLogger(fields Fields) *Logger`

Crea un nuevo `Logger` con los campos especificados. Usa configuración por defecto (Debug level, sin colores).

**Parámetros:**
- `fields`: Implementación de `Fields` (puede ser `nil`)

**Retorna:**
- `*Logger`: Nueva instancia de Logger

**Ejemplo:**
```go
logger := loguru.NewLogger(nil)
logger := loguru.NewLogger(httpFields)
logger := loguru.NewLogger(wsFields)
```

#### `NewLoggerWithConfig(fields Fields, config Config) *Logger`

Crea un nuevo `Logger` con los campos y configuración especificados.

**Parámetros:**
- `fields`: Implementación de `Fields` (puede ser `nil`)
- `config`: Configuración del logger (nivel y colores)

**Retorna:**
- `*Logger`: Nueva instancia de Logger

**Ejemplo:**
```go
config := loguru.Config{
    Level:     logrus.InfoLevel,
    Colorable: true,
}
logger := loguru.NewLoggerWithConfig(nil, config)
```

#### `SetConfig(config Config)`

Actualiza la configuración del logger (nivel y colores) en tiempo de ejecución.

**Parámetros:**
- `config`: Nueva configuración del logger

**Ejemplo:**
```go
logger.SetConfig(loguru.Config{
    Level:     logrus.WarnLevel,
    Colorable: false,
})
```

#### `AddHook(hook logrus.Hook)`

Agrega un hook de Logrus al logger. Los hooks permiten enviar logs a sistemas externos (Loki, Sentry, etc.).

**Parámetros:**
- `hook`: Hook de Logrus que implementa la interfaz `logrus.Hook`

**Ejemplo:**
```go
lokiHook := hooks.NewLokiBufferedHook(url, batchSize, baseLabels)
logger.AddHook(lokiHook)
```

#### `GetLogrusLogger() *logrus.Logger`

Retorna la instancia subyacente de Logrus. Útil para acceder a funcionalidades avanzadas de Logrus.

**Retorna:**
- `*logrus.Logger`: Instancia de Logrus subyacente

**Ejemplo:**
```go
logrusLogger := logger.GetLogrusLogger()
logrusLogger.SetOutput(customWriter)
```

#### Métodos de Logging

El `Logger` implementa todos los métodos de la interfaz `logger.ILogger`:

```go
Debug(message string, args ...any)
Info(message string, args ...any)
Warn(message string, args ...any)
Error(message string, args ...any)
Fatal(message string, args ...any)
Panic(message string, args ...any)
```

**Características:**
- Todos los métodos aceptan formato con `fmt.Printf`
- Los fields se incluyen automáticamente en el mensaje y en los entries de Logrus
- La información del caller se captura automáticamente

**Ejemplo:**
```go
logger.Info("Usuario autenticado: %s", userID)
logger.Error("Error al procesar solicitud: %v", err)
```

### Config

#### `Config`

```go
type Config struct {
    Level     logrus.Level
    Colorable bool
}
```

`Config` define la configuración del logger.

**Campos:**
- `Level`: Nivel mínimo de logging (logrus.Level)
- `Colorable`: Si es `true`, los logs se mostrarán con colores ANSI

**Valores de Level:**
- `logrus.DebugLevel`: Todos los logs
- `logrus.InfoLevel`: Info, Warn, Error, Fatal, Panic
- `logrus.WarnLevel`: Warn, Error, Fatal, Panic
- `logrus.ErrorLevel`: Error, Fatal, Panic
- `logrus.FatalLevel`: Fatal, Panic
- `logrus.PanicLevel`: Solo Panic

**Ejemplo:**
```go
config := loguru.Config{
    Level:     logrus.InfoLevel,
    Colorable: true,
}
```

### Fields

El logger usa la interfaz `Fields` para manejar campos estructurados. Los campos se pueden usar para:
- Formatear mensajes de log
- Agregar metadata a los entries de Logrus (útil para hooks como Loki)
- Proporcionar contexto estructurado a los logs

#### Interfaz `Fields`

```go
type Fields interface {
    ToFields() map[string]any
    UpdateAll(fields map[string]any)
    UpdateOne(key string, value any)
    Format() string
}
```

**Implementaciones incluidas:**
- `HTTPFileds`: Campos para solicitudes HTTP
- `WSFields`: Campos para conexiones WebSocket

#### HTTPFileds

```go
type HTTPFileds struct {
    Method   string
    Path     string
    TraceID  string
    CallerID string
    ClientIP string
    Metadata map[string]any
}
```

`HTTPFileds` proporciona campos específicos para contexto HTTP.

**Campos:**
- `Method`: Método HTTP (GET, POST, etc.)
- `Path`: Ruta de la solicitud
- `TraceID`: ID de traza
- `CallerID`: ID del llamador
- `ClientIP`: IP del cliente
- `Metadata`: Metadata adicional (mapa abierto)

**Ejemplo:**
```go
httpFields := &fields.HTTPFileds{}
httpFields.UpdateOne("method", "GET")
httpFields.UpdateOne("path", "/api/v1/users")
httpFields.UpdateOne("trace_id", "1234567890")
httpFields.UpdateOne("caller_id", "9876543210")
httpFields.UpdateOne("client_ip", "192.168.1.100")

metadata := map[string]any{
    "user_id": "user-123",
    "role":    "admin",
}
httpFields.UpdateOne("metadata", metadata)

logger := loguru.NewLogger(httpFields)
logger.Info("Solicitud procesada")
```

#### WSFields

```go
type WSFields struct {
    ServerID string
    ClientID string
    TraceID  string
    Path     string
    Metadata map[string]any
}
```

`WSFields` proporciona campos específicos para conexiones WebSocket.

**Campos:**
- `ServerID`: ID del servidor
- `ClientID`: ID del cliente
- `TraceID`: ID de traza
- `Path`: Ruta de la conexión
- `Metadata`: Metadata adicional (mapa abierto)

**Ejemplo:**
```go
wsFields := fields.WSFields{
    ServerID: "server-123",
    ClientID: "client-456",
    TraceID:  "trace-789",
    Path:     "/chat",
    Metadata: map[string]any{
        "room_id": "room-001",
    },
}

logger := loguru.NewLogger(&wsFields)
logger.Info("Conexión WebSocket establecida")
```

### Hooks

Los hooks permiten enviar logs a sistemas externos. El logger incluye un hook para Loki.

#### LokiBufferedHook

```go
type LokiBufferedHook struct {
    URL                  string
    BatchSize            int
    BaseLabels           map[string]string
    Headers              map[string]string
    IncludeCallerAsLabel bool
    IncludeFieldsInLine  bool
}
```

`LokiBufferedHook` es un hook de Logrus que envía logs a Grafana Loki en batches.

**Características:**
- Batching: Agrupa logs antes de enviarlos (configurable)
- Labels automáticos: Los fields se convierten automáticamente en labels de Loki
- Base labels: Labels estáticos para todos los logs (app, environment, etc.)
- Reintentos: Mantiene el buffer si el envío falla

**Campos:**
- `URL`: Endpoint de Loki (ej: `http://localhost:3100/loki/api/v1/push`)
- `BatchSize`: Número de logs por batch (default: 50)
- `BaseLabels`: Labels estáticos (app, environment, service, etc.)
- `Headers`: Headers HTTP opcionales (X-Scope-OrgID, Authorization, etc.)
- `IncludeCallerAsLabel`: Si incluir caller como label (alto riesgo de cardinalidad)
- `IncludeFieldsInLine`: Si incluir fields en la línea de log como JSON

#### `NewLokiBufferedHook(url string, batchSize int, baseLabels map[string]string) *LokiBufferedHook`

Crea un nuevo hook de Loki.

**Parámetros:**
- `url`: URL completa del endpoint de push de Loki
- `batchSize`: Tamaño del batch (si <= 0, usa 50 por defecto)
- `baseLabels`: Labels base para todos los logs

**Retorna:**
- `*LokiBufferedHook`: Nueva instancia del hook

**Ejemplo:**
```go
lokiHook := hooks.NewLokiBufferedHook(
    "http://localhost:3100/loki/api/v1/push",
    10, // batch size
    map[string]string{
        "app":         "my-app",
        "environment": "production",
        "service":     "api",
    },
)

logger.AddHook(lokiHook)
```

#### `Flush() error`

Fuerza el envío de logs pendientes en el buffer.

**Retorna:**
- `error`: Error si el envío falla, `nil` si tiene éxito

**Ejemplo:**
```go
// Al finalizar la aplicación
if err := lokiHook.Flush(); err != nil {
    log.Printf("Error al enviar logs pendientes: %v", err)
}
```

## 💡 Ejemplos

### Ejemplo 1: Logger Básico

```go
package main

import (
    "github.com/foundathyon/base/observability/logger/loguru"
    "github.com/sirupsen/logrus"
)

func main() {
    // Logger básico sin campos
    logger := loguru.NewLogger(nil)
    logger.Debug("Debug message")
    logger.Info("Info message")
    logger.Warn("Warning message")
    logger.Error("Error message")

    // Logger con configuración
    config := loguru.Config{
        Level:     logrus.InfoLevel,
        Colorable: true,
    }
    logger2 := loguru.NewLoggerWithConfig(nil, config)
    logger2.Info("Info with colors")
}
```

### Ejemplo 2: Logger con Campos HTTP

```go
package main

import (
    "github.com/foundathyon/base/observability/logger/loguru"
    "github.com/foundathyon/base/observability/logger/loguru/fields"
)

func main() {
    httpFields := &fields.HTTPFileds{}
    httpFields.UpdateOne("method", "GET")
    httpFields.UpdateOne("path", "/api/v1/users")
    httpFields.UpdateOne("trace_id", "1234567890")
    httpFields.UpdateOne("caller_id", "9876543210")
    httpFields.UpdateOne("client_ip", "192.168.1.100")

    metadata := map[string]any{
        "user_id": "user-123",
        "role":    "admin",
    }
    httpFields.UpdateOne("metadata", metadata)

    logger := loguru.NewLogger(httpFields)
    logger.Info("Usuario autenticado correctamente")
    logger.Info("Solicitud procesada exitosamente")
}
```

### Ejemplo 3: Logger con Campos WebSocket

```go
package main

import (
    "github.com/foundathyon/base/observability/logger/loguru"
    "github.com/foundathyon/base/observability/logger/loguru/fields"
)

func main() {
    wsFields := fields.WSFields{
        ServerID: "server-123",
        ClientID: "client-456",
        TraceID:  "trace-789",
        Path:     "/chat",
        Metadata: map[string]any{
            "room_id": "room-001",
        },
    }

    logger := loguru.NewLogger(&wsFields)
    logger.Info("Conexión WebSocket establecida")
    logger.Info("Mensaje recibido del cliente")
    logger.Error("Error al procesar mensaje WebSocket")
}
```

### Ejemplo 4: Logger con Hook de Loki

```go
package main

import (
    "github.com/foundathyon/base/observability/logger/loguru"
    "github.com/foundathyon/base/observability/logger/loguru/fields"
    "github.com/foundathyon/base/observability/logger/loguru/hooks"
)

func main() {
    httpFields := &fields.HTTPFileds{}
    httpFields.UpdateOne("method", "GET")
    httpFields.UpdateOne("path", "/api/v1/users")
    httpFields.UpdateOne("trace_id", "1234567890")
    httpFields.UpdateOne("client_ip", "192.168.1.100")

    logger := loguru.NewLogger(httpFields)

    // Crear hook de Loki
    lokiHook := hooks.NewLokiBufferedHook(
        "http://localhost:3100/loki/api/v1/push",
        10, // batch size
        map[string]string{
            "app":         "my-app",
            "environment": "production",
        },
    )

    logger.AddHook(lokiHook)

    // Los logs se enviarán a Loki en batches
    logger.Info("Primer log")
    logger.Info("Segundo log")
    // ... más logs
    // Cuando llegue al batch size, se enviará automáticamente

    // Forzar envío de logs pendientes
    lokiHook.Flush()
}
```

### Ejemplo 5: Configuración Avanzada

```go
package main

import (
    "github.com/foundathyon/base/observability/logger/loguru"
    "github.com/foundathyon/base/observability/logger/loguru/fields"
    "github.com/foundathyon/base/observability/logger/loguru/hooks"
    "github.com/sirupsen/logrus"
)

func main() {
    httpFields := &fields.HTTPFileds{}
    httpFields.UpdateOne("method", "POST")
    httpFields.UpdateOne("path", "/api/v1/advanced")
    httpFields.UpdateOne("trace_id", "trace-001")

    // Logger con configuración personalizada
    config := loguru.Config{
        Level:     logrus.DebugLevel,
        Colorable: true,
    }
    logger := loguru.NewLoggerWithConfig(httpFields, config)

    // Hook de Loki
    lokiHook := hooks.NewLokiBufferedHook(
        "http://localhost:3100/loki/api/v1/push",
        5,
        map[string]string{
            "app":         "my-app",
            "environment": "production",
            "version":     "1.0.0",
        },
    )
    logger.AddHook(lokiHook)

    // Cambiar configuración dinámicamente
    logger.SetConfig(loguru.Config{
        Level:     logrus.WarnLevel,
        Colorable: false,
    })

    logger.Debug("No se mostrará")
    logger.Info("No se mostrará")
    logger.Warn("Se mostrará")
    logger.Error("Se mostrará")

    lokiHook.Flush()
}
```

## 🎯 Casos de Uso

### Middleware HTTP

```go
func loggingMiddleware(logger *loguru.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        httpFields := &fields.HTTPFileds{}
        httpFields.UpdateOne("method", c.Request.Method)
        httpFields.UpdateOne("path", c.Request.URL.Path)
        httpFields.UpdateOne("trace_id", c.GetString("trace_id"))
        httpFields.UpdateOne("client_ip", c.ClientIP())

        requestLogger := loguru.NewLogger(httpFields)
        c.Set("logger", requestLogger)

        c.Next()

        requestLogger.Info("Request processed")
    }
}
```

### Conexiones WebSocket

```go
func handleWebSocket(conn *websocket.Conn) {
    wsFields := fields.WSFields{
        ServerID: serverID,
        ClientID: clientID,
        TraceID:  generateTraceID(),
        Path:     conn.Request().URL.Path,
        Metadata: map[string]any{
            "session_id": sessionID,
        },
    }

    logger := loguru.NewLogger(&wsFields)
    logger.Info("WebSocket connection established")

    // ... manejo de conexión

    logger.Info("WebSocket connection closed")
}
```

### Integración con Loki en Producción

```go
func setupProductionLogger() *loguru.Logger {
    httpFields := &fields.HTTPFileds{}

    config := loguru.Config{
        Level:     logrus.InfoLevel,
        Colorable: false, // Sin colores en producción
    }
    logger := loguru.NewLoggerWithConfig(httpFields, config)

    // Hook de Loki con batch grande para producción
    lokiHook := hooks.NewLokiBufferedHook(
        os.Getenv("LOKI_URL"),
        50, // Batch grande para producción
        map[string]string{
            "app":         os.Getenv("APP_NAME"),
            "environment": os.Getenv("ENVIRONMENT"),
            "service":     os.Getenv("SERVICE_NAME"),
            "version":     os.Getenv("APP_VERSION"),
        },
    )

    logger.AddHook(lokiHook)

    return logger
}
```

## 📝 Mejores Prácticas

### Configuración de Niveles

- **Desarrollo**: `logrus.DebugLevel` con colores
- **Staging**: `logrus.InfoLevel` sin colores
- **Producción**: `logrus.InfoLevel` o `logrus.WarnLevel` sin colores

### Uso de Fields

- Usa `HTTPFileds` para solicitudes HTTP
- Usa `WSFields` para conexiones WebSocket
- Actualiza fields dinámicamente según el contexto de la solicitud
- Usa `Metadata` para campos adicionales específicos del dominio

### Hooks de Loki

- **Batch Size**: Usa batches más grandes (50+) en producción para reducir el número de requests
- **Base Labels**: Limita los base labels a valores de baja cardinalidad (app, environment, service)
- **Field Labels**: Los fields se convierten automáticamente en labels - evita campos de alta cardinalidad
- **Flush**: Llama a `Flush()` al finalizar la aplicación para enviar logs pendientes

### Formato de Logs

El logger usa un formato personalizado:
```
DATE_TIME | LEVEL | FILE.FUNCTION:LINE | [fields] | MESSAGE
```

Los campos se formatean según su implementación (`HTTPFileds.Format()`, `WSFields.Format()`, etc.).

### Seguridad

- No incluyas información sensible en los logs (contraseñas, tokens, etc.)
- Usa campos estructurados en lugar de formatear manualmente
- Considera el costo de cardinalidad en Loki labels

## 🔗 Ver También

- [Logger Core](../core/logger/README.md) - Interfaz base del logger
- [Ejemplos](../../../observability/logger/loguru/examples/) - Ejemplos de uso
- [Logrus](https://github.com/sirupsen/logrus) - Biblioteca subyacente
- [Grafana Loki](https://grafana.com/docs/loki/latest/) - Sistema de agregación de logs

## 📚 Referencias

- [Tests](../../../observability/logger/loguru/) - Implementación y tests
- [Mock de Loki](../../../mock/loki/) - Servidor mock para testing
