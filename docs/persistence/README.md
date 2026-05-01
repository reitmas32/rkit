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
