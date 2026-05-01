package main

import (
	"context"
	"fmt"
	"time"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
)

func main() {
	fmt.Printf("=== Ejemplos del paquete customctx ===\n\n")

	// Ejemplo 1: Crear y usar un CustomContext básico
	exampleBasicContext()

	// Ejemplo 2: Acumular múltiples errores
	exampleMultipleErrors()

	// Ejemplo 3: Usar con context estándar
	exampleWithStandardContext()

	// Ejemplo 4: Caso de uso práctico - Validación de formulario
	exampleFormValidation()

	// Ejemplo 5: Uso concurrente
	exampleConcurrentUse()

	// Ejemplo 6: Almacenamiento de valores (metadata técnica)
	exampleValueStorage()
}

// exampleBasicContext muestra cómo crear y usar un CustomContext básico
func exampleBasicContext() {
	fmt.Println("1. CustomContext básico:")
	parent := context.Background()
	ctx := customctx.New(parent)

	// Agregar un error
	err := kerrors.NewKError("Error de configuración", 500, map[string]any{
		"config_key": "database_url",
	})
	ctx.AddError(err)

	fmt.Printf("   ¿Tiene errores? %v\n", ctx.HasErrors())
	fmt.Printf("   Número de errores: %d\n", len(ctx.Errors()))

	if ctx.HasErrors() {
		firstErr := ctx.FirstError()
		fmt.Printf("   Primer error: %s\n", firstErr.Error.Message)
		fmt.Printf("   Registrado en: %s\n", firstErr.CallIn)
	}
	fmt.Println()
}

// exampleMultipleErrors muestra cómo acumular múltiples errores
func exampleMultipleErrors() {
	fmt.Println("2. Acumulación de múltiples errores:")

	ctx := customctx.New(context.Background())

	// Agregar varios errores
	ctx.AddError(kerrors.NewKError("Campo 'nombre' es requerido", 400, map[string]any{
		"field": "nombre",
	}))
	ctx.AddError(kerrors.NewKError("Campo 'email' tiene formato inválido", 400, map[string]any{
		"field": "email",
		"value": "invalid-email",
	}))
	ctx.AddError(kerrors.NewKError("Campo 'edad' debe ser mayor a 0", 400, map[string]any{
		"field": "edad",
		"value": -5,
	}))

	fmt.Printf("   Total de errores: %d\n", len(ctx.Errors()))
	fmt.Printf("   Errores acumulados:\n")
	for i, wrapErr := range ctx.Errors() {
		fmt.Printf("     %d. %s (Código: %d)\n", i+1, wrapErr.Error.Message, wrapErr.Error.Code)
		if wrapErr.Error.Metadata != nil {
			fmt.Printf("        Metadata: %v\n", wrapErr.Error.Metadata)
		}
		fmt.Printf("        Registrado en: %s\n", wrapErr.CallIn)
	}

	fmt.Printf("   Primer error: %s\n", ctx.FirstError().Error.Message)
	fmt.Printf("   Último error: %s\n", ctx.LastError().Error.Message)
	fmt.Println()
}

// exampleWithStandardContext muestra cómo usar CustomContext con context estándar
func exampleWithStandardContext() {
	fmt.Println("3. Integración con context estándar:")

	// Crear un context con timeout
	parent, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Crear un context con un valor
	parent = context.WithValue(parent, "user_id", 123)

	ctx := customctx.New(parent)

	// Acceder a valores del context padre
	userID := ctx.Value("user_id")
	fmt.Printf("   User ID del context: %v\n", userID)

	// Verificar deadline
	deadline, ok := ctx.Deadline()
	if ok {
		fmt.Printf("   Deadline: %v\n", deadline)
	}

	// Agregar errores mientras se usa el context
	ctx.AddError(kerrors.NewKError("Error al procesar", 500, nil))

	fmt.Printf("   ¿Tiene errores? %v\n", ctx.HasErrors())
	fmt.Println()
}

// exampleFormValidation muestra un caso de uso práctico de validación
func exampleFormValidation() {
	fmt.Println("4. Caso de uso - Validación de formulario:")

	type FormData struct {
		Nombre string
		Email  string
		Edad   int
	}

	validateForm := func(ctx *customctx.CustomContext, form FormData) {
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
		} else if len(form.Email) < 5 {
			ctx.AddError(kerrors.NewKError(
				"El email es demasiado corto",
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

	// Caso 1: Formulario válido
	fmt.Println("   Caso 1: Formulario válido")
	ctx1 := customctx.New(context.Background())
	validateForm(ctx1, FormData{
		Nombre: "Juan",
		Email:  "juan@example.com",
		Edad:   25,
	})
	if ctx1.HasErrors() {
		fmt.Printf("     ✗ Errores encontrados: %d\n", len(ctx1.Errors()))
	} else {
		fmt.Printf("     ✓ Formulario válido\n")
	}

	// Caso 2: Formulario con errores
	fmt.Println("   Caso 2: Formulario con errores")
	ctx2 := customctx.New(context.Background())
	validateForm(ctx2, FormData{
		Nombre: "",
		Email:  "a@b",
		Edad:   15,
	})
	if ctx2.HasErrors() {
		fmt.Printf("     ✗ Errores encontrados: %d\n", len(ctx2.Errors()))
		for i, wrapErr := range ctx2.Errors() {
			fmt.Printf("       %d. %s\n", i+1, wrapErr.Error.Message)
		}
	}
	fmt.Println()
}

// exampleConcurrentUse muestra el uso concurrente del CustomContext
func exampleConcurrentUse() {
	fmt.Println("5. Uso concurrente:")

	ctx := customctx.New(context.Background())

	// Simular múltiples goroutines agregando errores
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(id int) {
			ctx.AddError(kerrors.NewKError(
				fmt.Sprintf("Error de goroutine %d", id),
				500,
				map[string]any{"goroutine_id": id},
			))
			done <- true
		}(i)
	}

	// Esperar a que todas las goroutines terminen
	for i := 0; i < 5; i++ {
		<-done
	}

	fmt.Printf("   Total de errores (concurrentes): %d\n", len(ctx.Errors()))
	fmt.Printf("   ¿Tiene errores? %v\n", ctx.HasErrors())

	// Limpiar errores
	ctx.Clear()
	fmt.Printf("   Después de Clear(), ¿tiene errores? %v\n", ctx.HasErrors())
	fmt.Println()
}

// exampleValueStorage muestra cómo almacenar valores técnicos en CustomContext
func exampleValueStorage() {
	fmt.Println("6. Almacenamiento de valores (metadata técnica):")

	ctxOriginal := customctx.New(context.Background())

	// Almacenar metadata técnica (permitido) - creando nuevos contextos inmutables
	// Usar strings como claves para identificar valores
	ctx := ctxOriginal.WithValue("request_id", "req-abc-123")
	ctx = ctx.WithValue("trace_id", "trace-xyz-456")
	ctx = ctx.WithValue("user_id", "user-789") // ID técnico para logging/tracing

	fmt.Println("   Valores almacenados:")
	fmt.Printf("     Request ID: %v\n", ctx.GetValue("request_id"))
	fmt.Printf("     Trace ID: %v\n", ctx.GetValue("trace_id"))
	fmt.Printf("     User ID (técnico): %v\n", ctx.GetValue("user_id"))

	// Verificar existencia de valores
	if ctx.HasValue("request_id") {
		fmt.Println("   ✓ Request ID existe en el contexto")
	}

	// El método Value() busca primero en CustomContext (para string keys), luego en el parent
	fmt.Println("   Usando Value() (busca en CustomContext y parent para string keys):")
	fmt.Printf("     Request ID via Value(): %v\n", ctx.Value("request_id"))

	// Inmutabilidad: crear nuevos contextos con más valores
	ctx1 := ctxOriginal.WithValue("request_id", "req-123")
	ctx2 := ctx1.WithValue("trace_id", "trace-456")

	fmt.Println("   Inmutabilidad:")
	fmt.Printf("     ctx original - Request ID: %v (debe ser nil)\n", ctxOriginal.GetValue("request_id"))
	fmt.Printf("     ctx1 - Request ID: %v\n", ctx1.GetValue("request_id"))
	fmt.Printf("     ctx2 - Request ID: %v, Trace ID: %v\n",
		ctx2.GetValue("request_id"), ctx2.GetValue("trace_id"))

	// Integración con parent context
	parent := context.WithValue(context.Background(), "parent_key", "parent_value")
	ctxWithParent := customctx.New(parent)
	ctxWithParent = ctxWithParent.WithValue("request_id", "custom_value")

	fmt.Println("   Integración con context padre:")
	fmt.Printf("     Valor del CustomContext (string key): %v\n", ctxWithParent.Value("request_id"))
	fmt.Printf("     Valor del parent context (string key): %v\n", ctxWithParent.Value("parent_key"))

	// Valor() busca primero en CustomContext (para string keys), luego en parent
	if ctxWithParent.Value("request_id") != "custom_value" {
		fmt.Println("     ⚠ CustomContext tiene prioridad sobre parent para string keys")
	}

	fmt.Println()
	fmt.Println("   ⚠ IMPORTANTE: Solo almacenar metadata técnica,")
	fmt.Println("     no datos de negocio (entidades, estados, etc.)")
	fmt.Println("     Usar strings descriptivos como claves (ej: 'request_id', 'trace_id')")
	fmt.Println()
}
