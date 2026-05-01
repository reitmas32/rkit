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
