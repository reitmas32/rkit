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
