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
