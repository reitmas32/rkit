package main

import (
	"fmt"

	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
)

func main() {
	fmt.Printf("=== Ejemplos del paquete result ===\n\n")

	// Ejemplo 1: Crear un Result exitoso
	exampleOkResult()

	// Ejemplo 2: Crear un Result con error
	exampleErrResult()

	// Ejemplo 3: Verificar si un Result es exitoso
	exampleIsOk()

	// Ejemplo 4: Usar Result con diferentes tipos
	exampleWithDifferentTypes()

	// Ejemplo 5: Manejo de errores estructurados
	exampleWithStructuredErrors()
}

// exampleOkResult muestra cómo crear y usar un Result exitoso
func exampleOkResult() {
	fmt.Println("1. Result exitoso:")
	r := result.Ok(42)
	fmt.Printf("   Result: %+v\n", r)
	fmt.Printf("   IsOk: %v\n", r.IsOk())
	fmt.Printf("   Value: %d\n", r.Value())
	fmt.Printf("   Error: %v\n", r.Error())
	fmt.Println()
}

// exampleErrResult muestra cómo crear y usar un Result con error
func exampleErrResult() {
	fmt.Println("2. Result con error:")
	err := kerrors.NewKError("No se pudo procesar la solicitud", 500, map[string]any{
		"request_id": "req-123",
		"timestamp":  "2024-01-01T00:00:00Z",
	})
	r := result.Err[int](err)
	fmt.Printf("   Result: %+v\n", r)
	fmt.Printf("   IsOk: %v\n", r.IsOk())
	fmt.Printf("   Value: %d (valor cero)\n", r.Value())
	fmt.Printf("   Error: %v\n", r.Error())
	if r.Error() != nil {
		fmt.Printf("   Error Message: %s\n", r.Error().(*kerrors.KError).Message)
		fmt.Printf("   Error Code: %d\n", r.Error().(*kerrors.KError).Code)
	}
	fmt.Println()
}

// exampleIsOk muestra cómo verificar si un Result es exitoso antes de usar el valor
func exampleIsOk() {
	fmt.Println("3. Verificación de Result:")
	results := []result.Result[int]{
		result.Ok(100),
		result.Err[int](kerrors.NewKError("Error de validación", 400, nil)),
		result.Ok(200),
	}

	for i, r := range results {
		fmt.Printf("   Result %d:\n", i+1)
		if r.IsOk() {
			fmt.Printf("     ✓ Éxito: %d\n", r.Value())
		} else {
			fmt.Printf("     ✗ Error: %s\n", r.Error().(*kerrors.KError).Message)
		}
	}
	fmt.Println()
}

// exampleWithDifferentTypes muestra cómo usar Result con diferentes tipos
func exampleWithDifferentTypes() {
	fmt.Println("4. Result con diferentes tipos:")

	// Result con string
	strResult := result.Ok("Hola, mundo!")
	fmt.Printf("   String Result: %s\n", strResult.Value())

	// Result con slice
	sliceResult := result.Ok([]int{1, 2, 3, 4, 5})
	fmt.Printf("   Slice Result: %v\n", sliceResult.Value())

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
	fmt.Printf("   Struct Result: %+v\n", userResult.Value())
	fmt.Println()
}

// exampleWithStructuredErrors muestra cómo usar Result con errores estructurados
func exampleWithStructuredErrors() {
	fmt.Println("5. Result con errores estructurados:")

	// Simular una función que puede fallar
	divide := func(a, b int) result.Result[float64] {
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

	// Casos de prueba
	cases := []struct {
		a, b int
		desc string
	}{
		{10, 2, "División válida"},
		{10, 0, "División por cero"},
		{15, 3, "Otra división válida"},
	}

	for _, c := range cases {
		fmt.Printf("   %s (%d / %d):\n", c.desc, c.a, c.b)
		r := divide(c.a, c.b)
		if r.IsOk() {
			fmt.Printf("     ✓ Resultado: %.2f\n", r.Value())
		} else {
			kerr := r.Error().(*kerrors.KError)
			fmt.Printf("     ✗ Error: %s (Código: %d)\n", kerr.Message, kerr.Code)
			if kerr.Metadata != nil {
				fmt.Printf("       Metadata: %v\n", kerr.Metadata)
			}
		}
	}
	fmt.Println()
}
