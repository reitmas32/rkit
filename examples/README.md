# Ejemplos de Base Kit

Esta carpeta contiene ejemplos de uso para cada módulo del paquete base-kit, organizados en carpetas separadas dentro de `core`.

## Estructura

```
examples/
└── core/
    ├── result/
    │   └── result_example.go      # Ejemplos del módulo result
    ├── kerrors/
    │   └── kerrors_example.go     # Ejemplos del módulo kerrors
    └── customctx/
        └── customctx_example.go   # Ejemplos del módulo customctx
```

## Cómo ejecutar los ejemplos

Cada carpeta contiene un programa independiente con su propia función `main`. Para ejecutarlos:

### Ejemplo de result

```bash
go run examples/core/result/result_example.go
```

O desde la carpeta:

```bash
cd examples/core/result
go run result_example.go
```

### Ejemplo de kerrors

```bash
go run examples/core/kerrors/kerrors_example.go
```

O desde la carpeta:

```bash
cd examples/core/kerrors
go run kerrors_example.go
```

### Ejemplo de customctx

```bash
go run examples/core/customctx/customctx_example.go
```

O desde la carpeta:

```bash
cd examples/core/customctx
go run customctx_example.go
```

## Descripción de los ejemplos

### core/result/result_example.go

Muestra cómo usar el tipo `Result[T]` para manejar valores y errores de forma funcional:

- **Ejemplo 1**: Crear Results exitosos
- **Ejemplo 2**: Crear Results con errores
- **Ejemplo 3**: Verificar si un Result es exitoso
- **Ejemplo 4**: Usar Result con diferentes tipos (int, string, structs, slices)
- **Ejemplo 5**: Manejar errores estructurados con KError

### core/kerrors/kerrors_example.go

Demuestra el uso de errores estructurados con `KError`:

- **Ejemplo 1**: Crear errores básicos
- **Ejemplo 2**: Crear errores con metadata
- **Ejemplo 3**: Encadenar errores para preservar el contexto
- **Ejemplo 4**: Usar `errors.Is` y `errors.Unwrap` para inspeccionar cadenas de errores
- **Ejemplo 5**: Casos de uso prácticos de validación

### core/customctx/customctx_example.go

Ilustra cómo usar `CustomContext` para acumular errores durante la ejecución:

- **Ejemplo 1**: Crear y usar CustomContext básico
- **Ejemplo 2**: Acumular múltiples errores
- **Ejemplo 3**: Integración con context estándar de Go
- **Ejemplo 4**: Caso de uso práctico: validación de formularios
- **Ejemplo 5**: Uso concurrente seguro
