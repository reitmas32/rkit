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
