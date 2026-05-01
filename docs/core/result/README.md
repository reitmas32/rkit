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
