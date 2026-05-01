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
