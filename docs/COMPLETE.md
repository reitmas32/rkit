# Base Kit - Complete Documentation

This file contains all documentation from the `docs/` directory combined into a single document.

**Generated on:** 2026-01-13 05:23:41 UTC

---

## Table of Contents

- [core/customctx](#core-customctx)
- [core/eventbus](#core-eventbus)
- [core/http](#core-http)
- [core/kerrors](#core-kerrors)
- [core/logger](#core-logger)
- [core/result](#core-result)
- [infrastructure/eventbus/inmemory](#infrastructure-eventbus-inmemory)
- [infrastructure/eventbus/rabbit](#infrastructure-eventbus-rabbit)
- [infrastructure/eventbus](#infrastructure-eventbus)
- [infrastructure/http](#infrastructure-http)
- [observability/logger/loguru](#observability-logger-loguru)
- [persistence/contracts](#persistence-contracts)
- [persistence/criteria](#persistence-criteria)
- [persistence/inmemory](#persistence-inmemory)
- [persistence/pagination](#persistence-pagination)
- [persistence/postgres](#persistence-postgres)
- [persistence](#persistence)

---


## core/customctx

**Source:** `docs/core/customctx/README.md`

---

# CustomContext

`CustomContext` es una implementación personalizada de `context.Context` que extiende el contexto estándar de Go con la capacidad de acumular errores estructurados (`KError`) a lo largo del flujo de ejecución, sin detener la ejecución. Esto es especialmente útil para recopilar múltiples errores durante el procesamiento y rastrear dónde se registró cada error en la pila de llamadas.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Ejemplos](#ejemplos)
- [Casos de Uso](#casos-de-uso)
- [Almacenamiento de Valores](#almacenamiento-de-valores)
- [Thread Safety](#thread-safety)
- [Integración con Context Estándar](#integración-con-context-estándar)

## ✨ Características

- **Acumulación de errores**: Permite recolectar múltiples errores estructurados sin detener la ejecución
- **Tracking de call sites**: Captura automáticamente la ubicación donde se registró cada error (función y línea)
- **Thread-safe**: Todos los métodos son seguros para uso concurrente
- **Compatible con context.Context**: Implementa completamente la interfaz estándar de Go
- **Integración transparente**: Funciona con todos los métodos del context estándar (timeout, cancelación, valores, etc.)
- **Gestión de logger**: Soporte opcional para inyección de logger
- **Almacenamiento de valores**: Permite almacenar metadata técnica (request IDs, trace IDs, etc.) usando strings como claves

## 📦 Instalación

```bash
go get github.com/foundathyon/base/core/customctx
```

## 🚀 Uso Básico

```go
import (
    "context"
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/core/kerrors"
)

// Crear un CustomContext
ctx := customctx.New(context.Background())

// Agregar errores
ctx.AddError(kerrors.NewKError("Error de validación", 400, map[string]any{
    "field": "email",
}))

// Verificar si hay errores
if ctx.HasErrors() {
    errors := ctx.Errors()
    firstErr := ctx.FirstError()
    lastErr := ctx.LastError()
}

// Almacenar metadata técnica
ctx = ctx.WithValue("request_id", "req-123")
requestID := ctx.GetValue("request_id")
```

## 📚 API

### Tipos

#### `CustomContext`

```go
type CustomContext struct {
    parent context.Context
    mu     sync.Mutex
    errors []WrapError
    logger logger.ILogger
}
```

`CustomContext` es una implementación de `context.Context` que acumula errores estructurados. Envuelve un contexto padre y mantiene toda la funcionalidad estándar del contexto mientras añade la capacidad de recolectar múltiples errores durante la ejecución.

**Características:**
- Todos los métodos son thread-safe
- Implementa completamente la interfaz `context.Context`
- Delega todas las operaciones estándar al contexto padre

#### `WrapError`

```go
type WrapError struct {
    Error  *kerrors.KError `json:"error"`   // Error estructurado registrado
    CallIn string          `json:"call_in"` // Ubicación donde se registró: "functionName:lineNumber"
}
```

`WrapError` asocia un error estructurado con la ubicación donde fue registrado. Contiene tanto el `KError` como la información del call site (nombre de función y número de línea).

### Funciones de Construcción

#### `New(parent context.Context) *CustomContext`

Crea un nuevo `CustomContext` que envuelve el contexto padre.

**Parámetros:**
- `parent`: El contexto padre a envolver

**Retorna:**
- `*CustomContext`: Una nueva instancia de CustomContext

**Ejemplo:**
```go
parent := context.Background()
ctx := customctx.New(parent)
```

#### `NewCustomContext(parent context.Context) *CustomContext`

Alias de `New` para compatibilidad. **Deprecated**: usar `New` en su lugar.

### Métodos de Context.Context

`CustomContext` implementa completamente la interfaz `context.Context`, delegando todas las operaciones al contexto padre:

#### `Deadline() (time.Time, bool)`

Retorna el tiempo límite para cancelar el trabajo realizado en nombre de este contexto. Delega al método `Deadline` del contexto padre.

#### `Done() <-chan struct{}`

Retorna un canal que se cierra cuando el trabajo realizado en nombre de este contexto debe ser cancelado. Delega al contexto padre.

#### `Err() error`

Retorna `nil` si `Done` no está cerrado. Si `Done` está cerrado, retorna el error que explica por qué. Delega al contexto padre.

#### `Value(key interface{}) interface{}`

Retorna el valor asociado con este contexto para la clave, o `nil` si no hay valor asociado. Delega al contexto padre.

### Métodos de Acceso al Contexto Padre

#### `Context() context.Context`

Retorna el `context.Context` padre que este `CustomContext` envuelve. Esto permite acceder al contexto subyacente para operaciones que requieren la interfaz estándar de contexto.

**Ejemplo:**
```go
parent := context.WithValue(context.Background(), "key", "value")
ctx := customctx.New(parent)
originalCtx := ctx.Context() // Retorna el contexto padre
```

### Métodos de Gestión de Errores

#### `AddError(err *kerrors.KError) *kerrors.KError`

Registra un error estructurado en el contexto, capturando automáticamente la información del caller (nombre de función y número de línea) donde se llamó este método.

**Parámetros:**
- `err`: El error estructurado a registrar. Si es `nil`, el método no hace nada.

**Retorna:**
- `*kerrors.KError`: El mismo error que se registró (o `nil` si el parámetro era `nil`)

**Características:**
- Thread-safe
- Si `err` es `nil`, no hace nada y retorna `nil`
- Captura automáticamente el call site donde se llamó

**Ejemplo:**
```go
err := kerrors.NewKError("Error de validación", 400, map[string]any{
    "field": "email",
})
ctx.AddError(err) // El error se registra con información del call site
```

#### `NewError(err *kerrors.KError) *kerrors.KError`

Alias de `AddError` para compatibilidad. **Deprecated**: usar `AddError` en su lugar.

#### `Errors() []WrapError`

Retorna todos los errores acumulados con sus call sites asociados. El slice retornado contiene todos los errores que han sido registrados via `AddError`, ordenados por la secuencia en que fueron agregados.

**Retorna:**
- `[]WrapError`: Slice con todos los errores acumulados

**Nota:** El slice retornado es el slice interno, por lo que modificarlo puede afectar el estado interno del contexto. Considera hacer una copia si necesitas modificar el slice.

**Ejemplo:**
```go
errors := ctx.Errors()
for _, wrapErr := range errors {
    fmt.Printf("Error: %s, Registrado en: %s\n", 
        wrapErr.Error.Message, 
        wrapErr.CallIn)
}
```

#### `HasErrors() bool`

Retorna `true` si hay al menos un error registrado en el contexto. Este método es thread-safe y puede ser usado para verificar si hay errores antes de acceder a los detalles de los errores.

**Retorna:**
- `bool`: `true` si hay errores, `false` en caso contrario

**Ejemplo:**
```go
if ctx.HasErrors() {
    // Procesar errores
    errors := ctx.Errors()
}
```

#### `FirstError() WrapError`

Retorna el primer error registrado (el error más antiguo en la colección). Este método hace panic si no hay errores registrados. Usa `HasErrors()` para verificar antes de llamar a este método.

**Retorna:**
- `WrapError`: El primer error registrado

**Panics:**
- Si no hay errores registrados

**Ejemplo:**
```go
if ctx.HasErrors() {
    firstErr := ctx.FirstError()
    fmt.Printf("Primer error: %s\n", firstErr.Error.Message)
}
```

#### `LastError() WrapError`

Retorna el último error registrado (el error más reciente en la colección). Este método hace panic si no hay errores registrados. Usa `HasErrors()` para verificar antes de llamar a este método.

**Retorna:**
- `WrapError`: El último error registrado

**Panics:**
- Si no hay errores registrados

**Ejemplo:**
```go
if ctx.HasErrors() {
    lastErr := ctx.LastError()
    fmt.Printf("Último error: %s\n", lastErr.Error.Message)
}
```

#### `Clear()`

Remueve todos los errores acumulados del contexto, reiniciando la colección de errores. Después de llamar a `Clear()`, `HasErrors()` retornará `false` y `Errors()` retornará un slice vacío. Este método es thread-safe.

**Ejemplo:**
```go
ctx.AddError(kerrors.NewKError("Error 1", 500, nil))
ctx.AddError(kerrors.NewKError("Error 2", 400, nil))
ctx.Clear() // Limpia todos los errores
if !ctx.HasErrors() {
    fmt.Println("No hay errores después de Clear()")
}
```

### Métodos de Logger

#### `WithLogger(logger logger.ILogger) *CustomContext`

Crea un nuevo `CustomContext` con el logger especificado. El contexto original no se modifica.

**Parámetros:**
- `logger`: El logger a asociar con el contexto

**Retorna:**
- `*CustomContext`: Nueva instancia con el logger configurado

**Ejemplo:**
```go
logger := logger.NewSimpleLogger(logger.LevelInfo.String())
ctxWithLogger := ctx.WithLogger(logger)
```

#### `Logger() logger.ILogger`

Retorna el logger asociado con el contexto. Si no hay logger configurado, retorna un logger simple con nivel Debug por defecto.

**Retorna:**
- `logger.ILogger`: El logger configurado o un logger por defecto

**Ejemplo:**
```go
logger := ctx.Logger()
logger.Info("Mensaje de log")
```

### Métodos de Almacenamiento de Valores (Metadata Técnica)

#### `WithValue(key string, val any) *CustomContext`

Retorna un nuevo `CustomContext` con el par clave-valor agregado. El contexto retornado es inmutable - el contexto original no se modifica.

**IMPORTANTE:** Solo debe usarse para metadata técnica request-scoped:
- Request IDs, trace IDs, correlation IDs
- User IDs (para propósitos técnicos como logging/tracing, no lógica de negocio)
- Nombres de servicios, nombres de operaciones
- Flags técnicos o configuración

**NO usar para datos de negocio:**
- Entidades de dominio (usuarios, órdenes, productos, etc.)
- Estado o contexto de negocio
- Valores específicos del dominio de la aplicación

Los datos de negocio deben pasarse como parámetros explícitos de función, no almacenarse en el contexto.

**Parámetros:**
- `key`: Clave string descriptiva (ej: `"request_id"`, `"trace_id"`)
- `val`: Valor de cualquier tipo (`any`)

**Retorna:**
- `*CustomContext`: Nueva instancia con el valor agregado

**Características:**
- Thread-safe
- Inmutable (retorna nuevo contexto)
- Comparte parent context, errores y logger, pero tiene su propio mapa de valores

**Ejemplo:**
```go
ctx := customctx.New(context.Background())
ctx = ctx.WithValue("request_id", "req-123")
ctx = ctx.WithValue("trace_id", "trace-456")
ctx = ctx.WithValue("user_id", "user-789") // ID técnico para logging
```

#### `GetValue(key string) any`

Retorna el valor asociado con la clave string dada en este `CustomContext` únicamente, sin verificar el contexto padre. Retorna `nil` si la clave no se encuentra.

Este es un método de conveniencia equivalente a verificar si la clave existe en el mapa de valores interno del `CustomContext`. Para verificar tanto `CustomContext` como el contexto padre, usa `Value()` en su lugar.

**Parámetros:**
- `key`: Clave string a buscar

**Retorna:**
- `any`: El valor asociado o `nil` si no se encuentra

**Características:**
- Thread-safe
- Solo busca en el `CustomContext`, no en el parent

**Ejemplo:**
```go
requestID := ctx.GetValue("request_id")
if requestID != nil {
    fmt.Printf("Request ID: %v\n", requestID)
}
```

#### `HasValue(key string) bool`

Retorna `true` si la clave string dada existe en el mapa de valores de este `CustomContext`, sin verificar el contexto padre.

Este es un método de conveniencia para verificar la existencia de una clave. Para verificar tanto `CustomContext` como el contexto padre, usa `Value()` y verifica por `nil`.

**Parámetros:**
- `key`: Clave string a verificar

**Retorna:**
- `bool`: `true` si la clave existe, `false` en caso contrario

**Características:**
- Thread-safe
- Solo verifica en el `CustomContext`, no en el parent

**Ejemplo:**
```go
if ctx.HasValue("request_id") {
    fmt.Println("Request ID está presente")
}
```

#### `Value(key interface{}) interface{}` (Mejorado)

Retorna el valor asociado con este contexto para la clave, o `nil` si no hay valor asociado. Primero verifica el mapa de valores interno del `CustomContext` (si la clave es un string), luego delega al método `Value` del contexto padre.

**Parámetros:**
- `key`: Clave a buscar (puede ser `interface{}` para compatibilidad con `context.Context`)

**Retorna:**
- `interface{}`: El valor asociado o `nil`

**Comportamiento:**
- Si `key` es un `string`, busca primero en el `CustomContext`
- Si `key` no es un `string` o no se encuentra en `CustomContext`, delega al contexto padre
- Mantiene compatibilidad total con `context.Context`

**Ejemplo:**
```go
// String key - busca primero en CustomContext
requestID := ctx.Value("request_id")

// Non-string key - delega directamente al parent
userID := ctx.Value(contextUserIDKey{})
```

## 💡 Ejemplos

### Ejemplo 1: Uso Básico

```go
package main

import (
    "context"
    "fmt"
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/core/kerrors"
)

func main() {
    ctx := customctx.New(context.Background())
    
    // Agregar un error
    err := kerrors.NewKError("Error de configuración", 500, map[string]any{
        "config_key": "database_url",
    })
    ctx.AddError(err)
    
    // Verificar errores
    if ctx.HasErrors() {
        fmt.Printf("Errores encontrados: %d\n", len(ctx.Errors()))
        firstErr := ctx.FirstError()
        fmt.Printf("Primer error: %s\n", firstErr.Error.Message)
        fmt.Printf("Registrado en: %s\n", firstErr.CallIn)
    }
}
```

### Ejemplo 2: Acumulación de Múltiples Errores

```go
func validateForm(ctx *customctx.CustomContext, form FormData) {
    if form.Nombre == "" {
        ctx.AddError(kerrors.NewKError(
            "El nombre es requerido",
            400,
            map[string]any{"field": "nombre"},
        ))
    }
    
    if form.Email == "" {
        ctx.AddError(kerrors.NewKError(
            "El email es requerido",
            400,
            map[string]any{"field": "email"},
        ))
    } else if !isValidEmail(form.Email) {
        ctx.AddError(kerrors.NewKError(
            "El email tiene formato inválido",
            400,
            map[string]any{
                "field": "email",
                "value": form.Email,
            },
        ))
    }
    
    if form.Edad < 18 {
        ctx.AddError(kerrors.NewKError(
            "Debe ser mayor de edad",
            400,
            map[string]any{
                "field": "edad",
                "value": form.Edad,
                "min":   18,
            },
        ))
    }
}

// Uso
ctx := customctx.New(context.Background())
validateForm(ctx, formData)

if ctx.HasErrors() {
    fmt.Printf("Errores de validación encontrados:\n")
    for i, wrapErr := range ctx.Errors() {
        fmt.Printf("%d. %s (Código: %d)\n", 
            i+1, 
            wrapErr.Error.Message, 
            wrapErr.Error.Code)
        fmt.Printf("   Registrado en: %s\n", wrapErr.CallIn)
    }
}
```

### Ejemplo 3: Integración con Context Estándar

```go
// Crear un context con timeout
parent, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

// Crear un context con un valor
parent = context.WithValue(parent, "user_id", 123)

ctx := customctx.New(parent)

// Acceder a valores del context padre
userID := ctx.Value("user_id")
fmt.Printf("User ID: %v\n", userID)

// Verificar deadline
deadline, ok := ctx.Deadline()
if ok {
    fmt.Printf("Deadline: %v\n", deadline)
}

// Agregar errores mientras se usa el context
ctx.AddError(kerrors.NewKError("Error al procesar", 500, nil))

// El context funciona normalmente con cancelación
select {
case <-ctx.Done():
    fmt.Println("Context cancelado o timeout")
case <-time.After(1 * time.Second):
    fmt.Println("Todavía ejecutando...")
}
```

### Ejemplo 4: Almacenamiento de Valores

```go
ctx := customctx.New(context.Background())

// Almacenar metadata técnica
ctx = ctx.WithValue("request_id", "req-abc-123")
ctx = ctx.WithValue("trace_id", "trace-xyz-456")
ctx = ctx.WithValue("user_id", "user-789")

// Recuperar valores
fmt.Printf("Request ID: %v\n", ctx.GetValue("request_id"))
fmt.Printf("Trace ID: %v\n", ctx.GetValue("trace_id"))

// Verificar existencia
if ctx.HasValue("request_id") {
    fmt.Println("Request ID está presente")
}

// Usar Value() que busca en CustomContext y parent
requestID := ctx.Value("request_id")

// Inmutabilidad
ctxOriginal := customctx.New(context.Background())
ctx1 := ctxOriginal.WithValue("request_id", "req-123")
ctx2 := ctx1.WithValue("trace_id", "trace-456")

fmt.Printf("Original: %v\n", ctxOriginal.GetValue("request_id")) // nil
fmt.Printf("ctx1: %v\n", ctx1.GetValue("request_id"))             // "req-123"
fmt.Printf("ctx2: %v, %v\n", ctx2.GetValue("request_id"), ctx2.GetValue("trace_id"))
```

### Ejemplo 5: Uso Concurrente

```go
func processConcurrently(ctx *customctx.CustomContext) {
    // Simular múltiples goroutines agregando errores
    done := make(chan bool, 5)
    for i := 0; i < 5; i++ {
        go func(id int) {
            // Simular algún procesamiento que puede fallar
            if id%2 == 0 {
                ctx.AddError(kerrors.NewKError(
                    fmt.Sprintf("Error de goroutine %d", id),
                    500,
                    map[string]any{"goroutine_id": id},
                ))
            }
            done <- true
        }(i)
    }
    
    // Esperar a que todas las goroutines terminen
    for i := 0; i < 5; i++ {
        <-done
    }
    
    fmt.Printf("Total de errores (concurrentes): %d\n", len(ctx.Errors()))
}

ctx := customctx.New(context.Background())
processConcurrently(ctx)
```

## 🎯 Casos de Uso

### Validación de Formularios

`CustomContext` es ideal para validar formularios donde necesitas recopilar todos los errores de validación antes de responder al usuario:

```go
func validateRegistrationForm(ctx *customctx.CustomContext, form RegistrationForm) {
    validateEmail(ctx, form.Email)
    validatePassword(ctx, form.Password)
    validateAge(ctx, form.Age)
    validateTerms(ctx, form.AcceptedTerms)
}

// Al final, puedes retornar todos los errores juntos
if ctx.HasErrors() {
    return respondWithValidationErrors(ctx.Errors())
}
```

### Procesamiento en Lotes

Cuando procesas múltiples elementos y quieres continuar incluso si algunos fallan:

```go
func processBatch(ctx *customctx.CustomContext, items []Item) {
    for _, item := range items {
        if err := processItem(ctx, item); err != nil {
            ctx.AddError(kerrors.NewKError(
                fmt.Sprintf("Error procesando item %d", item.ID),
                500,
                map[string]any{"item_id": item.ID},
            ))
        }
    }
}

// Al final, puedes reportar todos los fallos
if ctx.HasErrors() {
    logFailedItems(ctx.Errors())
}
```

### Middleware de Errores en HTTP

En un middleware HTTP, puedes acumular errores durante el procesamiento de la request:

```go
func errorAccumulatingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := customctx.New(r.Context())
        r = r.WithContext(ctx)
        
        next.ServeHTTP(w, r)
        
        // Después del handler, verificar si hubo errores
        if ctx.HasErrors() {
            respondWithErrors(w, ctx.Errors())
        }
    })
}
```

## 🔒 Thread Safety

Todos los métodos de `CustomContext` son **thread-safe** y pueden ser llamados desde múltiples goroutines simultáneamente sin necesidad de sincronización adicional. Esto incluye:

- `AddError()` - Puede ser llamado concurrentemente
- `Errors()` - Puede ser llamado concurrentemente
- `HasErrors()` - Puede ser llamado concurrentemente
- `FirstError()` / `LastError()` - Puede ser llamado concurrentemente
- `Clear()` - Puede ser llamado concurrentemente

La sincronización interna está garantizada mediante un `sync.Mutex`.

**Nota importante sobre `Errors()`**: El slice retornado es una referencia al slice interno. Si necesitas modificar el slice o iterar sobre él mientras otras goroutines pueden estar agregando errores, debes hacer una copia:

```go
// Copia segura
errors := make([]customctx.WrapError, len(ctx.Errors()))
copy(errors, ctx.Errors())
```

## 🔑 Almacenamiento de Valores

`CustomContext` permite almacenar valores usando strings como claves, similar a `context.WithValue()`, pero con su propio almacenamiento interno.

### Uso Básico

```go
ctx := customctx.New(context.Background())

// Almacenar valores (solo metadata técnica)
ctx = ctx.WithValue("request_id", "req-abc-123")
ctx = ctx.WithValue("trace_id", "trace-xyz-456")
ctx = ctx.WithValue("user_id", "user-789") // ID técnico para logging

// Recuperar valores
requestID := ctx.GetValue("request_id")
if ctx.HasValue("request_id") {
    fmt.Printf("Request ID: %v\n", requestID)
}

// Usar Value() que busca en CustomContext y parent
requestID = ctx.Value("request_id")
```

### Inmutabilidad

Al igual que `context.WithValue()`, `WithValue()` retorna un nuevo contexto sin modificar el original:

```go
ctxOriginal := customctx.New(context.Background())
ctx1 := ctxOriginal.WithValue("request_id", "req-123")
ctx2 := ctx1.WithValue("trace_id", "trace-456")

// ctxOriginal no tiene valores
fmt.Println(ctxOriginal.GetValue("request_id")) // nil

// ctx1 tiene request_id
fmt.Println(ctx1.GetValue("request_id")) // "req-123"

// ctx2 tiene ambos
fmt.Println(ctx2.GetValue("request_id")) // "req-123"
fmt.Println(ctx2.GetValue("trace_id"))   // "trace-456"
```

### Integración con Parent Context

El método `Value()` busca primero en el `CustomContext` (para string keys), luego en el contexto padre:

```go
// Valor en parent context
parent := context.WithValue(context.Background(), "parent_key", "parent_value")

// Valor en CustomContext
ctx := customctx.New(parent)
ctx = ctx.WithValue("custom_key", "custom_value")

// Value() busca primero en CustomContext (para strings), luego en parent
fmt.Println(ctx.Value("custom_key"))  // "custom_value" (de CustomContext)
fmt.Println(ctx.Value("parent_key"))  // "parent_value" (de parent)
```

### ⚠️ Importante: Solo Metadata Técnica

**Usar para:**
- ✅ Request IDs, trace IDs, correlation IDs
- ✅ User IDs técnicos (para logging/tracing)
- ✅ Nombres de servicios u operaciones
- ✅ Flags técnicos o configuración

**NO usar para:**
- ❌ Entidades de dominio (usuarios, órdenes, productos)
- ❌ Estado o contexto de negocio
- ❌ Valores específicos del dominio

Los datos de negocio deben pasarse como parámetros explícitos de función.

## 🔗 Integración con Context Estándar

`CustomContext` es completamente compatible con `context.Context`. Todas las operaciones estándar funcionan normalmente:

- **Timeouts**: `context.WithTimeout()` funciona normalmente
- **Cancelación**: `context.WithCancel()` funciona normalmente
- **Valores**: `context.WithValue()` funciona normalmente
- **Deadline**: Se propaga desde el contexto padre
- **Done channel**: Se propaga desde el contexto padre

```go
// Todas estas operaciones funcionan normalmente
parent, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

parent = context.WithValue(parent, "request_id", "abc-123")

ctx := customctx.New(parent)

// Acceder a valores
requestID := ctx.Value("request_id")

// Verificar deadline
deadline, _ := ctx.Deadline()

// Usar cancelación
select {
case <-ctx.Done():
    // Context cancelado
default:
    // Continuar procesamiento
}
```

## 📝 Notas Adicionales

### Captura de Call Sites

La información del call site se captura usando `runtime.Caller(1)`, que obtiene información de la función que llama a `AddError()`. El formato es:

```
functionName:lineNumber
```

Por ejemplo: `github.com/user/project.validateForm:45`

### Rendimiento

`CustomContext` está diseñado para ser eficiente. Las operaciones principales tienen complejidad:

- `AddError()`: O(1) amortizado
- `Errors()`: O(1) (retorna referencia al slice)
- `HasErrors()`: O(1)
- `FirstError()` / `LastError()`: O(1)
- `Clear()`: O(1)

### Migración desde Deprecated

Si estás usando `NewCustomContext()` o `NewError()`, migra a:

- `NewCustomContext()` → `New()`
- `NewError()` → `AddError()`

Estos métodos deprecated seguirán funcionando pero pueden ser removidos en versiones futuras.

## 🔗 Ver También

- [KError](../kerrors/README.md) - Documentación del tipo de error estructurado
- [Result](../result/README.md) - Documentación del tipo Result genérico
- [Logger](../logger/README.md) - Documentación del sistema de logging

## 📚 Referencias

- [Go Context Package](https://pkg.go.dev/context)
- [Ejemplos de uso](../../../examples/core/customctx/customctx_example.go)
- [Tests](../../../core/customctx/customctx_test.go)

---


## core/eventbus

**Source:** `docs/core/eventbus/README.md`

---

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

---


## core/http

**Source:** `docs/core/http/README.md`

---

# HTTP Client - Core

El módulo `core/http` proporciona abstracciones e interfaces para realizar peticiones HTTP de forma tipada y segura. Sigue los principios de Clean Architecture, separando las interfaces del core de las implementaciones concretas.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Tipos](#tipos)
- [Ejemplos](#ejemplos)
- [Casos de Uso](#casos-de-uso)

## ✨ Características

- **Interfaz abstracta**: `Client` interface que permite diferentes implementaciones
- **Respuestas tipadas**: Soporte para respuestas genéricas con `TypedResponse[T]`
- **Request/Response tipados**: Objetos estructurados para requests y responses
- **Serialización automática**: Marshal automático de objetos a JSON en requests
- **Funcional options**: Patrón de opciones funcionales para configuración
- **Criterios de éxito configurables**: Define qué códigos de estado se consideran exitosos
- **Medición de tiempo**: Tracking automático de duración de requests
- **Type-safe**: Uso de genéricos para seguridad de tipos

## 📦 Instalación

```bash
go get github.com/foundathyon/base/core/http
```

## 🚀 Uso Básico

### Request Simple

```go
import (
    "github.com/foundathyon/base/core/customctx"
    corehttp "github.com/foundathyon/base/core/http"
    "github.com/foundathyon/base/infrastructure/http"
)

ctx := customctx.New(context.Background())
client := http.NewClient(http.DefaultConfig())

// GET request simple
resp, err := client.Get(ctx, "https://api.example.com/users")
if err != nil {
    log.Fatal(err)
}
defer resp.Close()
```

### Request con Respuesta Tipada

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// GET con respuesta tipada
resp, err := corehttp.GetTyped[User](
    client,
    ctx,
    "https://api.example.com/users/1",
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("User: %+v\n", resp.Body)
fmt.Printf("Status: %d\n", resp.StatusCode)
fmt.Printf("Duration: %v\n", resp.Duration)
```

### POST con Body Tipado

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

createReq := CreateUserRequest{
    Name:  "John Doe",
    Email: "john@example.com",
}

// POST con request y response tipados
// TRequest = CreateUserRequest, TResponse = UserResponse
resp, err := corehttp.PostTyped[CreateUserRequest, UserResponse](
    client,
    ctx,
    "https://api.example.com/users",
    createReq, // El objeto se serializa automáticamente a JSON
    corehttp.WithContentType("application/json"),
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created User: %+v\n", resp.Body)
```

## 📚 API

### Client Interface

La interfaz `Client` define los métodos para realizar peticiones HTTP:

```go
type Client interface {
    Do(ctx *customctx.CustomContext, req *Request) (*Response, error)
    Get(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)
    Post(ctx *customctx.CustomContext, url string, body io.Reader, opts ...RequestOption) (*Response, error)
    Put(ctx *customctx.CustomContext, url string, body io.Reader, opts ...RequestOption) (*Response, error)
    Patch(ctx *customctx.CustomContext, url string, body io.Reader, opts ...RequestOption) (*Response, error)
    Delete(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)
    Head(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)
    Options(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)
}
```

### Request

Estructura que representa una petición HTTP:

```go
type Request struct {
    Method      string
    URL         string
    Headers     map[string]string
    Body        io.Reader
    ContentType string
    QueryParams map[string]string
    Timeout     *int // en segundos
}
```

### Response

Estructura que representa una respuesta HTTP:

```go
type Response struct {
    StatusCode    int
    Status        string
    Headers       map[string]string
    Body          io.ReadCloser
    ContentLength int64
}
```

### TypedResponse

Wrapper genérico para respuestas tipadas:

```go
type TypedResponse[T any] struct {
    StatusCode              int
    Status                  string
    Headers                 map[string]string
    Body                    T              // Respuesta parseada del tipo T
    RawBody                 []byte         // Body crudo (útil para binarios)
    ContentType             string
    ContentLength           int64
    RequestTime             time.Time
    ResponseTime            time.Time
    Duration                time.Duration
    ExpectedStatusCode      *int
    SuccessStatusCodeRange  *StatusCodeRange
}
```

### Funciones Helper Tipadas

Funciones genéricas para realizar peticiones con respuestas tipadas:

```go
// GET con respuesta tipada
func GetTyped[T any](client Client, ctx *customctx.CustomContext, url string, opts ...RequestOption) (*TypedResponse[T], error)

// POST con request y response tipados
func PostTyped[TRequest any, TResponse any](client Client, ctx *customctx.CustomContext, url string, body TRequest, opts ...RequestOption) (*TypedResponse[TResponse], error)

// PUT con request y response tipados
func PutTyped[TRequest any, TResponse any](client Client, ctx *customctx.CustomContext, url string, body TRequest, opts ...RequestOption) (*TypedResponse[TResponse], error)

// PATCH con request y response tipados
func PatchTyped[TRequest any, TResponse any](client Client, ctx *customctx.CustomContext, url string, body TRequest, opts ...RequestOption) (*TypedResponse[TResponse], error)

// DELETE con respuesta tipada
func DeleteTyped[T any](client Client, ctx *customctx.CustomContext, url string, opts ...RequestOption) (*TypedResponse[T], error)

// HEAD con respuesta tipada
func HeadTyped[T any](client Client, ctx *customctx.CustomContext, url string, opts ...RequestOption) (*TypedResponse[T], error)

// OPTIONS con respuesta tipada
func OptionsTyped[T any](client Client, ctx *customctx.CustomContext, url string, opts ...RequestOption) (*TypedResponse[T], error)
```

### Request Options

Opciones funcionales para configurar requests:

```go
// Añadir un header
WithHeader(key, value string) RequestOption

// Añadir múltiples headers
WithHeaders(headers map[string]string) RequestOption

// Añadir un query parameter
WithQueryParam(key, value string) RequestOption

// Añadir múltiples query parameters
WithQueryParams(params map[string]string) RequestOption

// Establecer Content-Type
WithContentType(contentType string) RequestOption

// Establecer timeout en segundos
WithTimeout(seconds int) RequestOption
```

### TypedResponse Options

Opciones para configurar criterios de éxito:

```go
// Establecer código de estado esperado
WithExpectedStatusCode[T](code int) func(*TypedResponse[T])

// Establecer rango de códigos de estado exitosos
WithSuccessStatusCodeRange[T](min, max int) func(*TypedResponse[T])
```

### RequestBody Interface

Interfaz para tipos que pueden serializarse a request body:

```go
type RequestBody interface {
    ToReader() (io.Reader, error)
}
```

Si un tipo implementa `RequestBody`, se usará su método `ToReader()`. Si no, se intentará hacer JSON marshal automático.

## 📖 Tipos

### RequestBody

Interfaz para serialización personalizada de request bodies:

```go
type RequestBody interface {
    ToReader() (io.Reader, error)
}
```

**Implementaciones incluidas:**
- `NewJSONRequestBody(value any) RequestBody` - Serializa cualquier valor a JSON
- `NewReaderRequestBody(reader io.Reader) RequestBody` - Envuelve un `io.Reader` existente

### StatusCodeRange

Rango de códigos de estado HTTP:

```go
type StatusCodeRange struct {
    Min int // Mínimo (inclusivo)
    Max int // Máximo (inclusivo)
}
```

## 💡 Ejemplos

### Ejemplo 1: GET Simple

```go
ctx := customctx.New(context.Background())
client := http.NewClient(http.DefaultConfig())

resp, err := client.Get(ctx, "https://api.example.com/data")
if err != nil {
    log.Fatal(err)
}
defer resp.Close()

body, _ := resp.ReadBodyString()
fmt.Println(body)
```

### Ejemplo 2: GET con Respuesta Tipada

```go
type Product struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Price float64 `json:"price"`
}

resp, err := corehttp.GetTyped[Product](
    client,
    ctx,
    "https://api.example.com/products/1",
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Product: %+v\n", resp.Body)
fmt.Printf("Is Success: %v\n", resp.IsSuccess())
```

### Ejemplo 3: POST con Body Automático

```go
type CreateOrder struct {
    UserID    int    `json:"user_id"`
    ProductID int    `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

type Order struct {
    ID       int    `json:"id"`
    UserID   int    `json:"user_id"`
    Total    float64 `json:"total"`
    Status   string  `json:"status"`
}

order := CreateOrder{
    UserID:    123,
    ProductID: 456,
    Quantity:  2,
}

resp, err := corehttp.PostTyped[CreateOrder, Order](
    client,
    ctx,
    "https://api.example.com/orders",
    order, // Se serializa automáticamente a JSON
    corehttp.WithContentType("application/json"),
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Order created: %+v\n", resp.Body)
```

### Ejemplo 4: Criterios de Éxito Personalizados

```go
// Esperar específicamente código 201
resp, err := corehttp.PostTyped[CreateUser, User](
    client,
    ctx,
    "https://api.example.com/users",
    user,
    corehttp.WithExpectedStatusCode[User](201),
)
if err != nil {
    log.Fatal(err)
}

if resp.IsSuccess() {
    fmt.Println("User created successfully!")
}

// O usar un rango
resp.SetSuccessStatusCodeRange(&corehttp.StatusCodeRange{Min: 200, Max: 299})
```

### Ejemplo 5: Request con Query Parameters

```go
resp, err := corehttp.GetTyped[[]Product](
    client,
    ctx,
    "https://api.example.com/products",
    corehttp.WithQueryParam("category", "electronics"),
    corehttp.WithQueryParam("limit", "10"),
    corehttp.WithQueryParam("offset", "0"),
)
```

### Ejemplo 6: Request con Timeout

```go
resp, err := corehttp.GetTyped[Data](
    client,
    ctx,
    "https://api.example.com/data",
    corehttp.WithTimeout(5), // 5 segundos
)
```

### Ejemplo 7: Respuesta Binaria

```go
// Para respuestas binarias, usa []byte como tipo
resp, err := corehttp.GetTyped[[]byte](
    client,
    ctx,
    "https://api.example.com/image.jpg",
)
if err != nil {
    log.Fatal(err)
}

// El body ya está en resp.Body como []byte
ioutil.WriteFile("image.jpg", resp.Body, 0644)
```

## 🎯 Casos de Uso

### Integración con APIs REST

El módulo está diseñado para facilitar la integración con APIs REST externas:

```go
type APIClient struct {
    httpClient corehttp.Client
    baseURL    string
}

func (c *APIClient) GetUser(id int) (*User, error) {
    ctx := customctx.New(context.Background())
    resp, err := corehttp.GetTyped[User](
        c.httpClient,
        ctx,
        fmt.Sprintf("%s/users/%d", c.baseURL, id),
    )
    if err != nil {
        return nil, err
    }
    return &resp.Body, nil
}
```

### Testing con Mocks

La interfaz `Client` permite crear mocks fácilmente para testing:

```go
type MockClient struct{}

func (m *MockClient) Get(ctx *customctx.CustomContext, url string, opts ...corehttp.RequestOption) (*corehttp.Response, error) {
    // Implementación mock
}
```

### Medición de Performance

El `TypedResponse` incluye información de timing:

```go
resp, err := corehttp.GetTyped[Data](client, ctx, url)
if err != nil {
    return err
}

fmt.Printf("Request took: %v\n", resp.Duration)
fmt.Printf("Started at: %v\n", resp.RequestTime)
fmt.Printf("Completed at: %v\n", resp.ResponseTime)
```

## 🔗 Referencias

- [Infrastructure HTTP Implementation](../infrastructure/http/README.md) - Implementación concreta usando `net/http`
- [Ejemplos](../../../examples/infrastructure/http/) - Ejemplos de uso completos

---


## core/kerrors

**Source:** `docs/core/kerrors/README.md`

---

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

---


## core/logger

**Source:** `docs/core/logger/README.md`

---

# Logger

`Logger` proporciona una abstracción de logging basada en interfaces para el base-kit. Define un contrato claro para logging que permite diferentes implementaciones, incluyendo una implementación simple (`SimpleLogger`) para desarrollo y pruebas.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Niveles de Log](#niveles-de-log)
- [Ejemplos](#ejemplos)
- [Casos de Uso](#casos-de-uso)
- [Implementación Personalizada](#implementación-personalizada)

## ✨ Características

- **Interfaz simple**: Contrato claro y fácil de implementar
- **Múltiples niveles**: Debug, Info, Warn, Error, Fatal, Panic
- **Niveles configurables**: Control de qué mensajes se registran
- **Implementación simple**: `SimpleLogger` incluida para desarrollo
- **Extensible**: Fácil de implementar con backends personalizados

## 📦 Instalación

```bash
go get github.com/foundathyon/base/core/logger
```

## 🚀 Uso Básico

```go
import "github.com/foundathyon/base/core/logger"

// Crear un logger simple
logger := logger.NewSimpleLogger("info")

// Usar los diferentes niveles
logger.Debug("Mensaje de debug")
logger.Info("Mensaje informativo")
logger.Warn("Advertencia")
logger.Error("Error")
logger.Fatal("Error fatal")
logger.Panic("Error crítico")
```

## 📚 API

### Interfaces

#### `ILogger`

```go
type ILogger interface {
    Debug(message string, args ...any)
    Info(message string, args ...any)
    Warn(message string, args ...any)
    Error(message string, args ...any)
    Fatal(message string, args ...any)
    Panic(message string, args ...any)
}
```

`ILogger` define el contrato para todas las implementaciones de logger en el base-kit. Cada método acepta un mensaje y argumentos opcionales que se pasan a `fmt.Printf`.

**Métodos:**
- `Debug()`: Mensajes de depuración (solo en desarrollo)
- `Info()`: Mensajes informativos
- `Warn()`: Advertencias
- `Error()`: Errores que no detienen la ejecución
- `Fatal()`: Errores fatales que deberían detener la aplicación
- `Panic()`: Errores críticos que deberían causar panic

### Implementaciones

#### `SimpleLogger`

```go
type SimpleLogger struct {
    Level string
}
```

`SimpleLogger` es una implementación básica de `ILogger` que escribe a la salida estándar usando `fmt.Printf`. Ideal para desarrollo y pruebas.

**Campos:**
- `Level`: Nivel mínimo de logging (string): `"debug"`, `"info"`, `"warn"`, `"error"`, `"fatal"`, `"panic"`

#### `NewSimpleLogger(level string) *SimpleLogger`

Crea un nuevo `SimpleLogger` con el nivel especificado.

**Parámetros:**
- `level`: Nivel mínimo de logging como string (ej: `"info"`, `"debug"`, `"error"`)

**Retorna:**
- `*SimpleLogger`: Nueva instancia de SimpleLogger

**Ejemplo:**
```go
logger := logger.NewSimpleLogger("info")
logger := logger.NewSimpleLogger("debug")
logger := logger.NewSimpleLogger("error")
```

### Niveles de Logging

El paquete `logger` define los siguientes niveles de logging, ordenados de menor a mayor severidad:

- `LevelAll`: Registra todos los mensajes
- `LevelDebug`: Mensajes de depuración
- `LevelInfo`: Mensajes informativos
- `LevelWarn`: Advertencias
- `LevelError`: Errores
- `LevelFatal`: Errores fatales
- `LevelPanic`: Errores críticos

#### `Level`

```go
type Level int
```

`Level` representa un nivel de logging como constante de tipo entero.

**Constantes:**
```go
LevelAll   Level = iota // todos los mensajes
LevelDebug              // debug, info, warn, error, fatal, panic
LevelInfo               // info, warn, error, fatal, panic
LevelWarn               // warn, error, fatal, panic
LevelError              // error, fatal, panic
LevelFatal              // fatal, panic
LevelPanic              // panic
```

#### `(Level) String() string`

Convierte un `Level` a su representación string.

**Retorna:**
- `string`: Representación string del nivel (ej: `"debug"`, `"info"`, `"error"`)

**Ejemplo:**
```go
level := logger.LevelInfo
fmt.Println(level.String()) // "info"
```

#### `ParseLevel(s string) (lvl Level, ok bool)`

Convierte un string de nivel (`"debug"`, `"info"`, ...) en un `Level`. Retorna `ok=false` si el input es desconocido.

**Parámetros:**
- `s`: String del nivel a parsear

**Retorna:**
- `lvl`: El `Level` correspondiente
- `ok`: `true` si el string es válido, `false` en caso contrario

**Valores aceptados:**
- `"all"` → `LevelAll`
- `"debug"` → `LevelDebug`
- `"info"` → `LevelInfo`
- `"warn"` o `"warning"` → `LevelWarn`
- `"error"` → `LevelError`
- `"fatal"` → `LevelFatal`
- `"panic"` → `LevelPanic`

**Ejemplo:**
```go
level, ok := logger.ParseLevel("info")
if ok {
    fmt.Printf("Nivel: %s\n", level.String()) // "info"
}
```

#### `IsLoggable(minLevel, msgLevel Level) bool`

Retorna `true` si un mensaje en `msgLevel` debe ser registrado cuando el nivel mínimo configurado es `minLevel`.

**Parámetros:**
- `minLevel`: Nivel mínimo configurado
- `msgLevel`: Nivel del mensaje a verificar

**Retorna:**
- `bool`: `true` si el mensaje debe ser registrado, `false` en caso contrario

**Comportamiento:**
- `LevelAll` registra todo
- Niveles más altos (mayor severidad) tienen valores numéricos mayores
- Un mensaje se registra si `msgLevel >= minLevel`

**Ejemplo:**
```go
minLevel := logger.LevelInfo
msgLevel := logger.LevelDebug

if logger.IsLoggable(minLevel, msgLevel) {
    // No se registra porque Debug < Info
}

msgLevel = logger.LevelError
if logger.IsLoggable(minLevel, msgLevel) {
    // Se registra porque Error >= Info
}
```

#### `IsLoggableLevel(minLevelStr, msgLevelStr string) bool`

Wrapper de conveniencia para inputs string (env/config). Niveles desconocidos retornan `false` (fail-closed).

**Parámetros:**
- `minLevelStr`: Nivel mínimo como string
- `msgLevelStr`: Nivel del mensaje como string

**Retorna:**
- `bool`: `true` si el mensaje debe ser registrado, `false` en caso contrario

**Ejemplo:**
```go
shouldLog := logger.IsLoggableLevel("info", "debug") // false
shouldLog = logger.IsLoggableLevel("info", "error")  // true
```

## 💡 Ejemplos

### Ejemplo 1: Uso Básico

```go
package main

import (
    "github.com/foundathyon/base/core/logger"
)

func main() {
    logger := logger.NewSimpleLogger("info")
    
    logger.Debug("Este mensaje no se mostrará (nivel debug < info)")
    logger.Info("Este mensaje se mostrará")
    logger.Warn("Este mensaje se mostrará")
    logger.Error("Este mensaje se mostrará")
}
```

### Ejemplo 2: Diferentes Niveles

```go
// Logger con nivel debug (muestra todo)
debugLogger := logger.NewSimpleLogger("debug")
debugLogger.Debug("Mensaje de debug")
debugLogger.Info("Mensaje informativo")

// Logger con nivel error (solo errores y superiores)
errorLogger := logger.NewSimpleLogger("error")
errorLogger.Debug("No se muestra")
errorLogger.Info("No se muestra")
errorLogger.Warn("No se muestra")
errorLogger.Error("Este sí se muestra")
```

### Ejemplo 3: Uso con Formato

```go
logger := logger.NewSimpleLogger("info")

userID := 42
logger.Info("Usuario autenticado: ID=%d\n", userID)
logger.Warn("Intentos de login fallidos: count=%d\n", 3)
logger.Error("Error al conectar con base de datos: %v\n", err)
```

### Ejemplo 4: Integración con CustomContext

```go
import (
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/core/logger"
)

func main() {
    ctx := customctx.New(context.Background())
    
    logger := logger.NewSimpleLogger("info")
    ctxWithLogger := ctx.WithLogger(logger)
    
    // Usar el logger del contexto
    ctxLogger := ctxWithLogger.Logger()
    ctxLogger.Info("Mensaje usando logger del contexto\n")
}
```

### Ejemplo 5: Niveles Programáticos

```go
// Parsear nivel desde configuración
levelStr := os.Getenv("LOG_LEVEL") // ej: "info"
level, ok := logger.ParseLevel(levelStr)
if !ok {
    level = logger.LevelInfo // default
}

logger := logger.NewSimpleLogger(level.String())

// Verificar si un nivel debe ser registrado
if logger.IsLoggable(logger.LevelInfo, logger.LevelDebug) {
    logger.Debug("Este mensaje se mostraría")
}
```

## 🎯 Casos de Uso

### Configuración por Entorno

```go
func NewLoggerFromEnv() logger.ILogger {
    levelStr := os.Getenv("LOG_LEVEL")
    if levelStr == "" {
        levelStr = "info" // default
    }
    
    return logger.NewSimpleLogger(levelStr)
}

// Uso
logger := NewLoggerFromEnv()
logger.Info("Aplicación iniciada\n")
```

### Logging Estructurado

```go
func logRequest(logger logger.ILogger, requestID string, method string, path string) {
    logger.Info("[%s] %s %s\n", requestID, method, path)
}

func logError(logger logger.ILogger, err error, context map[string]any) {
    logger.Error("Error: %v, Context: %+v\n", err, context)
}
```

### Integración con Servicios

```go
type Service struct {
    logger logger.ILogger
}

func NewService(logger logger.ILogger) *Service {
    return &Service{logger: logger}
}

func (s *Service) ProcessRequest(ctx context.Context, data Data) error {
    s.logger.Info("Procesando request\n")
    
    // ... lógica del servicio
    
    if err != nil {
        s.logger.Error("Error procesando request: %v\n", err)
        return err
    }
    
    s.logger.Info("Request procesado exitosamente\n")
    return nil
}
```

## 🔧 Implementación Personalizada

Puedes crear tu propia implementación de `ILogger` para integrar con sistemas de logging como zap, logrus, o cualquier otro backend:

```go
type CustomLogger struct {
    zapLogger *zap.Logger
}

func NewCustomLogger(zapLogger *zap.Logger) *CustomLogger {
    return &CustomLogger{zapLogger: zapLogger}
}

func (l *CustomLogger) Debug(message string, args ...any) {
    l.zapLogger.Debug(fmt.Sprintf(message, args...))
}

func (l *CustomLogger) Info(message string, args ...any) {
    l.zapLogger.Info(fmt.Sprintf(message, args...))
}

func (l *CustomLogger) Warn(message string, args ...any) {
    l.zapLogger.Warn(fmt.Sprintf(message, args...))
}

func (l *CustomLogger) Error(message string, args ...any) {
    l.zapLogger.Error(fmt.Sprintf(message, args...))
}

func (l *CustomLogger) Fatal(message string, args ...any) {
    l.zapLogger.Fatal(fmt.Sprintf(message, args...))
}

func (l *CustomLogger) Panic(message string, args ...any) {
    l.zapLogger.Panic(fmt.Sprintf(message, args...))
}
```

Ahora puedes usar tu logger personalizado con `CustomContext`:

```go
zapLogger, _ := zap.NewProduction()
customLogger := NewCustomLogger(zapLogger)

ctx := customctx.New(context.Background())
ctxWithLogger := ctx.WithLogger(customLogger)
```

## 📝 Mejores Prácticas

### Selección de Niveles

- **Debug**: Información detallada para depuración, solo en desarrollo
- **Info**: Eventos normales del flujo de la aplicación
- **Warn**: Situaciones inusuales que no son errores pero merecen atención
- **Error**: Errores que no detienen la aplicación pero deben ser investigados
- **Fatal**: Errores críticos que deberían detener la aplicación
- **Panic**: Errores críticos que deberían causar panic

### Mensajes de Log

- Usa mensajes descriptivos y específicos
- Incluye contexto relevante (IDs, parámetros, etc.)
- Evita información sensible (contraseñas, tokens, etc.)
- Usa formato consistente

```go
// ❌ Mal
logger.Info("Error\n")

// ✅ Bien
logger.Info("Usuario autenticado: user_id=%d, ip=%s\n", userID, ipAddress)
```

### Configuración

- Configura el nivel desde variables de entorno
- Usa niveles más verbosos en desarrollo (debug)
- Usa niveles más restrictivos en producción (info o warn)

```go
func GetLogLevel() string {
    if env := os.Getenv("LOG_LEVEL"); env != "" {
        return env
    }
    
    // Default basado en entorno
    if os.Getenv("ENV") == "production" {
        return "info"
    }
    
    return "debug"
}
```

## 🔗 Ver También

- [CustomContext](../customctx/README.md) - Contexto que puede contener un logger
- [KErrors](../kerrors/README.md) - Errores estructurados para logging
- [Result](../result/README.md) - Tipo Result

## 📚 Referencias

- [Ejemplos de uso](../../../examples/core/customctx/customctx_example.go) - Ve cómo se usa logger con CustomContext
- [Tests](../../../core/logger/) - Implementación y contratos

---


## core/result

**Source:** `docs/core/result/README.md`

---

# Result

`Result` proporciona un tipo genérico `Result[T]` que representa un valor exitoso de tipo `T` o un error. Similar al tipo `Result` de Rust o el tipo `Either` de programación funcional, permite manejo explícito de errores sin depender del patrón tradicional de retorno de errores de Go.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Ejemplos](#ejemplos)
- [Casos de Uso](#casos-de-uso)
- [Comparación con Errores Estándar](#comparación-con-errores-estándar)

## ✨ Características

- **Tipo genérico**: Funciona con cualquier tipo `T`
- **Manejo funcional**: Enfoque declarativo para manejo de errores
- **Integración con KError**: Diseñado para trabajar con errores estructurados
- **Explícito**: Fuerza el manejo explícito de éxito y error
- **Type-safe**: El compilador ayuda a prevenir errores

## 📦 Instalación

```bash
go get github.com/foundathyon/base/core/result
```

## 🚀 Uso Básico

```go
import (
    "github.com/foundathyon/base/core/result"
    "github.com/foundathyon/base/core/kerrors"
)

// Result exitoso
r := result.Ok(42)
if r.IsOk() {
    value := r.Value() // 42
}

// Result con error
err := kerrors.NewKError("Error", 500, nil)
r := result.Err[int](err)
if !r.IsOk() {
    error := r.Error() // KError
}
```

## 📚 API

### Tipos

#### `Result[T]`

```go
type Result[T any] struct {
    value  T
    _error *kerrors.KError
}
```

`Result[T]` es un tipo genérico que representa un valor exitoso de tipo `T` o un error. Proporciona un enfoque funcional al manejo de errores, permitiendo encadenar operaciones y manejar errores de forma más declarativa.

**Campos internos:**
- `value`: El valor exitoso (si `IsOk()` es `true`)
- `_error`: El error estructurado (si `IsOk()` es `false`)

### Funciones de Construcción

#### `Ok[T any](value T) Result[T]`

Crea un nuevo `Result` exitoso con el valor dado. El `Result` no tendrá error y `IsOk()` retornará `true`.

**Parámetros:**
- `value`: El valor exitoso de tipo `T`

**Retorna:**
- `Result[T]`: Nueva instancia de Result exitosa

**Ejemplo:**
```go
r := result.Ok(42)
r := result.Ok("success")
r := result.Ok(user)
```

#### `Err[T any](err *kerrors.KError) Result[T]`

Crea un nuevo `Result` con error con el error dado. El `Result` tendrá un valor cero para `T` y `IsOk()` retornará `false`.

**Parámetros:**
- `err`: El error estructurado (`*kerrors.KError`)

**Retorna:**
- `Result[T]`: Nueva instancia de Result con error

**Ejemplo:**
```go
err := kerrors.NewKError("Error", 500, nil)
r := result.Err[int](err)
r := result.Err[string](err)
```

#### `NewOkResult[T any](value T) Result[T]`

Crea un nuevo `Result` exitoso con el valor dado. Equivalente a `Ok()` pero con nombre más explícito.

**Parámetros:**
- `value`: El valor exitoso de tipo `T`

**Retorna:**
- `Result[T]`: Nueva instancia de Result exitosa

#### `NewErrResult[T any](err *kerrors.KError) Result[T]`

Crea un nuevo `Result` con error con el error dado. Equivalente a `Err()` pero con nombre más explícito.

**Parámetros:**
- `err`: El error estructurado (`*kerrors.KError`)

**Retorna:**
- `Result[T]`: Nueva instancia de Result con error

#### `NewResult[T any](value T, err *kerrors.KError) Result[T]`

Crea un nuevo `Result` con ambos valor y error. Si `err` es `nil`, el `Result` se considera exitoso (Ok). Si `err` no es `nil`, el `Result` se considera fallido (Err).

**Parámetros:**
- `value`: El valor de tipo `T`
- `err`: El error estructurado (puede ser `nil`)

**Retorna:**
- `Result[T]`: Nueva instancia de Result

**Ejemplo:**
```go
// Result exitoso
r := result.NewResult(42, nil)

// Result con error
err := kerrors.NewKError("Error", 500, nil)
r := result.NewResult(0, err)
```

#### `Empty[T any]() Result[T]`

Crea un `Result` vacío (sin valor ni error). Útil para inicialización o casos especiales.

**Retorna:**
- `Result[T]`: Result vacío

### Métodos

#### `Value() T`

Retorna el valor contenido en el `Result`.

**Nota:** Este método retorna el valor independientemente de si el `Result` es Ok o Err. Usa `IsOk()` para verificar si el `Result` es exitoso antes de acceder al valor.

**Retorna:**
- `T`: El valor contenido (puede ser valor cero si es Err)

**Ejemplo:**
```go
r := result.Ok(42)
value := r.Value() // 42

err := kerrors.NewKError("Error", 500, nil)
r = result.Err[int](err)
value = r.Value() // 0 (valor cero de int)
```

#### `Error() error`

Retorna el error contenido en el `Result`. Si el `Result` es exitoso (Ok), esto retornará `nil`. Si el `Result` es fallido (Err), esto retornará el error que fue configurado.

**Retorna:**
- `error`: El error estructurado o `nil` si es exitoso

**Ejemplo:**
```go
r := result.Ok(42)
if r.Error() == nil {
    fmt.Println("No hay error")
}

err := kerrors.NewKError("Error", 500, nil)
r = result.Err[int](err)
if r.Error() != nil {
    fmt.Printf("Error: %v\n", r.Error())
}
```

#### `IsOk() bool`

Retorna `true` si el `Result` representa un valor exitoso (sin error). Retorna `false` si el `Result` representa un error.

**Retorna:**
- `bool`: `true` si es exitoso, `false` si hay error

**Ejemplo:**
```go
r := result.Ok(42)
if r.IsOk() {
    value := r.Value() // Acceder al valor
}

err := kerrors.NewKError("Error", 500, nil)
r = result.Err[int](err)
if !r.IsOk() {
    error := r.Error() // Manejar el error
}
```

#### `IsEmpty() bool`

Retorna `true` si el `Result` está vacío (sin valor ni error). Retorna `false` si el `Result` tiene un valor o un error.

**Retorna:**
- `bool`: `true` si está vacío, `false` en caso contrario

**Ejemplo:**
```go
r := result.Empty[int]()
if r.IsEmpty() {
    fmt.Println("Result está vacío")
}
```

#### `ToKError() *kerrors.KError`

Retorna el error estructurado si el `Result` tiene un error, o `nil` si es exitoso. Útil para obtener directamente el `*kerrors.KError` sin necesidad de type assertion.

**Retorna:**
- `*kerrors.KError`: El error estructurado o `nil` si es exitoso

**Ejemplo:**
```go
err := kerrors.NewKError("Error", 500, nil)
r := result.Err[int](err)

kerr := r.ToKError()
if kerr != nil {
    fmt.Printf("Error Code: %d\n", kerr.Code)
    fmt.Printf("Error Message: %s\n", kerr.Message)
}
```

## 💡 Ejemplos

### Ejemplo 1: Result Básico

```go
package main

import (
    "fmt"
    "github.com/foundathyon/base/core/result"
)

func main() {
    // Result exitoso
    r := result.Ok(42)
    fmt.Printf("IsOk: %v\n", r.IsOk())
    fmt.Printf("Value: %d\n", r.Value())
    fmt.Printf("Error: %v\n", r.Error())
}
```

### Ejemplo 2: Result con Error

```go
import (
    "github.com/foundathyon/base/core/kerrors"
    "github.com/foundathyon/base/core/result"
)

err := kerrors.NewKError("No se pudo procesar la solicitud", 500, map[string]any{
    "request_id": "req-123",
})
r := result.Err[int](err)

fmt.Printf("IsOk: %v\n", r.IsOk()) // false
fmt.Printf("Value: %d\n", r.Value()) // 0 (valor cero)
fmt.Printf("Error: %v\n", r.Error())

if kerr, ok := r.Error().(*kerrors.KError); ok {
    fmt.Printf("Error Code: %d\n", kerr.Code)
    fmt.Printf("Error Message: %s\n", kerr.Message)
}
```

### Ejemplo 3: Verificación de Result

```go
results := []result.Result[int]{
    result.Ok(100),
    result.Err[int](kerrors.NewKError("Error de validación", 400, nil)),
    result.Ok(200),
}

for i, r := range results {
    fmt.Printf("Result %d:\n", i+1)
    if r.IsOk() {
        fmt.Printf("  ✓ Éxito: %d\n", r.Value())
    } else {
        kerr := r.ToKError()
        fmt.Printf("  ✗ Error: %s (Código: %d)\n", kerr.Message, kerr.Code)
    }
}
```

### Ejemplo 4: Result con Diferentes Tipos

```go
// Result con string
strResult := result.Ok("Hola, mundo!")
fmt.Printf("String: %s\n", strResult.Value())

// Result con slice
sliceResult := result.Ok([]int{1, 2, 3, 4, 5})
fmt.Printf("Slice: %v\n", sliceResult.Value())

// Result con struct
type Usuario struct {
    ID     int
    Nombre string
    Email  string
}

userResult := result.Ok(Usuario{
    ID:     1,
    Nombre: "Juan Pérez",
    Email:  "juan@example.com",
})
fmt.Printf("Usuario: %+v\n", userResult.Value())

// Result con error
err := kerrors.NewKError("Usuario no encontrado", 404, nil)
userErrResult := result.Err[Usuario](err)
if !userErrResult.IsOk() {
    fmt.Printf("Error: %s\n", userErrResult.ToKError().Message)
}
```

### Ejemplo 5: Función que Retorna Result

```go
func divide(a, b int) result.Result[float64] {
    if b == 0 {
        return result.Err[float64](kerrors.NewKError(
            "División por cero",
            400,
            map[string]any{
                "operation": "divide",
                "dividend":  a,
                "divisor":   b,
            },
        ))
    }
    return result.Ok(float64(a) / float64(b))
}

// Uso
r := divide(10, 2)
if r.IsOk() {
    fmt.Printf("Resultado: %.2f\n", r.Value())
} else {
    kerr := r.ToKError()
    fmt.Printf("Error: %s (Código: %d)\n", kerr.Message, kerr.Code)
}

r = divide(10, 0)
if !r.IsOk() {
    kerr := r.ToKError()
    fmt.Printf("Error: %s\n", kerr.Message)
    fmt.Printf("Metadata: %v\n", kerr.Metadata)
}
```

### Ejemplo 6: Encadenamiento de Operaciones

```go
func getUser(userID int) result.Result[User] {
    // Simular validación
    if userID <= 0 {
        return result.Err[User](kerrors.NewKError(
            "ID de usuario inválido",
            400,
            map[string]any{"user_id": userID},
        ))
    }
    
    // Simular obtención de usuario
    if userID == 999 {
        return result.Err[User](kerrors.NewKError(
            "Usuario no encontrado",
            404,
            map[string]any{"user_id": userID},
        ))
    }
    
    return result.Ok(User{
        ID:    userID,
        Name:  "Juan Pérez",
        Email: "juan@example.com",
    })
}

func getProfile(user User) result.Result[Profile] {
    // Simular obtención de perfil
    return result.Ok(Profile{
        UserID:   user.ID,
        Bio:      "Desarrollador",
        Location: "México",
    })
}

// Uso encadenado
userResult := getUser(42)
if !userResult.IsOk() {
    fmt.Printf("Error al obtener usuario: %s\n", userResult.ToKError().Message)
    return
}

profileResult := getProfile(userResult.Value())
if !profileResult.IsOk() {
    fmt.Printf("Error al obtener perfil: %s\n", profileResult.ToKError().Message)
    return
}

profile := profileResult.Value()
fmt.Printf("Perfil: %+v\n", profile)
```

## 🎯 Casos de Uso

### Validación con Result

`Result` es ideal para funciones de validación donde quieres retornar un valor o un error:

```go
func validateEmail(email string) result.Result[string] {
    if email == "" {
        return result.Err[string](kerrors.NewKError(
            "Email es requerido",
            400,
            map[string]any{"field": "email"},
        ))
    }
    
    if !strings.Contains(email, "@") {
        return result.Err[string](kerrors.NewKError(
            "Email inválido",
            400,
            map[string]any{
                "field": "email",
                "value": email,
            },
        ))
    }
    
    return result.Ok(email)
}
```

### Operaciones de Base de Datos

```go
func findUserByID(ctx context.Context, userID int) result.Result[User] {
    user, err := db.QueryUser(ctx, userID)
    if err != nil {
        return result.Err[User](kerrors.NewKErrorWithCause(
            "Error al obtener usuario",
            500,
            map[string]any{"user_id": userID},
            err,
        ))
    }
    
    if user == nil {
        return result.Err[User](kerrors.NewKError(
            "Usuario no encontrado",
            404,
            map[string]any{"user_id": userID},
        ))
    }
    
    return result.Ok(*user)
}
```

### Procesamiento de Datos

```go
func processData(input Data) result.Result[Output] {
    // Validar input
    if err := validateInput(input); err != nil {
        return result.Err[Output](err)
    }
    
    // Procesar
    output, err := transform(input)
    if err != nil {
        return result.Err[Output](kerrors.NewKErrorWithCause(
            "Error al procesar datos",
            500,
            map[string]any{"input": input},
            err,
        ))
    }
    
    return result.Ok(output)
}
```

## 🔄 Comparación con Errores Estándar

### Patrón Estándar de Go

```go
// Go tradicional
func getUser(id int) (User, error) {
    if id <= 0 {
        return User{}, errors.New("ID inválido")
    }
    // ...
    return user, nil
}

// Uso
user, err := getUser(42)
if err != nil {
    return err
}
// usar user
```

### Con Result

```go
// Con Result
func getUser(id int) result.Result[User] {
    if id <= 0 {
        return result.Err[User](kerrors.NewKError("ID inválido", 400, nil))
    }
    // ...
    return result.Ok(user)
}

// Uso
userResult := getUser(42)
if !userResult.IsOk() {
    return userResult.ToKError()
}
user := userResult.Value()
// usar user
```

### Ventajas de Result

1. **Explícito**: Fuerza el manejo explícito de éxito y error
2. **Type-safe**: El compilador ayuda a prevenir accesos incorrectos
3. **Funcional**: Permite encadenar operaciones de forma más natural
4. **Integración con KError**: Diseñado para trabajar con errores estructurados
5. **Valor cero explícito**: Puedes verificar si hay valor antes de acceder

### Cuándo Usar Result

**Usar Result cuando:**
- ✅ Quieres un enfoque más funcional
- ✅ Necesitas encadenar operaciones
- ✅ Trabajas con errores estructurados (KError)
- ✅ Quieres forzar manejo explícito de errores

**Usar errores estándar cuando:**
- ✅ Quieres seguir el patrón tradicional de Go
- ✅ Necesitas compatibilidad con bibliotecas estándar
- ✅ Prefieres el patrón establecido de Go

## 📝 Mejores Prácticas

### Verificar Antes de Acceder

Siempre verifica `IsOk()` antes de acceder al valor:

```go
// ❌ Mal
value := r.Value() // Puede ser valor cero si hay error

// ✅ Bien
if r.IsOk() {
    value := r.Value() // Seguro
} else {
    error := r.Error() // Manejar error
}
```

### Usar ToKError()

Cuando trabajas con `KError`, usa `ToKError()` para obtener el error directamente:

```go
// ❌ Menos claro
if !r.IsOk() {
    if kerr, ok := r.Error().(*kerrors.KError); ok {
        // usar kerr
    }
}

// ✅ Mejor
if !r.IsOk() {
    kerr := r.ToKError()
    if kerr != nil {
        // usar kerr
    }
}
```

### Mensajes de Error Claros

Cuando creas errores para `Result`, usa mensajes descriptivos:

```go
// ❌ Mal
return result.Err[int](kerrors.NewKError("Error", 500, nil))

// ✅ Bien
return result.Err[int](kerrors.NewKError(
    "Error al procesar solicitud",
    500,
    map[string]any{
        "operation": "process_request",
        "request_id": requestID,
    },
))
```

## 🔗 Ver También

- [KErrors](../kerrors/README.md) - Errores estructurados usados en Result
- [CustomContext](../customctx/README.md) - Contexto que acumula errores
- [Logger](../logger/README.md) - Sistema de logging

## 📚 Referencias

- [Ejemplos de uso](../../../examples/core/result/result_example.go)
- [Tests](../../../core/result/result_test.go)
- [Rust Result Type](https://doc.rust-lang.org/std/result/) - Inspiración para este diseño

---


## infrastructure/eventbus/inmemory

**Source:** `docs/infrastructure/eventbus/inmemory/README.md`

---

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

---


## infrastructure/eventbus/rabbit

**Source:** `docs/infrastructure/eventbus/rabbit/README.md`

---

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

---


## infrastructure/eventbus

**Source:** `docs/infrastructure/eventbus/README.md`

---

# Event Bus - Infrastructure

Este módulo contiene las implementaciones del event bus para diferentes backends.

## Implementaciones Disponibles

### In-Memory

Implementación en memoria para testing y desarrollo.

**Ubicación:** [`inmemory/`](./inmemory/README.md)

**Características:**
- Thread-safe
- Soporte para múltiples consumidores
- Delay mediante goroutines
- Sin persistencia

**Ideal para:**
- Testing unitario e integración
- Desarrollo local
- Prototipado rápido

### RabbitMQ

Implementación usando RabbitMQ para producción.

**Ubicación:** [`rabbit/`](./rabbit/README.md)

**Características:**
- Persistencia de mensajes
- Alta disponibilidad
- Delay usando DLX con TTL
- Broadcast con múltiples queues
- Soporte para microservicios

**Ideal para:**
- Producción
- Microservicios distribuidos
- Aplicaciones que requieren persistencia

## Elegir una Implementación

### Usa In-Memory si:
- ✅ Estás haciendo testing
- ✅ Desarrollando localmente sin infraestructura
- ✅ Tienes una aplicación simple y monolítica
- ✅ No necesitas persistencia

### Usa RabbitMQ si:
- ✅ Estás en producción
- ✅ Tienes múltiples servicios/microservicios
- ✅ Necesitas persistencia de eventos
- ✅ Necesitas alta disponibilidad
- ✅ Necesitas comunicación entre procesos/máquinas

## Ejemplos

### In-Memory

Ver `examples/infrastructure/eventbus/inmemory/main.go` para un ejemplo completo.

### RabbitMQ

Ver `examples/infrastructure/eventbus/rabbit/` para múltiples ejemplos:

- **Basic**: Publicación y consumo básico
- **Delayed**: Mensajes con delay
- **Signup**: Patrón broadcast con múltiples consumidores

## Referencias

- [Event Bus Core Documentation](../../core/eventbus/README.md)
- [In-Memory Implementation](./inmemory/README.md)
- [RabbitMQ Implementation](./rabbit/README.md)

---


## infrastructure/http

**Source:** `docs/infrastructure/http/README.md`

---

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

---


## observability/logger/loguru

**Source:** `docs/observability/logger/loguru/README.md`

---

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

---


## persistence/contracts

**Source:** `docs/persistence/contracts/README.md`

---

# Contracts

`Contracts` define las interfaces fundamentales para el módulo de persistencia. Proporciona los contratos del dominio (`IEntity`, `IModel`) y funciones de conversión entre entidades y modelos, manteniendo la separación entre el dominio de negocio y la infraestructura.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Ejemplos](#ejemplos)
- [Mejores Prácticas](#mejores-prácticas)

## ✨ Características

- **Separación de dominio e infraestructura**: Interfaces claras para entidades y modelos
- **Conversiones type-safe**: Funciones genéricas para convertir Entity ↔ Model
- **Integración con Result**: Conversiones retornan `Result[T]` para manejo funcional
- **Manejo de errores estructurado**: Usa `KError` para errores de conversión

## 📦 Instalación

```bash
go get github.com/foundathyon/base/persistence/contracts
```

## 🚀 Uso Básico

```go
import "github.com/foundathyon/base/persistence/contracts"

// Definir entidad
type User struct {
    ID    string
    Name  string
    Email string
}

func (u User) GetID() string { return u.ID }

// Definir modelo
type UserModel struct {
    ID        string
    Name      string
    Email     string
    CreatedAt time.Time
}

func (m UserModel) GetID() string { return m.ID }

// Convertir Entity a Model
entity := User{ID: "123", Name: "Juan", Email: "juan@example.com"}
modelResult := contracts.EntityToModel[User, UserModel](entity)

// Convertir Model a Entity
model := UserModel{ID: "123", Name: "Juan", Email: "juan@example.com"}
entityResult := contracts.ModelToEntity[User, UserModel](model)
```

## 📚 API

### Interfaces

#### `IEntity`

```go
type IEntity interface {
    GetID() string
}
```

`IEntity` representa una entidad de dominio. Las entidades son parte del dominio de negocio y no deben contener campos técnicos de infraestructura.

**Métodos:**
- `GetID() string`: Retorna el identificador único de la entidad

#### `IModel`

```go
type IModel interface {
    GetID() string
}
```

`IModel` representa un modelo de persistencia. Los modelos son parte de la capa de infraestructura y pueden incluir campos técnicos como `CreatedAt`, `UpdatedAt`, `IsRemoved`, etc.

**Métodos:**
- `GetID() string`: Retorna el identificador único del modelo

### Funciones de Conversión

#### `ToJSON[E IEntity](entity E) []byte`

Convierte una entidad a JSON indentado. Útil para serialización y debugging.

**Parámetros:**
- `entity`: La entidad a convertir

**Retorna:**
- `[]byte`: Representación JSON de la entidad

**Nota:** Si hay error en la serialización, retorna `nil` y escribe el error en stdout.

**Ejemplo:**
```go
user := User{ID: "123", Name: "Juan", Email: "juan@example.com"}
jsonData := contracts.ToJSON(user)
fmt.Println(string(jsonData))
```

#### `FromJSON[E IEntity](m map[string]interface{}) (E, error)`

Convierte un mapa a una entidad usando JSON como intermediario.

**Parámetros:**
- `m`: Mapa con los datos de la entidad

**Retorna:**
- `E`: La entidad convertida
- `error`: Error si la conversión falla

**Ejemplo:**
```go
data := map[string]interface{}{
    "id":    "123",
    "name":  "Juan",
    "email": "juan@example.com",
}

user, err := contracts.FromJSON[User](data)
if err != nil {
    // Manejar error
}
```

#### `EntityToModel[E IEntity, M IModel](entity IEntity) result.Result[M]`

Convierte una entidad de dominio a un modelo de persistencia. Retorna un `Result[M]` para manejo funcional de errores.

**Parámetros:**
- `entity`: La entidad a convertir

**Retorna:**
- `result.Result[M]`: Resultado con el modelo o error estructurado

**Errores posibles:**
- `ErrorConvertEntityToJSON`: Error al convertir entidad a JSON
- `ErrorConvertJSONToMap`: Error al convertir JSON a mapa
- `ErrorConvertMapToModel`: Error al convertir mapa a modelo

**Ejemplo:**
```go
user := User{ID: "123", Name: "Juan", Email: "juan@example.com"}
modelResult := contracts.EntityToModel[User, UserModel](user)

if modelResult.IsOk() {
    model := modelResult.Value()
    // Usar modelo
} else {
    kerr := modelResult.ToKError()
    fmt.Printf("Error: %s\n", kerr.Message)
}
```

#### `ModelToEntity[E IEntity, M IModel](model IModel) result.Result[E]`

Convierte un modelo de persistencia a una entidad de dominio. Retorna un `Result[E]` para manejo funcional de errores.

**Parámetros:**
- `model`: El modelo a convertir

**Retorna:**
- `result.Result[E]`: Resultado con la entidad o error estructurado

**Errores posibles:**
- `ErrorConvertModelToEntity`: Error al convertir modelo a entidad
- `ErrorConvertEntityToModel`: Error en conversión intermedia

**Ejemplo:**
```go
model := UserModel{
    ID:        "123",
    Name:      "Juan",
    Email:     "juan@example.com",
    CreatedAt: time.Now(),
}

entityResult := contracts.ModelToEntity[User, UserModel](model)
if entityResult.IsOk() {
    entity := entityResult.Value()
    // Usar entidad (sin campos técnicos)
} else {
    kerr := entityResult.ToKError()
    fmt.Printf("Error: %s\n", kerr.Message)
}
```

### Errores Predefinidos

El paquete define varios errores estáticos para conversiones:

```go
var (
    ErrorConvertEntityToJSON = kerrors.NewKError("Error al convertir la entidad a JSON", 500, nil)
    ErrorConvertJSONToMap    = kerrors.NewKError("Error al convertir los bytes JSON a un mapa", 500, nil)
    ErrorConvertMapToModel   = kerrors.NewKError("Error al convertir el mapa a modelo", 500, nil)
    ErrorConvertModelToEntity = kerrors.NewKError("Error al convertir el modelo a entidad", 500, nil)
    ErrorConvertEntityToModel = kerrors.NewKError("Error al convertir la entidad a modelo", 500, nil)
)
```

Estos errores pueden tener un `Cause` asociado cuando se usan con `WithCause()`.

## 💡 Ejemplos

### Ejemplo 1: Definir Entidad y Modelo

```go
package domain

import "github.com/foundathyon/base/persistence/contracts"

// Entidad de dominio (sin campos técnicos)
type User struct {
    ID    string
    Name  string
    Email string
    Age   int
}

func (u User) GetID() string {
    return u.ID
}

// Implementa IEntity
var _ contracts.IEntity = User{}
```

```go
package infrastructure

import (
    "time"
    "github.com/foundathyon/base/persistence/contracts"
)

// Modelo de infraestructura (con campos técnicos)
type UserModel struct {
    ID        string    `gorm:"primary_key" json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    IsRemoved bool      `json:"is_removed"`
}

func (m UserModel) GetID() string {
    return m.ID
}

// Implementa IModel
var _ contracts.IModel = UserModel{}
```

### Ejemplo 2: Conversión Entity → Model

```go
import (
    "github.com/foundathyon/base/persistence/contracts"
)

func main() {
    // Crear entidad de dominio
    user := domain.User{
        ID:    "user-123",
        Name:  "Juan Pérez",
        Email: "juan@example.com",
        Age:   30,
    }
    
    // Convertir a modelo
    modelResult := contracts.EntityToModel[domain.User, infrastructure.UserModel](user)
    
    if modelResult.IsOk() {
        model := modelResult.Value()
        fmt.Printf("Modelo: %+v\n", model)
        // El modelo tendrá CreatedAt, UpdatedAt, etc. con valores por defecto
    } else {
        kerr := modelResult.ToKError()
        fmt.Printf("Error de conversión: %s\n", kerr.Message)
        if kerr.Cause != nil {
            fmt.Printf("Causa: %v\n", kerr.Cause)
        }
    }
}
```

### Ejemplo 3: Conversión Model → Entity

```go
func main() {
    // Crear modelo de infraestructura
    model := infrastructure.UserModel{
        ID:        "user-123",
        Name:      "Juan Pérez",
        Email:     "juan@example.com",
        Age:       30,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        IsRemoved: false,
    }
    
    // Convertir a entidad
    entityResult := contracts.ModelToEntity[domain.User, infrastructure.UserModel](model)
    
    if entityResult.IsOk() {
        entity := entityResult.Value()
        fmt.Printf("Entidad: %+v\n", entity)
        // La entidad NO tendrá CreatedAt, UpdatedAt, IsRemoved
        // Solo campos de negocio: ID, Name, Email, Age
    } else {
        kerr := entityResult.ToKError()
        fmt.Printf("Error de conversión: %s\n", kerr.Message)
    }
}
```

### Ejemplo 4: Uso en Repositorio

```go
func (r *Repository) Save(ctx *customctx.CustomContext, user domain.User) result.Result[domain.User] {
    // Convertir Entity a Model para guardar
    modelResult := contracts.EntityToModel[domain.User, infrastructure.UserModel](user)
    if !modelResult.IsOk() {
        kerr := modelResult.ToKError()
        ctx.AddError(kerr)
        return result.Err[domain.User](kerr)
    }
    
    model := modelResult.Value()
    model.CreatedAt = time.Now() // Setear campos técnicos
    model.UpdatedAt = time.Now()
    
    // Guardar modelo (operación de infraestructura)
    savedModel := saveToDatabase(model)
    
    // Convertir Model a Entity para retornar
    entityResult := contracts.ModelToEntity[domain.User, infrastructure.UserModel](savedModel)
    if !entityResult.IsOk() {
        kerr := entityResult.ToKError()
        ctx.AddError(kerr)
        return result.Err[domain.User](kerr)
    }
    
    return result.Ok(entityResult.Value())
}
```

### Ejemplo 5: Serialización a JSON

```go
func main() {
    user := domain.User{
        ID:    "user-123",
        Name:  "Juan Pérez",
        Email: "juan@example.com",
        Age:   30,
    }
    
    // Convertir a JSON
    jsonData := contracts.ToJSON(user)
    if jsonData != nil {
        fmt.Println(string(jsonData))
    }
    
    // Salida:
    // {
    //   "ID": "user-123",
    //   "Name": "Juan Pérez",
    //   "Email": "juan@example.com",
    //   "Age": 30
    // }
}
```

## 📝 Mejores Prácticas

### Separación Clara de Responsabilidades

**✅ Correcto:**
- Entidades solo contienen lógica de negocio
- Modelos contienen campos técnicos de infraestructura
- Las conversiones se hacen en adaptadores (repositorios)

**❌ Incorrecto:**
- Mezclar campos técnicos en entidades
- Incluir lógica de negocio en modelos
- Convertir en el dominio en lugar de en adaptadores

### Manejo de Errores

Siempre verificar el resultado de las conversiones:

```go
// ✅ Correcto
modelResult := contracts.EntityToModel[User, UserModel](user)
if !modelResult.IsOk() {
    kerr := modelResult.ToKError()
    ctx.AddError(kerr)
    return result.Err[User](kerr)
}
model := modelResult.Value()

// ❌ Incorrecto (no verificar errores)
model := contracts.EntityToModel[User, UserModel](user).Value() // Puede panic
```

### Campos Técnicos en Modelos

Los modelos deben incluir campos técnicos comunes:

```go
type UserModel struct {
    // Campos de negocio (mapean de Entity)
    ID    string
    Name  string
    Email string
    
    // Campos técnicos de infraestructura
    CreatedAt time.Time
    UpdatedAt time.Time
    IsRemoved bool      // Soft delete
    Version   int       // Optimistic locking (opcional)
}
```

### Type Safety

Usa los tipos genéricos correctamente:

```go
// ✅ Correcto - tipos explícitos
modelResult := contracts.EntityToModel[domain.User, infrastructure.UserModel](entity)

// ❌ Incorrecto - tipos inferidos incorrectamente
modelResult := contracts.EntityToModel(entity) // No compila
```

## 🔗 Ver También

- [Persistence Overview](../README.md) - Visión general del módulo
- [InMemory](../inmemory/README.md) - Implementación en memoria
- [Postgres](../postgres/README.md) - Implementación PostgreSQL
- [Result](../../core/result/README.md) - Tipo Result usado en conversiones
- [KErrors](../../core/kerrors/README.md) - Errores estructurados

## 📚 Referencias

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [Código fuente](../../../persistence/contracts/)

---


## persistence/criteria

**Source:** `docs/persistence/criteria/README.md`

---

# Criteria

`Criteria` proporciona un sistema de filtrado y búsqueda para repositorios. Permite construir consultas usando filtros con diferentes operadores SQL, manteniendo la abstracción sobre cómo se implementan las búsquedas en cada repositorio.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Ejemplos](#ejemplos)
- [Operadores Disponibles](#operadores-disponibles)

## ✨ Características

- **Filtros genéricos**: Sistema de filtrado flexible y type-safe
- **Múltiples operadores**: Soporta operadores SQL estándar
- **Composición de filtros**: Permite combinar múltiples filtros
- **Independiente de implementación**: Los filtros se aplican según la implementación del repositorio

## 📦 Instalación

```bash
go get github.com/foundathyon/base/persistence/criteria
```

## 🚀 Uso Básico

```go
import "github.com/foundathyon/base/persistence/criteria"

// Crear un filtro
filter := criteria.Filter{
    Field:    "email",
    Operator: criteria.OperatorEqual,
    Value:    "juan@example.com",
}

// Crear múltiples filtros
filters := criteria.NewFilters([]criteria.Filter{
    {Field: "email", Operator: criteria.OperatorEqual, Value: "juan@example.com"},
    {Field: "age", Operator: criteria.OperatorGreaterThan, Value: 18},
})

// Crear criteria con filtros
crit := criteria.Criteria{
    Filters: *filters,
}
```

## 📚 API

### Tipos

#### `Criteria`

```go
type Criteria struct {
    Filters Filters
}
```

`Criteria` representa un criterio de búsqueda que contiene una colección de filtros.

**Campos:**
- `Filters`: Colección de filtros a aplicar

#### `Filter`

```go
type Filter struct {
    Field    FilterField
    Operator Operator
    Value    interface{}
}
```

`Filter` representa una condición de filtro individual.

**Campos:**
- `Field`: Nombre del campo sobre el que se aplica el filtro
- `Operator`: Operador a usar (ver operadores disponibles)
- `Value`: Valor a comparar (puede ser `int` o `string`)

#### `FilterField`

```go
type FilterField string
```

`FilterField` representa el nombre del campo sobre el que se aplica el filtro.

#### `Filters`

```go
type Filters struct {
    filters []Filter
}
```

`Filters` es una colección de filtros.

**Métodos:**
- `Get() []Filter`: Retorna todos los filtros

#### `NewFilters(filters []Filter) *Filters`

Crea una nueva colección de filtros.

**Parámetros:**
- `filters`: Slice de filtros

**Retorna:**
- `*Filters`: Nueva colección de filtros

### Operadores

#### `Operator`

```go
type Operator string
```

`Operator` representa un operador SQL válido.

**Operadores disponibles:**
- `OperatorEqual` (`"="`): Igual a
- `OperatorNotEqual` (`"<>"`): No igual a
- `OperatorGreaterThan` (`">"`): Mayor que
- `OperatorGreaterEqual` (`">="`): Mayor o igual que
- `OperatorLessThan` (`"<"`): Menor que
- `OperatorLessEqual` (`"<="`): Menor o igual que
- `OperatorLike` (`"LIKE"`): Coincide con patrón (para strings)
- `OperatorNotLike` (`"NOT LIKE"`): No coincide con patrón
- `OperatorIn` (`"IN"`): Está en lista
- `OperatorNotIn` (`"NOT IN"`): No está en lista

## 💡 Ejemplos

### Ejemplo 1: Filtro Simple

```go
import "github.com/foundathyon/base/persistence/criteria"

// Filtrar por email igual
filter := criteria.Filter{
    Field:    criteria.FilterField("email"),
    Operator: criteria.OperatorEqual,
    Value:    "juan@example.com",
}

filters := criteria.NewFilters([]criteria.Filter{filter})
crit := criteria.Criteria{
    Filters: *filters,
}
```

### Ejemplo 2: Múltiples Filtros

```go
// Filtrar usuarios con email específico y edad mayor a 18
filters := criteria.NewFilters([]criteria.Filter{
    {
        Field:    criteria.FilterField("email"),
        Operator: criteria.OperatorEqual,
        Value:    "juan@example.com",
    },
    {
        Field:    criteria.FilterField("age"),
        Operator: criteria.OperatorGreaterThan,
        Value:    18,
    },
})

crit := criteria.Criteria{
    Filters: *filters,
}
```

### Ejemplo 3: Operadores de Comparación

```go
// Mayor que
filter := criteria.Filter{
    Field:    criteria.FilterField("age"),
    Operator: criteria.OperatorGreaterThan,
    Value:    18,
}

// Menor o igual que
filter := criteria.Filter{
    Field:    criteria.FilterField("price"),
    Operator: criteria.OperatorLessEqual,
    Value:    100,
}

// Entre valores (usando dos filtros)
filters := criteria.NewFilters([]criteria.Filter{
    {
        Field:    criteria.FilterField("age"),
        Operator: criteria.OperatorGreaterEqual,
        Value:    18,
    },
    {
        Field:    criteria.FilterField("age"),
        Operator: criteria.OperatorLessEqual,
        Value:    65,
    },
})
```

### Ejemplo 4: Operadores LIKE

```go
// Búsqueda parcial de email
filter := criteria.Filter{
    Field:    criteria.FilterField("email"),
    Operator: criteria.OperatorLike,
    Value:    "%@example.com", // Emails que terminan en @example.com
}

// Búsqueda que no contiene
filter := criteria.Filter{
    Field:    criteria.FilterField("name"),
    Operator: criteria.OperatorNotLike,
    Value:    "%admin%", // Nombres que no contienen "admin"
}
```

### Ejemplo 5: Operadores IN y NOT IN

```go
// Valores en lista (requiere implementación especial en repositorio)
filter := criteria.Filter{
    Field:    criteria.FilterField("status"),
    Operator: criteria.OperatorIn,
    Value:    []string{"active", "pending"}, // Lista de valores
}

// Valores no en lista
filter := criteria.Filter{
    Field:    criteria.FilterField("role"),
    Operator: criteria.OperatorNotIn,
    Value:    []string{"admin", "superuser"},
}
```

### Ejemplo 6: Uso con Repositorio

```go
import (
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/persistence/criteria"
)

func searchUsers(ctx *customctx.CustomContext, repo Repository, email string, minAge int) {
    // Construir criteria
    filters := criteria.NewFilters([]criteria.Filter{
        {
            Field:    criteria.FilterField("email"),
            Operator: criteria.OperatorEqual,
            Value:    email,
        },
        {
            Field:    criteria.FilterField("age"),
            Operator: criteria.OperatorGreaterEqual,
            Value:    minAge,
        },
    })
    
    crit := criteria.Criteria{
        Filters: *filters,
    }
    
    // Usar en repositorio (la implementación aplicará los filtros)
    result := repo.FindByCriteria(ctx, crit)
    if result.IsOk() {
        users := result.Value()
        // Procesar usuarios
    }
}
```

## 📝 Mejores Prácticas

### Nombres de Campos

Usa nombres de campos consistentes. Los nombres deben coincidir con los campos en los modelos:

```go
// ✅ Correcto - nombre de campo coincide con modelo
filter := criteria.Filter{
    Field:    criteria.FilterField("email"), // Campo "email" en modelo
    Operator: criteria.OperatorEqual,
    Value:    "user@example.com",
}

// ❌ Incorrecto - nombre no coincide
filter := criteria.Filter{
    Field:    criteria.FilterField("user_email"), // Campo no existe
    Operator: criteria.OperatorEqual,
    Value:    "user@example.com",
}
```

### Valores de Filtros

Los valores deben ser del tipo correcto para el campo:

```go
// ✅ Correcto - tipo correcto
filter := criteria.Filter{
    Field:    criteria.FilterField("age"),
    Operator: criteria.OperatorGreaterThan,
    Value:    18, // int para campo numérico
}

// ✅ Correcto - string para campo de texto
filter := criteria.Filter{
    Field:    criteria.FilterField("email"),
    Operator: criteria.OperatorEqual,
    Value:    "user@example.com", // string
}
```

### Operadores LIKE

Para `OperatorLike`, el valor debe incluir wildcards (`%`):

```go
// ✅ Correcto - con wildcards
filter := criteria.Filter{
    Field:    criteria.FilterField("name"),
    Operator: criteria.OperatorLike,
    Value:    "%juan%", // Contiene "juan"
}

// Buscar al inicio
filter := criteria.Filter{
    Field:    criteria.FilterField("name"),
    Operator: criteria.OperatorLike,
    Value:    "juan%", // Empieza con "juan"
}

// Buscar al final
filter := criteria.Filter{
    Field:    criteria.FilterField("email"),
    Operator: criteria.OperatorLike,
    Value:    "%@example.com", // Termina con "@example.com"
}
```

## 🔗 Ver También

- [Persistence Overview](../README.md) - Visión general del módulo
- [Pagination](../pagination/README.md) - Paginación y ordenamiento
- [InMemory](../inmemory/README.md) - Implementación en memoria
- [Postgres](../postgres/README.md) - Implementación PostgreSQL

## 📚 Referencias

- [Código fuente](../../../persistence/criteria/)
- [Operadores SQL](https://www.w3schools.com/sql/sql_operators.asp)

---


## persistence/inmemory

**Source:** `docs/persistence/inmemory/README.md`

---

# InMemory Repository

`InMemory` proporciona una implementación de repositorio en memoria basada en mapas. Es ideal para testing, desarrollo y prototipado, ya que no requiere una base de datos externa. Implementa todas las operaciones CRUD, filtrado y paginación en memoria.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Ejemplos](#ejemplos)
- [Operaciones Disponibles](#operaciones-disponibles)
- [Limitaciones](#limitaciones)

## ✨ Características

- **Implementación completa**: CRUD, filtrado, paginación y ordenamiento
- **Sin dependencias externas**: No requiere base de datos
- **Thread-safe básico**: Para uso en testing
- **Filtrado en memoria**: Soporta todos los operadores de Criteria
- **Paginación y ordenamiento**: Implementación completa de Pageable
- **Conversión automática**: Entity ↔ Model automática

## 📦 Instalación

```bash
go get github.com/foundathyon/base/persistence/inmemory
```

## 🚀 Uso Básico

```go
import (
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/persistence/inmemory"
    "github.com/foundathyon/base/persistence/contracts"
)

// Definir tipos
type User struct {
    ID    string
    Name  string
    Email string
}
func (u User) GetID() string { return u.ID }

type UserModel struct {
    ID    string
    Name  string
    Email string
}
func (m UserModel) GetID() string { return m.ID }

// Crear repositorio
repo := inmemory.NewInMemoryMapRepository[User, UserModel]()

// Crear contexto
ctx := customctx.New(context.Background())

// Guardar entidad
user := User{ID: "123", Name: "Juan", Email: "juan@example.com"}
result := repo.Save(ctx, user)
```

## 📚 API

### Tipos

#### `InMemoryMapRepository[E IEntity, M IModel]`

```go
type InMemoryMapRepository[E contracts.IEntity, M contracts.IModel] struct {
    items map[string]M
}
```

`InMemoryMapRepository` es un repositorio genérico que almacena modelos en un mapa en memoria.

**Campos:**
- `items`: Mapa que almacena los modelos por ID

**Parámetros de tipo:**
- `E`: Tipo de entidad (debe implementar `IEntity`)
- `M`: Tipo de modelo (debe implementar `IModel`)

### Funciones de Construcción

#### `NewInMemoryMapRepository[E IEntity, M IModel]() *InMemoryMapRepository[E, M]`

Crea un nuevo repositorio en memoria vacío.

**Retorna:**
- `*InMemoryMapRepository[E, M]`: Nueva instancia de repositorio

**Ejemplo:**
```go
repo := inmemory.NewInMemoryMapRepository[User, UserModel]()
```

### Métodos

#### `Save(cc *customctx.CustomContext, item E) result.Result[E]`

Guarda o actualiza una entidad en el repositorio. Si ya existe un item con el mismo ID, lo reemplaza.

**Parámetros:**
- `cc`: CustomContext para acumular errores y logging
- `item`: Entidad a guardar

**Retorna:**
- `result.Result[E]`: Entidad guardada o error estructurado

**Errores posibles:**
- `ErrorItemIDRequired`: Si el ID de la entidad está vacío
- Errores de conversión Entity ↔ Model

**Ejemplo:**
```go
user := User{ID: "123", Name: "Juan", Email: "juan@example.com"}
result := repo.Save(ctx, user)
if result.IsOk() {
    savedUser := result.Value()
    fmt.Printf("Usuario guardado: %+v\n", savedUser)
}
```

#### `GetById(cc *customctx.CustomContext, id string) result.Result[E]`

Obtiene una entidad por su ID.

**Parámetros:**
- `cc`: CustomContext para acumular errores y logging
- `id`: ID de la entidad a buscar

**Retorna:**
- `result.Result[E]`: Entidad encontrada o error estructurado

**Errores posibles:**
- `ErrorItemIDRequired`: Si el ID está vacío
- `ErrorItemNotFound`: Si no se encuentra la entidad

**Ejemplo:**
```go
result := repo.GetById(ctx, "123")
if result.IsOk() {
    user := result.Value()
    fmt.Printf("Usuario: %+v\n", user)
} else {
    kerr := result.ToKError()
    fmt.Printf("Error: %s\n", kerr.Message)
}
```

#### `DeleteById(cc *customctx.CustomContext, id string) result.Result[E]`

Elimina una entidad por su ID.

**Parámetros:**
- `cc`: CustomContext para acumular errores y logging
- `id`: ID de la entidad a eliminar

**Retorna:**
- `result.Result[E]`: Result vacío (siempre `Empty[E]()`)

**Errores posibles:**
- `ErrorItemIDRequired`: Si el ID está vacío

**Nota:** Este método no retorna error si el item no existe, simplemente lo elimina si existe.

**Ejemplo:**
```go
result := repo.DeleteById(ctx, "123")
if result.IsEmpty() {
    fmt.Println("Item eliminado (o no existía)")
}
```

#### `UpdateByFields(cc *customctx.CustomContext, id string, fields map[string]any) result.Result[E]`

Actualiza campos específicos de una entidad sin reemplazarla completamente.

**Parámetros:**
- `cc`: CustomContext para acumular errores y logging
- `id`: ID de la entidad a actualizar
- `fields`: Mapa con los campos a actualizar

**Retorna:**
- `result.Result[E]`: Entidad actualizada o error estructurado

**Errores posibles:**
- `ErrorItemIDRequired`: Si el ID está vacío
- `ErrorItemFieldsRequired`: Si fields está vacío
- `ErrorItemNotFound`: Si no se encuentra la entidad
- Errores de conversión

**Ejemplo:**
```go
fields := map[string]any{
    "name":  "Juan Pérez",
    "email": "juan.perez@example.com",
}
result := repo.UpdateByFields(ctx, "123", fields)
if result.IsOk() {
    updatedUser := result.Value()
    fmt.Printf("Usuario actualizado: %+v\n", updatedUser)
}
```

#### `Matching(cc *customctx.CustomContext, crit criteria.Criteria, pageable *pagination.Pageable) result.Result[[]E]`

Busca entidades que coincidan con los criterios especificados, aplicando filtros, ordenamiento y paginación.

**Parámetros:**
- `cc`: CustomContext para acumular errores y logging
- `crit`: Criterios de búsqueda (filtros)
- `pageable`: Parámetros de paginación y ordenamiento

**Retorna:**
- `result.Result[[]E]`: Slice de entidades que coinciden o error estructurado

**Errores posibles:**
- `ErrorPageableRequired`: Si pageable es nil o inválido
- `ErrorApplyFilters`: Si hay error aplicando filtros
- Errores de conversión Model → Entity

**Ejemplo:**
```go
import (
    "github.com/foundathyon/base/persistence/criteria"
    "github.com/foundathyon/base/persistence/pagination"
)

// Crear criterios
filters := criteria.NewFilters([]criteria.Filter{
    {
        Field:    criteria.FilterField("email"),
        Operator: criteria.OperatorEqual,
        Value:    "juan@example.com",
    },
})
crit := criteria.Criteria{Filters: *filters}

// Crear paginación
pageReq := pagination.NewPageRequest(0, 10)
sort := pagination.NewSortASC("name")
pageable := pagination.NewPageable(pageReq, sort)

// Buscar
result := repo.Matching(ctx, crit, pageable)
if result.IsOk() {
    users := result.Value()
    fmt.Printf("Encontrados %d usuarios\n", len(users))
}
```

## 💡 Ejemplos

### Ejemplo 1: Operaciones CRUD Completas

```go
func main() {
    repo := inmemory.NewInMemoryMapRepository[User, UserModel]()
    ctx := customctx.New(context.Background())
    
    // CREATE - Guardar nueva entidad
    user := User{
        ID:    "user-123",
        Name:  "Juan Pérez",
        Email: "juan@example.com",
    }
    result := repo.Save(ctx, user)
    if result.IsOk() {
        fmt.Println("Usuario creado")
    }
    
    // READ - Obtener por ID
    result = repo.GetById(ctx, "user-123")
    if result.IsOk() {
        savedUser := result.Value()
        fmt.Printf("Usuario: %+v\n", savedUser)
    }
    
    // UPDATE - Actualizar campos
    fields := map[string]any{
        "name": "Juan Carlos Pérez",
    }
    result = repo.UpdateByFields(ctx, "user-123", fields)
    if result.IsOk() {
        updatedUser := result.Value()
        fmt.Printf("Usuario actualizado: %+v\n", updatedUser)
    }
    
    // DELETE - Eliminar
    result = repo.DeleteById(ctx, "user-123")
    fmt.Println("Usuario eliminado")
}
```

### Ejemplo 2: Búsqueda con Filtros

```go
func searchUsers(repo *inmemory.InMemoryMapRepository[User, UserModel], email string, minAge int) {
    ctx := customctx.New(context.Background())
    
    // Crear filtros
    filters := criteria.NewFilters([]criteria.Filter{
        {
            Field:    criteria.FilterField("email"),
            Operator: criteria.OperatorLike,
            Value:    "%" + email + "%",
        },
        {
            Field:    criteria.FilterField("age"),
            Operator: criteria.OperatorGreaterEqual,
            Value:    minAge,
        },
    })
    
    crit := criteria.Criteria{Filters: *filters}
    pageReq := pagination.NewPageRequest(0, 10)
    pageable := pagination.NewPageable(pageReq, nil)
    
    result := repo.Matching(ctx, crit, pageable)
    if result.IsOk() {
        users := result.Value()
        fmt.Printf("Encontrados %d usuarios\n", len(users))
        for _, user := range users {
            fmt.Printf("- %s (%s)\n", user.Name, user.Email)
        }
    }
}
```

### Ejemplo 3: Paginación y Ordenamiento

```go
func listUsersPaginated(repo *inmemory.InMemoryMapRepository[User, UserModel], page, size int) {
    ctx := customctx.New(context.Background())
    
    // Crear paginación con ordenamiento
    pageReq := pagination.NewPageRequest(page, size)
    sort := pagination.NewSortDESC("created_at") // Más recientes primero
    pageable := pagination.NewPageable(pageReq, sort)
    
    // Buscar sin filtros (todos los usuarios)
    emptyFilters := criteria.NewFilters([]criteria.Filter{})
    crit := criteria.Criteria{Filters: *emptyFilters}
    
    result := repo.Matching(ctx, crit, pageable)
    if result.IsOk() {
        users := result.Value()
        fmt.Printf("Página %d: %d usuarios\n", page, len(users))
        for _, user := range users {
            fmt.Printf("- %s\n", user.Name)
        }
    }
}
```

### Ejemplo 4: Testing

```go
func TestUserRepository(t *testing.T) {
    repo := inmemory.NewInMemoryMapRepository[User, UserModel]()
    ctx := customctx.New(context.Background())
    
    // Test Save
    user := User{ID: "test-123", Name: "Test User", Email: "test@example.com"}
    result := repo.Save(ctx, user)
    if !result.IsOk() {
        t.Fatalf("Error al guardar: %v", result.ToKError())
    }
    
    // Test GetById
    result = repo.GetById(ctx, "test-123")
    if !result.IsOk() {
        t.Fatalf("Error al obtener: %v", result.ToKError())
    }
    savedUser := result.Value()
    if savedUser.Name != "Test User" {
        t.Errorf("Nombre incorrecto: esperado 'Test User', obtenido '%s'", savedUser.Name)
    }
    
    // Test UpdateByFields
    fields := map[string]any{"name": "Updated User"}
    result = repo.UpdateByFields(ctx, "test-123", fields)
    if !result.IsOk() {
        t.Fatalf("Error al actualizar: %v", result.ToKError())
    }
    
    // Verificar actualización
    result = repo.GetById(ctx, "test-123")
    updatedUser := result.Value()
    if updatedUser.Name != "Updated User" {
        t.Errorf("Actualización falló: esperado 'Updated User', obtenido '%s'", updatedUser.Name)
    }
    
    // Test DeleteById
    result = repo.DeleteById(ctx, "test-123")
    if !result.IsEmpty() {
        t.Error("DeleteById debería retornar Empty")
    }
    
    // Verificar eliminación
    result = repo.GetById(ctx, "test-123")
    if result.IsOk() {
        t.Error("Usuario debería haber sido eliminado")
    }
}
```

## 📝 Mejores Prácticas

### Uso en Testing

`InMemoryMapRepository` es ideal para testing porque:
- No requiere configuración externa
- Es rápido (no hay I/O)
- Puede resetearse fácilmente

```go
func setupTestRepo() *inmemory.InMemoryMapRepository[User, UserModel] {
    return inmemory.NewInMemoryMapRepository[User, UserModel]()
}

func TestWithCleanRepo(t *testing.T) {
    repo := setupTestRepo()
    ctx := customctx.New(context.Background())
    
    // Test aislado con repositorio limpio
    // ...
}
```

### No Usar en Producción

**⚠️ Importante:** `InMemoryMapRepository` NO debe usarse en producción:
- Los datos se pierden al reiniciar la aplicación
- No es thread-safe para uso concurrente intensivo
- No escala más allá de la memoria disponible
- No soporta transacciones

**Usar solo para:**
- ✅ Testing unitario
- ✅ Desarrollo local
- ✅ Prototipado
- ✅ Ejemplos y demos

**Para producción:**
- ✅ Usar `PostgresRepository` o implementaciones similares
- ✅ Usar bases de datos persistentes

### Manejo de Errores

Siempre verificar errores y usar el contexto:

```go
// ✅ Correcto
ctx := customctx.New(parentContext)
result := repo.Save(ctx, user)
if !result.IsOk() {
    // Errores ya están en el contexto
    errors := ctx.Errors()
    // Procesar errores
}

// ❌ Incorrecto - ignorar errores
result := repo.Save(ctx, user)
user = result.Value() // Puede ser valor cero si hay error
```

## ⚠️ Limitaciones

1. **Persistencia**: Los datos se pierden al reiniciar la aplicación
2. **Concurrencia**: No es completamente thread-safe para escrituras concurrentes
3. **Escalabilidad**: Limitado por memoria disponible
4. **Rendimiento**: Filtrado y ordenamiento son O(n) en memoria
5. **Transacciones**: No soporta transacciones ACID

## 🔗 Ver También

- [Persistence Overview](../README.md) - Visión general del módulo
- [Contracts](../contracts/README.md) - Contratos de dominio
- [Criteria](../criteria/README.md) - Filtrado
- [Pagination](../pagination/README.md) - Paginación
- [Postgres](../postgres/README.md) - Implementación para producción

## 📚 Referencias

- [Código fuente](../../../persistence/inmemory/)
- [CustomContext](../../core/customctx/README.md) - Contexto usado por repositorios
- [Result](../../core/result/README.md) - Tipo Result retornado

---


## persistence/pagination

**Source:** `docs/persistence/pagination/README.md`

---

# Pagination

`Pagination` proporciona un sistema completo de paginación y ordenamiento para repositorios. Permite consultar datos de forma paginada, ordenar resultados y combinar paginación con criterios de búsqueda.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Ejemplos](#ejemplos)
- [Mejores Prácticas](#mejores-prácticas)

## ✨ Características

- **Paginación completa**: PageRequest, PageResult con metadata completa
- **Ordenamiento**: Sort con direcciones ASC/DESC
- **Combinación de criterios**: CriteriaPageable combina filtros y paginación
- **Validación automática**: Validación de parámetros con valores por defecto
- **Metadata rica**: Información completa sobre paginación (IsFirst, IsLast, HasNext, etc.)

## 📦 Instalación

```bash
go get github.com/foundathyon/base/persistence/pagination
```

## 🚀 Uso Básico

```go
import "github.com/foundathyon/base/persistence/pagination"

// Crear PageRequest
pageReq := pagination.NewPageRequest(0, 10) // Página 0, tamaño 10

// Crear Sort
sort := pagination.NewSortASC("name") // Ordenar por name ASC

// Crear Pageable
pageable := pagination.NewPageable(pageReq, sort)

// Crear PageResult
content := []User{...}
totalElements := int64(100)
result := pagination.NewPageResult(content, totalElements, pageReq)
```

## 📚 API

### Tipos

#### `PageRequest`

```go
type PageRequest struct {
    Page int // Número de página (0-indexed)
    Size int // Tamaño de la página
}
```

`PageRequest` representa una solicitud de paginación con número de página y tamaño.

**Campos:**
- `Page`: Número de página (empieza en 0)
- `Size`: Número de elementos por página

**Métodos:**
- `Offset() int`: Calcula el offset (desplazamiento)
- `Limit() int`: Retorna el límite (equivalente a Size)
- `IsValid() bool`: Verifica si la PageRequest es válida

**Funciones:**
- `NewPageRequest(page, size int) PageRequest`: Crea una PageRequest con validación

#### `PageResult[T]`

```go
type PageResult[T any] struct {
    Content       []T   // Elementos de la página actual
    TotalElements int64 // Total de elementos que coinciden
    TotalPages    int   // Total de páginas disponibles
    Page          int   // Número de página actual
    Size          int   // Tamaño de la página
    IsFirst       bool  // Es la primera página
    IsLast        bool  // Es la última página
    HasNext       bool  // Hay página siguiente
    HasPrevious   bool  // Hay página anterior
}
```

`PageResult` representa el resultado de una consulta paginada con metadata completa.

**Métodos:**
- `NewPageResult[T](content []T, totalElements int64, pageRequest PageRequest) PageResult[T]`: Crea un PageResult con valores calculados
- `EmptyPageResult[T](pageRequest PageRequest) PageResult[T]`: Crea un PageResult vacío

#### `Sort`

```go
type Sort struct {
    Field     string        // Campo por el cual ordenar
    Direction SortDirection // Dirección del ordenamiento (ASC o DESC)
}
```

`Sort` representa información de ordenamiento.

**Campos:**
- `Field`: Nombre del campo por el cual ordenar
- `Direction`: Dirección del ordenamiento (`SortDirectionASC` o `SortDirectionDESC`)

**Métodos:**
- `IsValid() bool`: Verifica si el Sort es válido

**Funciones:**
- `NewSort(field string, direction SortDirection) *Sort`: Crea un Sort con validación
- `NewSortASC(field string) *Sort`: Crea un Sort ascendente
- `NewSortDESC(field string) *Sort`: Crea un Sort descendente

#### `SortDirection`

```go
type SortDirection string

const (
    SortDirectionASC  SortDirection = "ASC"  // Ascendente
    SortDirectionDESC SortDirection = "DESC" // Descendente
)
```

#### `Pageable`

```go
type Pageable struct {
    PageRequest PageRequest // Parámetros de paginación
    Sort        *Sort       // Ordenamiento (opcional)
}
```

`Pageable` combina paginación con ordenamiento.

**Métodos:**
- `Offset() int`: Retorna el offset
- `Limit() int`: Retorna el límite
- `IsValid() bool`: Verifica si es válido

**Funciones:**
- `NewPageable(pageRequest PageRequest, sort *Sort) *Pageable`: Crea un Pageable
- `NewPageableWithoutSort(page, size int) *Pageable`: Crea un Pageable sin ordenamiento

#### `CriteriaPageable`

```go
type CriteriaPageable struct {
    Criteria criteria.Criteria // Filtros a aplicar
    Pageable *Pageable         // Paginación y ordenamiento
}
```

`CriteriaPageable` combina criterios de búsqueda con paginación y ordenamiento.

**Métodos:**
- `Offset() int`: Retorna el offset
- `Limit() int`: Retorna el límite
- `IsValid() bool`: Verifica si es válido

**Funciones:**
- `NewCriteriaPageable(criteria criteria.Criteria, pageable *Pageable) *CriteriaPageable`: Crea un CriteriaPageable
- `NewCriteriaPageableSimple(criteria criteria.Criteria, page, size int) *CriteriaPageable`: Crea sin ordenamiento

### Funciones Utilitarias

```go
// Validar y normalizar PageRequest
ValidatePageRequest(page, size int) PageRequest

// Calcular offset
CalculateOffset(page, size int) int

// Calcular total de páginas
CalculateTotalPages(totalElements, size int64) int

// Verificar si una página es válida
IsValidPage(page, totalPages int) bool
```

## 💡 Ejemplos

### Ejemplo 1: Paginación Básica

```go
import "github.com/foundathyon/base/persistence/pagination"

// Crear PageRequest (página 0, tamaño 10)
pageReq := pagination.NewPageRequest(0, 10)

// Obtener datos paginados
offset := pageReq.Offset() // 0
limit := pageReq.Limit()   // 10

// Consultar datos (ejemplo con repositorio)
users := repo.FindAll(offset, limit)
total := repo.Count()

// Crear resultado paginado
result := pagination.NewPageResult(users, total, pageReq)

fmt.Printf("Página: %d\n", result.Page)
fmt.Printf("Total elementos: %d\n", result.TotalElements)
fmt.Printf("Total páginas: %d\n", result.TotalPages)
fmt.Printf("Es primera: %v\n", result.IsFirst)
fmt.Printf("Es última: %v\n", result.IsLast)
fmt.Printf("Tiene siguiente: %v\n", result.HasNext)
fmt.Printf("Tiene anterior: %v\n", result.HasPrevious)
```

### Ejemplo 2: Paginación con Ordenamiento

```go
// Crear PageRequest
pageReq := pagination.NewPageRequest(0, 10)

// Crear Sort (ordenar por name ascendente)
sort := pagination.NewSortASC("name")

// Crear Pageable
pageable := pagination.NewPageable(pageReq, sort)

// Usar en repositorio
result := repo.FindAll(pageable)
if result.IsOk() {
    pageResult := result.Value()
    for _, user := range pageResult.Content {
        fmt.Printf("User: %+v\n", user)
    }
}
```

### Ejemplo 3: Ordenamiento Descendente

```go
// Ordenar por fecha de creación descendente
sort := pagination.NewSortDESC("created_at")

pageReq := pagination.NewPageRequest(0, 20)
pageable := pagination.NewPageable(pageReq, sort)

result := repo.FindAll(pageable)
```

### Ejemplo 4: Paginación sin Ordenamiento

```go
// Crear Pageable sin ordenamiento
pageable := pagination.NewPageableWithoutSort(0, 10)

// O simplemente crear PageRequest directamente
pageReq := pagination.NewPageRequest(0, 10)
```

### Ejemplo 5: Paginación con Criterios

```go
import (
    "github.com/foundathyon/base/persistence/criteria"
    "github.com/foundathyon/base/persistence/pagination"
)

// Crear criterios de búsqueda
filters := criteria.NewFilters([]criteria.Filter{
    {
        Field:    criteria.FilterField("status"),
        Operator: criteria.OperatorEqual,
        Value:    "active",
    },
})

crit := criteria.Criteria{
    Filters: *filters,
}

// Crear paginación con ordenamiento
pageReq := pagination.NewPageRequest(0, 10)
sort := pagination.NewSortASC("name")
pageable := pagination.NewPageable(pageReq, sort)

// Combinar criterios y paginación
criteriaPageable := pagination.NewCriteriaPageable(crit, pageable)

// Usar en repositorio
result := repo.FindByCriteria(ctx, criteriaPageable)
if result.IsOk() {
    pageResult := result.Value()
    // Procesar resultados paginados y filtrados
}
```

### Ejemplo 6: Navegación de Páginas

```go
func getPageResult(page int, size int) pagination.PageResult[User] {
    pageReq := pagination.NewPageRequest(page, size)
    
    users := repo.FindAll(pageReq.Offset(), pageReq.Limit())
    total := repo.Count()
    
    result := pagination.NewPageResult(users, total, pageReq)
    return result
}

// Primera página
result := getPageResult(0, 10)
fmt.Printf("IsFirst: %v, HasNext: %v\n", result.IsFirst, result.HasNext)

// Página intermedia
result = getPageResult(5, 10)
fmt.Printf("HasPrevious: %v, HasNext: %v\n", result.HasPrevious, result.HasNext)

// Última página
result = getPageResult(result.TotalPages-1, 10)
fmt.Printf("IsLast: %v, HasPrevious: %v\n", result.IsLast, result.HasPrevious)
```

### Ejemplo 7: Validación y Valores por Defecto

```go
// NewPageRequest valida y normaliza automáticamente
pageReq := pagination.NewPageRequest(-1, 0)
// page = 0 (normalizado)
// size = 10 (valor por defecto)

pageReq = pagination.NewPageRequest(5, 25)
// page = 5
// size = 25

// Validar antes de usar
if pageReq.IsValid() {
    // Usar pageReq
}

// Validar Sort
sort := pagination.NewSort("name", "INVALID")
// Se normaliza a ASC por defecto

if sort.IsValid() {
    // Usar sort
}
```

### Ejemplo 8: Funciones Utilitarias

```go
// Validar PageRequest
pageReq := pagination.ValidatePageRequest(0, 10)

// Calcular offset manualmente
offset := pagination.CalculateOffset(5, 20) // 100

// Calcular total de páginas
totalPages := pagination.CalculateTotalPages(150, 20) // 8

// Verificar si una página es válida
isValid := pagination.IsValidPage(5, 8) // true
isValid = pagination.IsValidPage(10, 8) // false
```

## 📝 Mejores Prácticas

### Convención de Paginación

**✅ Usar 0-indexed:**
- La primera página es 0
- Esto simplifica cálculos de offset

```go
// ✅ Correcto - primera página es 0
pageReq := pagination.NewPageRequest(0, 10)

// ❌ Incorrecto - usar 1 como primera página
pageReq := pagination.NewPageRequest(1, 10) // Genera offset incorrecto
```

### Tamaños de Página

Usa tamaños razonables:

```go
// ✅ Tamaños razonables
pageReq := pagination.NewPageRequest(0, 10)  // Pequeño
pageReq := pagination.NewPageRequest(0, 25)  // Mediano
pageReq := pagination.NewPageRequest(0, 50)  // Grande
pageReq := pagination.NewPageRequest(0, 100) // Muy grande (solo si es necesario)

// ❌ Evitar tamaños extremos
pageReq := pagination.NewPageRequest(0, 1)    // Muy pequeño
pageReq := pagination.NewPageRequest(0, 1000) // Muy grande (problemas de rendimiento)
```

### Ordenamiento

Siempre especifica un ordenamiento por defecto para consultas paginadas:

```go
// ✅ Ordenamiento explícito
sort := pagination.NewSortASC("created_at")
pageable := pagination.NewPageable(pageReq, sort)

// ❌ Sin ordenamiento (resultados inconsistentes)
pageable := pagination.NewPageableWithoutSort(0, 10)
```

### Manejo de Resultados Vacíos

```go
result := pagination.NewPageResult(users, 0, pageReq)

if len(result.Content) == 0 {
    fmt.Println("No hay resultados")
    return
}

// O usar EmptyPageResult
result := pagination.EmptyPageResult[User](pageReq)
if len(result.Content) == 0 {
    // Ya está vacío
}
```

### Metadata de Navegación

Usa la metadata de `PageResult` para UI:

```go
type PaginationMetadata struct {
    CurrentPage    int  `json:"current_page"`
    TotalPages     int  `json:"total_pages"`
    TotalElements  int64 `json:"total_elements"`
    PageSize       int  `json:"page_size"`
    IsFirst        bool `json:"is_first"`
    IsLast         bool `json:"is_last"`
    HasNext        bool `json:"has_next"`
    HasPrevious    bool `json:"has_previous"`
}

func toMetadata(result pagination.PageResult[User]) PaginationMetadata {
    return PaginationMetadata{
        CurrentPage:   result.Page,
        TotalPages:    result.TotalPages,
        TotalElements: result.TotalElements,
        PageSize:      result.Size,
        IsFirst:       result.IsFirst,
        IsLast:        result.IsLast,
        HasNext:       result.HasNext,
        HasPrevious:   result.HasPrevious,
    }
}
```

## 🔗 Ver También

- [Persistence Overview](../README.md) - Visión general del módulo
- [Criteria](../criteria/README.md) - Filtrado y búsqueda
- [InMemory](../inmemory/README.md) - Implementación en memoria
- [Postgres](../postgres/README.md) - Implementación PostgreSQL

## 📚 Referencias

- [Código fuente](../../../persistence/pagination/)
- [Spring Data Pagination](https://docs.spring.io/spring-data/commons/docs/current/reference/html/#repositories.query-methods.query-creation) - Inspiración para este diseño

---


## persistence/postgres

**Source:** `docs/persistence/postgres/README.md`

---

# Postgres Repository

`Postgres` proporciona una implementación de repositorio para PostgreSQL usando GORM. Es la implementación recomendada para producción, proporcionando persistencia real, transacciones y optimizaciones de base de datos.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Ejemplos](#ejemplos)
- [Configuración](#configuración)
- [Mejores Prácticas](#mejores-prácticas)

## ✨ Características

- **Persistencia real**: Los datos se guardan en PostgreSQL
- **Usa GORM**: ORM popular y maduro para Go
- **Transacciones**: Soporte para transacciones ACID
- **Optimizaciones SQL**: Filtrado y paginación en la base de datos
- **Manejo de errores**: Manejo específico de errores de PostgreSQL
- **Conversión automática**: Entity ↔ Model automática

## 📦 Instalación

```bash
go get github.com/foundathyon/base/persistence/postgres
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

## 🚀 Uso Básico

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/persistence/postgres"
)

// Configurar conexión a PostgreSQL
dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
    panic("Error conectando a PostgreSQL")
}

// Crear repositorio
repo := &postgres.PostgresRepository[User, UserModel]{
    Connection: db,
}

// Crear contexto
ctx := customctx.New(context.Background())

// Guardar entidad
user := User{ID: "123", Name: "Juan", Email: "juan@example.com"}
result := repo.Save(ctx, user)
```

## 📚 API

### Tipos

#### `PostgresRepository[E IEntity, M IModel]`

```go
type PostgresRepository[E contracts.IEntity, M contracts.IModel] struct {
    Connection *gorm.DB
}
```

`PostgresRepository` es un repositorio genérico que usa GORM para acceder a PostgreSQL.

**Campos:**
- `Connection`: Conexión GORM a PostgreSQL (`*gorm.DB`)

**Parámetros de tipo:**
- `E`: Tipo de entidad (debe implementar `IEntity`)
- `M`: Tipo de modelo (debe implementar `IModel`)

### Métodos

#### `Save(cc *customctx.CustomContext, item E) result.Result[E]`

Guarda una entidad en PostgreSQL. Si ya existe un item con el mismo ID y hay restricción única, retorna error de clave duplicada.

**Parámetros:**
- `cc`: CustomContext para acumular errores y logging
- `item`: Entidad a guardar

**Retorna:**
- `result.Result[E]`: Entidad guardada o error estructurado

**Errores posibles:**
- `ErrorItemIDRequired`: Si el ID de la entidad está vacío
- `ErrorDuplicateKey`: Si ya existe un registro con el mismo ID (clave única)
- `ErrorDatabaseOperation`: Si hay error en la operación de base de datos
- Errores de conversión Entity ↔ Model

**Ejemplo:**
```go
user := User{ID: "123", Name: "Juan", Email: "juan@example.com"}
result := repo.Save(ctx, user)
if result.IsOk() {
    savedUser := result.Value()
    fmt.Printf("Usuario guardado: %+v\n", savedUser)
} else {
    kerr := result.ToKError()
    if kerr.Code == 409 {
        fmt.Println("Usuario ya existe")
    }
}
```

#### `GetById(cc *customctx.CustomContext, id string) result.Result[E]`

Obtiene una entidad por su ID desde PostgreSQL.

**Parámetros:**
- `cc`: CustomContext para acumular errores y logging
- `id`: ID de la entidad a buscar

**Retorna:**
- `result.Result[E]`: Entidad encontrada o error estructurado

**Errores posibles:**
- `ErrorItemIDRequired`: Si el ID está vacío
- `ErrorItemNotFound`: Si no se encuentra la entidad (GORM `ErrRecordNotFound`)
- `ErrorDatabaseOperation`: Si hay error en la consulta

**Ejemplo:**
```go
result := repo.GetById(ctx, "123")
if result.IsOk() {
    user := result.Value()
    fmt.Printf("Usuario: %+v\n", user)
} else {
    kerr := result.ToKError()
    if kerr.Code == 404 {
        fmt.Println("Usuario no encontrado")
    }
}
```

#### `DeleteById(cc *customctx.CustomContext, id string) result.Result[E]`

Elimina una entidad por su ID desde PostgreSQL.

**Parámetros:**
- `cc`: CustomContext para acumular errores y logging
- `id`: ID de la entidad a eliminar

**Retorna:**
- `result.Result[E]`: Result vacío o error estructurado

**Errores posibles:**
- `ErrorItemIDRequired`: Si el ID está vacío
- `ErrorDatabaseOperation`: Si hay error en la operación

**Nota:** Este método no retorna error si el item no existe, simplemente lo elimina si existe.

#### `UpdateByFields(cc *customctx.CustomContext, id string, fields map[string]any) result.Result[E]`

Actualiza campos específicos de una entidad sin reemplazarla completamente.

**Parámetros:**
- `cc`: CustomContext para acumular errores y logging
- `id`: ID de la entidad a actualizar
- `fields`: Mapa con los campos a actualizar

**Retorna:**
- `result.Result[E]`: Entidad actualizada o error estructurado

**Errores posibles:**
- `ErrorItemIDRequired`: Si el ID está vacío
- `ErrorItemFieldsRequired`: Si fields está vacío
- `ErrorItemNotFound`: Si no se encuentra la entidad
- `ErrorDatabaseOperation`: Si hay error en la actualización
- Errores de conversión

#### `Matching(cc *customctx.CustomContext, crit criteria.Criteria, pageable *pagination.Pageable) result.Result[[]E]`

Busca entidades que coincidan con los criterios especificados, aplicando filtros, ordenamiento y paginación usando SQL.

**Parámetros:**
- `cc`: CustomContext para acumular errores y logging
- `crit`: Criterios de búsqueda (filtros)
- `pageable`: Parámetros de paginación y ordenamiento

**Retorna:**
- `result.Result[[]E]`: Slice de entidades que coinciden o error estructurado

**Errores posibles:**
- `ErrorPageableRequired`: Si pageable es nil o inválido
- `ErrorDatabaseOperation`: Si hay error en la consulta SQL
- Errores de conversión Model → Entity

### Errores Predefinidos

```go
var (
    ErrorItemIDRequired     = kerrors.NewKError("Item ID is required", 400, nil)
    ErrorItemNotFound       = kerrors.NewKError("Item not found", 404, nil)
    ErrorItemFieldsRequired = kerrors.NewKError("Item fields are required", 400, nil)
    ErrorDatabaseOperation  = kerrors.NewKError("Database operation failed", 500, nil)
    ErrorDuplicateKey       = kerrors.NewKError("Duplicate key error: a record with the same unique key already exists", 409, nil)
    ErrorPageableRequired   = kerrors.NewKError("Pageable is required and must be valid", 400, nil)
    // Errores de conversión...
)
```

## 💡 Ejemplos

### Ejemplo 1: Configuración Inicial

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/foundathyon/base/persistence/postgres"
)

func setupRepository() (*postgres.PostgresRepository[User, UserModel], error) {
    // Configurar DSN
    dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"
    
    // Conectar a PostgreSQL
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Auto-migrate (opcional, para crear tablas)
    db.AutoMigrate(&UserModel{})
    
    // Crear repositorio
    repo := &postgres.PostgresRepository[User, UserModel]{
        Connection: db,
    }
    
    return repo, nil
}
```

### Ejemplo 2: Operaciones CRUD

```go
func main() {
    repo, err := setupRepository()
    if err != nil {
        panic(err)
    }
    
    ctx := customctx.New(context.Background())
    
    // CREATE - Guardar nueva entidad
    user := User{
        ID:    "user-123",
        Name:  "Juan Pérez",
        Email: "juan@example.com",
    }
    result := repo.Save(ctx, user)
    if result.IsOk() {
        fmt.Println("Usuario creado")
    } else {
        kerr := result.ToKError()
        if kerr.Code == 409 {
            fmt.Println("Usuario ya existe")
        }
    }
    
    // READ - Obtener por ID
    result = repo.GetById(ctx, "user-123")
    if result.IsOk() {
        savedUser := result.Value()
        fmt.Printf("Usuario: %+v\n", savedUser)
    }
    
    // UPDATE - Actualizar campos
    fields := map[string]any{
        "name": "Juan Carlos Pérez",
    }
    result = repo.UpdateByFields(ctx, "user-123", fields)
    if result.IsOk() {
        updatedUser := result.Value()
        fmt.Printf("Usuario actualizado: %+v\n", updatedUser)
    }
    
    // DELETE - Eliminar
    result = repo.DeleteById(ctx, "user-123")
    fmt.Println("Usuario eliminado")
}
```

### Ejemplo 3: Búsqueda con Filtros y Paginación

```go
import (
    "github.com/foundathyon/base/persistence/criteria"
    "github.com/foundathyon/base/persistence/pagination"
)

func searchUsers(repo *postgres.PostgresRepository[User, UserModel], email string) {
    ctx := customctx.New(context.Background())
    
    // Crear filtros
    filters := criteria.NewFilters([]criteria.Filter{
        {
            Field:    criteria.FilterField("email"),
            Operator: criteria.OperatorLike,
            Value:    "%" + email + "%",
        },
    })
    
    crit := criteria.Criteria{Filters: *filters}
    
    // Crear paginación con ordenamiento
    pageReq := pagination.NewPageRequest(0, 10)
    sort := pagination.NewSortASC("name")
    pageable := pagination.NewPageable(pageReq, sort)
    
    // Buscar (los filtros se aplican en SQL)
    result := repo.Matching(ctx, crit, pageable)
    if result.IsOk() {
        users := result.Value()
        fmt.Printf("Encontrados %d usuarios\n", len(users))
        for _, user := range users {
            fmt.Printf("- %s (%s)\n", user.Name, user.Email)
        }
    }
}
```

### Ejemplo 4: Manejo de Errores Específicos

```go
func saveUserWithErrorHandling(repo *postgres.PostgresRepository[User, UserModel], user User) {
    ctx := customctx.New(context.Background())
    
    result := repo.Save(ctx, user)
    if !result.IsOk() {
        kerr := result.ToKError()
        
        switch kerr.Code {
        case 400:
            fmt.Println("Datos inválidos:", kerr.Message)
        case 409:
            fmt.Println("Usuario ya existe (clave duplicada)")
            // Intentar actualizar en lugar de crear
            updateResult := repo.UpdateByFields(ctx, user.ID, map[string]any{
                "name":  user.Name,
                "email": user.Email,
            })
            // ...
        case 500:
            fmt.Println("Error de base de datos:", kerr.Message)
            if kerr.Cause != nil {
                fmt.Printf("Causa: %v\n", kerr.Cause)
            }
        default:
            fmt.Printf("Error desconocido: %s\n", kerr.Message)
        }
        
        // Errores también están en el contexto
        allErrors := ctx.Errors()
        fmt.Printf("Total de errores en contexto: %d\n", len(allErrors))
    }
}
```

### Ejemplo 5: Transacciones (con GORM)

```go
func saveUserWithTransaction(repo *postgres.PostgresRepository[User, UserModel], user User, profile Profile) error {
    ctx := customctx.New(context.Background())
    
    // Iniciar transacción
    err := repo.Connection.Transaction(func(tx *gorm.DB) error {
        // Crear repositorio temporal con transacción
        txRepo := &postgres.PostgresRepository[User, UserModel]{
            Connection: tx,
        }
        
        // Guardar usuario
        result := txRepo.Save(ctx, user)
        if !result.IsOk() {
            return result.ToKError()
        }
        
        // Guardar perfil (ejemplo con otro repositorio)
        // profileResult := profileRepo.Save(ctx, profile)
        // ...
        
        return nil // Commit automático
    })
    
    if err != nil {
        // Rollback automático si hay error
        fmt.Printf("Error en transacción: %v\n", err)
        return err
    }
    
    return nil
}
```

## ⚙️ Configuración

### DSN (Data Source Name)

```go
// Formato básico
dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"

// Con SSL
dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=require"

// Desde variables de entorno
dsn := fmt.Sprintf(
    "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_NAME"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_SSLMODE"),
)
```

### GORM Config

```go
config := &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info), // Log SQL queries
    NamingStrategy: schema.NamingStrategy{
        TablePrefix:   "app_",      // Prefijo para tablas
        SingularTable: true,        // Usar nombres singulares
    },
    PrepareStmt: true,              // Preparar statements (mejor rendimiento)
}

db, err := gorm.Open(postgres.Open(dsn), config)
```

### Auto-Migrate

```go
// Migrar modelo a tabla
err := db.AutoMigrate(&UserModel{})
if err != nil {
    panic("Error en migración: " + err.Error())
}

// Migrar múltiples modelos
err := db.AutoMigrate(
    &UserModel{},
    &ProductModel{},
    &OrderModel{},
)
```

## 📝 Mejores Prácticas

### Modelos con GORM Tags

Usa tags de GORM apropiados en los modelos:

```go
type UserModel struct {
    ID        string    `gorm:"type:uuid;primary_key" json:"id"`
    Name      string    `gorm:"type:varchar(255);not null" json:"name"`
    Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    IsRemoved bool      `gorm:"type:boolean;default:false;index" json:"is_removed"`
}
```

### Manejo de Errores

Siempre maneja errores específicos:

```go
// ✅ Correcto - verificar tipo de error
result := repo.Save(ctx, user)
if !result.IsOk() {
    kerr := result.ToKError()
    switch kerr.Code {
    case 409:
        // Clave duplicada - intentar actualizar
    case 500:
        // Error de BD - registrar y notificar
    default:
        // Otros errores
    }
}

// ❌ Incorrecto - ignorar errores
result := repo.Save(ctx, user)
user = result.Value() // Puede ser valor cero
```

### Uso de Transacciones

Usa transacciones para operaciones relacionadas:

```go
// ✅ Correcto - usar transacciones
err := db.Transaction(func(tx *gorm.DB) error {
    // Múltiples operaciones
    return nil // Commit
})

// ❌ Incorrecto - operaciones no relacionadas en transacción
// Usar transacciones solo cuando sea necesario
```

### Índices para Rendimiento

Define índices apropiados:

```go
type UserModel struct {
    Email string `gorm:"index"`              // Índice simple
    Age   int    `gorm:"index:idx_age_name"` // Índice compuesto
    Name  string `gorm:"index:idx_age_name"`
}
```

### Conexión Pool

Configura el pool de conexiones:

```go
sqlDB, err := db.DB()
if err != nil {
    return err
}

// Configurar pool
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

## ⚠️ Consideraciones

1. **Claves duplicadas**: `Save()` usa `Create()`, que falla si existe. Considera usar `Save()` o `FirstOrCreate()` de GORM directamente para upsert
2. **Transacciones**: Usa `db.Transaction()` para operaciones que deben ser atómicas
3. **Migraciones**: Usa herramientas de migración apropiadas (migrate, golang-migrate) en producción
4. **Rendimiento**: Los filtros se aplican en SQL, pero asegúrate de tener índices apropiados

## 🔗 Ver También

- [Persistence Overview](../README.md) - Visión general del módulo
- [Contracts](../contracts/README.md) - Contratos de dominio
- [Criteria](../criteria/README.md) - Filtrado
- [Pagination](../pagination/README.md) - Paginación
- [InMemory](../inmemory/README.md) - Implementación para testing
- [GORM Documentation](https://gorm.io/docs/)

## 📚 Referencias

- [Código fuente](../../../persistence/postgres/)
- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [CustomContext](../../core/customctx/README.md) - Contexto usado por repositorios
- [Result](../../core/result/README.md) - Tipo Result retornado

---


## persistence

**Source:** `docs/persistence/README.md`

---

# Persistence

`Persistence` proporciona abstracciones y implementaciones para persistencia de datos en el base-kit. Sigue los principios de Clean Architecture, separando las interfaces del dominio de las implementaciones de infraestructura, permitiendo múltiples implementaciones (in-memory, PostgreSQL, etc.) mientras mantiene el código de dominio independiente de detalles de almacenamiento.

## 📋 Tabla de Contenidos

- [Arquitectura](#arquitectura)
- [Conceptos Fundamentales](#conceptos-fundamentales)
- [Instalación](#instalación)
- [Módulos](#módulos)
- [Casos de Uso](#casos-de-uso)
- [Mejores Prácticas](#mejores-prácticas)

## 🏗️ Arquitectura

### Principios

El módulo de persistencia sigue Clean Architecture y Dependency Inversion:

- **Interfaces en el dominio**: `IEntity` y `IModel` definen contratos
- **Separación de capas**: Dominio, Contratos, Infraestructura
- **Independencia de implementación**: El dominio no conoce detalles de almacenamiento
- **Adaptadores**: Las implementaciones adaptan los contratos a tecnologías específicas

### Estructura

```
persistence/
├── contracts/          # Contratos del dominio (IEntity, IModel)
├── criteria/           # Filtrado y búsqueda
├── pagination/         # Paginación y ordenamiento
├── inmemory/           # Implementación en memoria (testing/development)
└── postgres/           # Implementación PostgreSQL/GORM (producción)
```

## 💡 Conceptos Fundamentales

### Entity vs Model

- **Entity (`IEntity`)**: Representa la entidad de dominio. Es independiente de cómo se almacena.
- **Model (`IModel`)**: Representa cómo se almacena la entidad. Incluye campos técnicos de infraestructura.

**Separación:**
- Las entidades son parte del dominio de negocio
- Los modelos son parte de la capa de infraestructura
- Las conversiones entre Entity y Model se realizan en la capa de adaptadores

### Repository Pattern

Los repositorios actúan como adaptadores entre el dominio y la infraestructura:

- Implementan operaciones CRUD sobre entidades
- Manejan conversiones Entity ↔ Model
- Proporcionan métodos para búsqueda, filtrado y paginación
- Manejan errores y los agregan al contexto

## 📦 Instalación

```bash
go get github.com/foundathyon/base/persistence/contracts
go get github.com/foundathyon/base/persistence/criteria
go get github.com/foundathyon/base/persistence/pagination
go get github.com/foundathyon/base/persistence/inmemory
go get github.com/foundathyon/base/persistence/postgres
```

## 📚 Módulos

### [Contracts](../persistence/contracts/README.md)

Define las interfaces fundamentales para entidades y modelos:
- `IEntity`: Interface para entidades de dominio
- `IModel`: Interface para modelos de persistencia
- Funciones de conversión Entity ↔ Model

### [Criteria](../persistence/criteria/README.md)

Sistema de filtrado y búsqueda:
- `Criteria`: Contenedor de filtros
- `Filter`: Filtros individuales con operadores
- `Filters`: Colección de filtros

### [Pagination](../persistence/pagination/README.md)

Sistema de paginación y ordenamiento:
- `Pageable`: Paginación con criterios
- `PageRequest`: Parámetros de paginación
- `PageResult`: Resultado paginado
- `Sort`: Ordenamiento

### [InMemory](../persistence/inmemory/README.md)

Implementación en memoria para testing y desarrollo:
- `InMemoryMapRepository`: Repositorio basado en mapas
- Operaciones CRUD básicas
- Filtrado y paginación en memoria

### [Postgres](../persistence/postgres/README.md)

Implementación PostgreSQL usando GORM:
- `PostgresRepository`: Repositorio con GORM
- Operaciones CRUD con base de datos
- Filtrado y paginación con SQL

## 🎯 Casos de Uso

### Definir una Entidad

```go
package domain

import "github.com/foundathyon/base/persistence/contracts"

type User struct {
    ID    string
    Name  string
    Email string
}

func (u User) GetID() string {
    return u.ID
}
```

### Definir un Modelo

```go
package infrastructure

import (
    "time"
    "github.com/foundathyon/base/persistence/contracts"
)

type UserModel struct {
    ID        string    `gorm:"primary_key" json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    IsRemoved bool      `json:"is_removed"`
}

func (m UserModel) GetID() string {
    return m.ID
}
```

### Usar un Repositorio

```go
import (
    "github.com/foundathyon/base/core/customctx"
    "github.com/foundathyon/base/persistence/inmemory"
)

// Crear repositorio
repo := inmemory.NewInMemoryMapRepository[domain.User, infrastructure.UserModel]()

// Crear contexto
ctx := customctx.New(context.Background())

// Guardar entidad
user := domain.User{
    ID:    "user-123",
    Name:  "Juan Pérez",
    Email: "juan@example.com",
}

result := repo.Save(ctx, user)
if result.IsOk() {
    savedUser := result.Value()
    fmt.Printf("Usuario guardado: %+v\n", savedUser)
} else {
    fmt.Printf("Error: %s\n", result.ToKError().Message)
}
```

## 📝 Mejores Prácticas

### Separación de Dominio e Infraestructura

**✅ Correcto:**
```go
// Dominio
type User struct {
    ID    string
    Name  string
    Email string
}

// Infraestructura
type UserModel struct {
    ID        string
    Name      string
    Email     string
    CreatedAt time.Time  // Campo técnico
    UpdatedAt time.Time  // Campo técnico
    IsRemoved bool       // Campo técnico
}
```

**❌ Incorrecto:**
```go
// No mezclar campos técnicos en entidades de dominio
type User struct {
    ID        string
    Name      string
    Email     string
    CreatedAt time.Time  // ❌ Campo técnico en dominio
}
```

### Uso de Context

Siempre usar `CustomContext` para:
- Acumular errores durante operaciones
- Agregar metadata técnica (request IDs, etc.)
- Logging contextual

```go
ctx := customctx.New(parentContext)
ctx = ctx.WithValue("request_id", "req-123")

result := repo.Save(ctx, entity)
if !result.IsOk() {
    // Errores ya están en el contexto
    errors := ctx.Errors()
    // Procesar errores
}
```

### Manejo de Errores

Los repositorios automáticamente:
- Validan entradas
- Agregan errores al contexto
- Retornan `Result[T]` para manejo funcional

```go
result := repo.GetById(ctx, "")
if !result.IsOk() {
    kerr := result.ToKError()
    // Error ya está en el contexto también
    // Puedes acceder con ctx.Errors()
}
```

## 🔗 Ver También

- [Contracts](../persistence/contracts/README.md) - Contratos de dominio
- [Criteria](../persistence/criteria/README.md) - Filtrado
- [Pagination](../persistence/pagination/README.md) - Paginación
- [InMemory](../persistence/inmemory/README.md) - Implementación en memoria
- [Postgres](../persistence/postgres/README.md) - Implementación PostgreSQL

## 📚 Referencias

- [CustomContext](../core/customctx/README.md) - Contexto personalizado
- [Result](../core/result/README.md) - Tipo Result
- [KErrors](../core/kerrors/README.md) - Errores estructurados

---

