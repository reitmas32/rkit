# KErrors

`KErrors` proporciona un tipo de error estructurado (`KError`) con información adicional que extiende el manejo de errores estándar de Go. Permite incluir códigos de error, metadata y capacidades de encadenamiento de errores para un mejor seguimiento y depuración.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Ejemplos](#ejemplos)
- [Casos de Uso](#casos-de-uso)
- [Integración con Errores Estándar](#integración-con-errores-estándar)

## ✨ Características

- **Errores estructurados**: Incluye mensaje, código numérico y metadata opcional
- **Encadenamiento de errores**: Soporta wrapping de errores con `Unwrap()`
- **Compatibilidad estándar**: Implementa la interfaz `error` de Go
- **Metadata contextual**: Permite agregar información adicional como request IDs, timestamps, etc.
- **Códigos de error**: Códigos numéricos para manejo programático
- **Integración con errors.Is()**: Compatible con las funciones estándar de Go

## 📦 Instalación

```bash
go get github.com/foundathyon/base/core/kerrors
```

## 🚀 Uso Básico

```go
import "github.com/foundathyon/base/core/kerrors"

// Error básico
err := kerrors.NewKError("Recurso no encontrado", 404, nil)

// Error con metadata
err := kerrors.NewKError("Error de validación", 400, map[string]any{
    "field": "email",
    "value": "invalid",
})

// Error encadenado
rootErr := errors.New("database connection failed")
err := kerrors.NewKErrorWithCause("Error al obtener usuario", 500, nil, rootErr)
```

## 📚 API

### Tipos

#### `KError`

```go
type KError struct {
    Message  string         // Mensaje legible del error
    Code     int            // Código numérico para manejo programático
    Metadata map[string]any // Metadata contextual opcional
    Cause    error          // Error subyacente que causó este error (opcional)
}
```

`KError` representa un error estructurado con mensaje, código de error, metadata opcional y un error subyacente. Implementa la interfaz `error` y soporta encadenamiento de errores.

**Campos:**
- `Message`: Descripción legible del error
- `Code`: Código numérico para categorización y manejo programático
- `Metadata`: Mapa opcional con información contextual adicional
- `Cause`: Error subyacente para encadenamiento (opcional)

### Funciones de Construcción

#### `NewKError(message string, code int, metadata map[string]any) *KError`

Crea un nuevo `KError` con el mensaje, código y metadata especificados. El error retornado no tendrá un error causa configurado.

**Parámetros:**
- `message`: Mensaje descriptivo del error
- `code`: Código numérico del error
- `metadata`: Mapa opcional con información contextual adicional (puede ser `nil`)

**Retorna:**
- `*KError`: Nueva instancia de KError

**Ejemplo:**
```go
err := kerrors.NewKError("Recurso no encontrado", 404, nil)

errWithMetadata := kerrors.NewKError("Error de validación", 400, map[string]any{
    "field":    "email",
    "value":    "invalid-email",
    "expected": "formato válido de email",
})
```

#### `NewKErrorWithCause(message string, code int, metadata map[string]any, cause error) *KError`

Crea un nuevo `KError` con mensaje, código, metadata y un error causa subyacente. Útil para envolver errores existentes con contexto adicional mientras se preserva la cadena de errores original.

**Parámetros:**
- `message`: Mensaje descriptivo del error
- `code`: Código numérico del error
- `metadata`: Mapa opcional con información contextual adicional (puede ser `nil`)
- `cause`: Error subyacente que causó este error

**Retorna:**
- `*KError`: Nueva instancia de KError con error causa

**Ejemplo:**
```go
dbErr := errors.New("database connection timeout")
err := kerrors.NewKErrorWithCause(
    "Error al obtener usuario",
    500,
    map[string]any{
        "operation": "get_user",
        "user_id":   123,
    },
    dbErr,
)
```

### Métodos

#### `Error() string`

Retorna el mensaje del error, implementando la interfaz `error`. Este método permite usar `KError` en cualquier lugar donde se espere un error estándar.

**Retorna:**
- `string`: El mensaje del error

**Ejemplo:**
```go
err := kerrors.NewKError("Error de validación", 400, nil)
fmt.Println(err.Error()) // "Error de validación"

// KError puede usarse como error estándar
var stdErr error = err
fmt.Println(stdErr.Error()) // "Error de validación"
```

#### `Unwrap() error`

Retorna el error subyacente (cause), si existe. Este método permite inspeccionar la cadena de errores usando `errors.Unwrap()` y `errors.Is()`.

**Retorna:**
- `error`: El error subyacente o `nil` si no hay causa

**Ejemplo:**
```go
rootErr := errors.New("root cause")
err := kerrors.NewKErrorWithCause("wrapped error", 500, nil, rootErr)

unwrapped := err.Unwrap()
fmt.Println(unwrapped.Error()) // "root cause"

// También funciona con errors.Unwrap()
unwrapped = errors.Unwrap(err)
fmt.Println(unwrapped.Error()) // "root cause"
```

#### `WithCause(cause error) *KError`

Crea un nuevo `KError` con el mismo mensaje, código y metadata, pero con un error causa diferente. Útil para modificar o agregar una causa a un error existente.

**Parámetros:**
- `cause`: El error causa a asociar

**Retorna:**
- `*KError`: Nueva instancia de KError con el error causa especificado

**Ejemplo:**
```go
err := kerrors.NewKError("Error de validación", 400, map[string]any{"field": "email"})
newCause := errors.New("database error")
errWithCause := err.WithCause(newCause)
```

## 💡 Ejemplos

### Ejemplo 1: Error Básico

```go
package main

import (
    "fmt"
    "github.com/foundathyon/base/core/kerrors"
)

func main() {
    err := kerrors.NewKError("Recurso no encontrado", 404, nil)
    
    fmt.Printf("Mensaje: %s\n", err.Message)
    fmt.Printf("Código: %d\n", err.Code)
    fmt.Printf("Error(): %s\n", err.Error())
}
```

### Ejemplo 2: Error con Metadata

```go
metadata := map[string]any{
    "request_id": "req-abc-123",
    "user_id":    42,
    "timestamp":  "2024-01-01T12:00:00Z",
    "endpoint":   "/api/users",
    "ip_address": "192.168.1.1",
}

err := kerrors.NewKError("Error de validación", 400, metadata)

fmt.Printf("Mensaje: %s\n", err.Message)
fmt.Printf("Código: %d\n", err.Code)
for key, value := range err.Metadata {
    fmt.Printf("  %s: %v\n", key, value)
}
```

### Ejemplo 3: Encadenamiento de Errores

```go
// Error original (raíz)
rootErr := errors.New("error de base de datos: conexión perdida")

// Error intermedio
intermediateErr := kerrors.NewKErrorWithCause(
    "Error al obtener usuario",
    500,
    map[string]any{
        "operation": "get_user",
        "user_id":   123,
    },
    rootErr,
)

// Error de nivel superior
topErr := kerrors.NewKErrorWithCause(
    "Error al procesar solicitud",
    500,
    map[string]any{
        "request_id": "req-xyz",
        "service":    "user_service",
    },
    intermediateErr,
)

// Desenvolver la cadena
current := topErr
level := 1
for current != nil {
    fmt.Printf("Nivel %d: %s (Código: %d)\n", level, current.Message, current.Code)
    
    unwrapped := current.Unwrap()
    if unwrapped == nil {
        break
    }
    
    if kerr, ok := unwrapped.(*kerrors.KError); ok {
        current = kerr
    } else {
        fmt.Printf("Nivel %d (raíz): %s\n", level+1, unwrapped.Error())
        break
    }
    level++
}
```

### Ejemplo 4: Uso con errors.Is() y errors.Unwrap()

```go
rootErr := errors.New("error de red")
wrappedErr := kerrors.NewKErrorWithCause(
    "Error al conectar con API externa",
    503,
    map[string]any{"service": "external_api"},
    rootErr,
)

// Usar errors.Is para verificar si el error raíz está en la cadena
if errors.Is(wrappedErr, rootErr) {
    fmt.Println("✓ errors.Is encontró el error raíz en la cadena")
}

// Usar errors.Unwrap para obtener el error subyacente
unwrapped := errors.Unwrap(wrappedErr)
if unwrapped != nil {
    fmt.Printf("✓ errors.Unwrap retornó: %s\n", unwrapped.Error())
}
```

### Ejemplo 5: Caso de Uso Práctico - Validación

```go
func validateUser(userID int, email string) *kerrors.KError {
    if userID <= 0 {
        return kerrors.NewKError(
            "ID de usuario inválido",
            400,
            map[string]any{
                "field":    "user_id",
                "value":    userID,
                "expected": "número positivo",
            },
        )
    }
    
    if email == "" {
        return kerrors.NewKError(
            "Email es requerido",
            400,
            map[string]any{
                "field": "email",
                "value": email,
            },
        )
    }
    
    // Simular error de base de datos
    dbErr := errors.New("database: connection timeout")
    return kerrors.NewKErrorWithCause(
        "Error al guardar usuario",
        500,
        map[string]any{
            "operation": "save_user",
            "user_id":   userID,
            "email":     email,
        },
        dbErr,
    )
}

// Uso
err := validateUser(42, "user@example.com")
if err != nil {
    fmt.Printf("Error: %s (Código: %d)\n", err.Message, err.Code)
    if err.Metadata != nil {
        fmt.Printf("Metadata: %v\n", err.Metadata)
    }
    if err.Cause != nil {
        fmt.Printf("Causa: %s\n", err.Cause.Error())
    }
}
```

## 🎯 Casos de Uso

### Validación de Datos

`KError` es ideal para validaciones donde necesitas proporcionar información detallada sobre qué falló y por qué:

```go
func validateForm(form FormData) *kerrors.KError {
    if form.Email == "" {
        return kerrors.NewKError(
            "Email es requerido",
            400,
            map[string]any{"field": "email"},
        )
    }
    
    if !isValidEmail(form.Email) {
        return kerrors.NewKError(
            "Email inválido",
            400,
            map[string]any{
                "field":    "email",
                "value":    form.Email,
                "expected": "formato válido de email",
            },
        )
    }
    
    return nil
}
```

### Manejo de Errores en Servicios

Cuando trabajas con servicios externos o APIs, puedes envolver errores con contexto adicional:

```go
func callExternalAPI(ctx context.Context, data Data) *kerrors.KError {
    resp, err := httpClient.Post(ctx, apiURL, data)
    if err != nil {
        return kerrors.NewKErrorWithCause(
            "Error al llamar API externa",
            503,
            map[string]any{
                "api_url":    apiURL,
                "request_id": ctx.Value("request_id"),
            },
            err,
        )
    }
    
    if resp.StatusCode != 200 {
        return kerrors.NewKError(
            "API externa retornó error",
            resp.StatusCode,
            map[string]any{
                "status_code": resp.StatusCode,
                "api_url":     apiURL,
            },
        )
    }
    
    return nil
}
```

### Logging y Monitoreo

La metadata en `KError` es útil para logging estructurado y monitoreo:

```go
err := kerrors.NewKError(
    "Error al procesar pago",
    500,
    map[string]any{
        "request_id":    requestID,
        "user_id":       userID,
        "payment_id":    paymentID,
        "amount":        amount,
        "timestamp":     time.Now(),
        "service":       "payment_service",
    },
)

// Log estructurado
logger.Error("Error procesando pago",
    "error", err.Message,
    "code", err.Code,
    "metadata", err.Metadata,
)
```

## 🔗 Integración con Errores Estándar

`KError` es completamente compatible con el manejo de errores estándar de Go:

### Uso como error estándar

```go
var err error = kerrors.NewKError("Error", 500, nil)
fmt.Println(err.Error()) // "Error"
```

### errors.Is() y errors.Unwrap()

```go
rootErr := errors.New("root error")
wrappedErr := kerrors.NewKErrorWithCause("wrapped", 500, nil, rootErr)

// errors.Is funciona
if errors.Is(wrappedErr, rootErr) {
    fmt.Println("Error raíz encontrado")
}

// errors.Unwrap funciona
unwrapped := errors.Unwrap(wrappedErr)
fmt.Println(unwrapped.Error()) // "root error"
```

### Verificación de tipo

```go
err := someFunction()

// Verificar si es un KError
if kerr, ok := err.(*kerrors.KError); ok {
    fmt.Printf("KError: %s (Código: %d)\n", kerr.Message, kerr.Code)
}

// O usando errors.As
var kerr *kerrors.KError
if errors.As(err, &kerr) {
    fmt.Printf("KError: %s (Código: %d)\n", kerr.Message, kerr.Code)
}
```

## 📝 Mejores Prácticas

### Mensajes de Error

- Usa mensajes descriptivos y específicos
- Evita mensajes genéricos como "Error" o "Algo salió mal"
- Incluye contexto sobre qué operación falló

```go
// ❌ Mal
kerrors.NewKError("Error", 500, nil)

// ✅ Bien
kerrors.NewKError("Error al conectar con base de datos", 500, map[string]any{
    "host": "db.example.com",
    "port": 5432,
})
```

### Códigos de Error

- Usa códigos HTTP estándar cuando sea apropiado (400, 404, 500, etc.)
- Define códigos personalizados para errores específicos del dominio
- Documenta los códigos de error en tu API

```go
const (
    ErrCodeValidation = 400
    ErrCodeNotFound   = 404
    ErrCodeInternal   = 500
    ErrCodeDatabase   = 501
    ErrCodeExternal   = 502
)
```

### Metadata

- Incluye información relevante para depuración
- Usa claves consistentes (ej: `"request_id"`, `"user_id"`, `"timestamp"`)
- Evita incluir información sensible (contraseñas, tokens, etc.)

```go
// ✅ Buen uso de metadata
kerrors.NewKError(
    "Error de validación",
    400,
    map[string]any{
        "field":     "email",
        "value":     email,
        "request_id": requestID,
        "timestamp":  time.Now(),
    },
)
```

### Encadenamiento

- Usa `NewKErrorWithCause` para preservar la cadena de errores
- No pierdas el error original al envolverlo
- Agrega contexto útil en cada nivel

```go
// ✅ Buen encadenamiento
dbErr := database.Query(ctx, query)
if dbErr != nil {
    return kerrors.NewKErrorWithCause(
        "Error al obtener usuario",
        500,
        map[string]any{"user_id": userID},
        dbErr, // Preserva el error original
    )
}
```

## 🔗 Ver También

- [CustomContext](../customctx/README.md) - Contexto personalizado que acumula errores
- [Result](../result/README.md) - Tipo Result para manejo funcional de errores
- [Logger](../logger/README.md) - Sistema de logging

## 📚 Referencias

- [Go Error Handling](https://go.dev/blog/error-handling-and-go)
- [Go errors package](https://pkg.go.dev/errors)
- [Ejemplos de uso](../../../examples/core/kerrors/kerrors_example.go)
- [Tests](../../../core/kerrors/errors_test.go)
