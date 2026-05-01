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
