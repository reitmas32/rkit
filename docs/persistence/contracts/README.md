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
