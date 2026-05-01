package main

import (
	"errors"
	"fmt"

	"github.com/reitmas32/rkit/core/kerrors"
)

func main() {
	fmt.Printf("=== Ejemplos del paquete kerrors ===\n\n")

	// Ejemplo 1: Crear un error básico
	exampleBasicError()

	// Ejemplo 2: Crear un error con metadata
	exampleErrorWithMetadata()

	// Ejemplo 3: Encadenar errores
	exampleErrorChaining()

	// Ejemplo 4: Usar errores con errors.Is y errors.Unwrap
	exampleErrorInspection()

	// Ejemplo 5: Caso de uso práctico
	examplePracticalUseCase()
}

// exampleBasicError muestra cómo crear un error básico
func exampleBasicError() {
	fmt.Println("1. Error básico:")
	err := kerrors.NewKError("Recurso no encontrado", 404, nil)
	fmt.Printf("   Mensaje: %s\n", err.Message)
	fmt.Printf("   Código: %d\n", err.Code)
	fmt.Printf("   Error(): %s\n", err.Error())
	fmt.Println()
}

// exampleErrorWithMetadata muestra cómo crear un error con información adicional
func exampleErrorWithMetadata() {
	fmt.Println("2. Error con metadata:")
	metadata := map[string]any{
		"request_id":  "req-abc-123",
		"user_id":     42,
		"timestamp":   "2024-01-01T12:00:00Z",
		"endpoint":    "/api/users",
		"http_method": "GET",
		"ip_address":  "192.168.1.1",
	}

	err := kerrors.NewKError("Error de validación", 400, metadata)
	fmt.Printf("   Mensaje: %s\n", err.Message)
	fmt.Printf("   Código: %d\n", err.Code)
	fmt.Printf("   Metadata:\n")
	for key, value := range err.Metadata {
		fmt.Printf("     %s: %v\n", key, value)
	}
	fmt.Println()
}

// exampleErrorChaining muestra cómo encadenar errores
func exampleErrorChaining() {
	fmt.Println("3. Encadenamiento de errores:")

	// Error original (raíz)
	rootErr := errors.New("error de base de datos: conexión perdida")

	// Error intermedio
	intermediateErr := kerrors.NewKErrorWithCause(
		"Error al obtener usuario",
		500,
		map[string]any{
			"operation": "get_user",
			"user_id":   123,
		},
		rootErr,
	)

	// Error de nivel superior
	topErr := kerrors.NewKErrorWithCause(
		"Error al procesar solicitud",
		500,
		map[string]any{
			"request_id": "req-xyz",
			"service":    "user_service",
		},
		intermediateErr,
	)

	fmt.Printf("   Error superior: %s\n", topErr.Message)
	fmt.Printf("   Desenvolviendo la cadena:\n")

	current := topErr
	level := 1
	for current != nil {
		fmt.Printf("     Nivel %d: %s (Código: %d)\n", level, current.Message, current.Code)
		if len(current.Metadata) > 0 {
			fmt.Printf("       Metadata: %v\n", current.Metadata)
		}

		unwrapped := current.Unwrap()
		if unwrapped == nil {
			break
		}

		// Verificar si el error desenvuelto es un KError
		if kerr, ok := unwrapped.(*kerrors.KError); ok {
			current = kerr
		} else {
			fmt.Printf("     Nivel %d (raíz): %s\n", level+1, unwrapped.Error())
			break
		}
		level++
	}
	fmt.Println()
}

// exampleErrorInspection muestra cómo usar errors.Is y errors.Unwrap
func exampleErrorInspection() {
	fmt.Println("4. Inspección de errores con errors.Is y errors.Unwrap:")

	rootErr := errors.New("error de red")
	wrappedErr := kerrors.NewKErrorWithCause(
		"Error al conectar con API externa",
		503,
		map[string]any{"service": "external_api"},
		rootErr,
	)

	// Usar errors.Is para verificar si el error raíz está en la cadena
	if errors.Is(wrappedErr, rootErr) {
		fmt.Printf("   ✓ errors.Is encontró el error raíz en la cadena\n")
	}

	// Usar errors.Unwrap para obtener el error subyacente
	unwrapped := errors.Unwrap(wrappedErr)
	if unwrapped != nil {
		fmt.Printf("   ✓ errors.Unwrap retornó: %s\n", unwrapped.Error())
	}

	// Verificar que el error implementa la interfaz error
	var err error = wrappedErr
	fmt.Printf("   ✓ KError implementa la interfaz error: %s\n", err.Error())
	fmt.Println()
}

// examplePracticalUseCase muestra un caso de uso práctico
func examplePracticalUseCase() {
	fmt.Println("5. Caso de uso práctico - Validación de usuario:")

	// Simular validación de usuario
	validateUser := func(userID int, email string) *kerrors.KError {
		if userID <= 0 {
			return kerrors.NewKError(
				"ID de usuario inválido",
				400,
				map[string]any{
					"field":    "user_id",
					"value":    userID,
					"expected": "número positivo",
				},
			)
		}

		if email == "" {
			return kerrors.NewKError(
				"Email es requerido",
				400,
				map[string]any{
					"field": "email",
					"value": email,
				},
			)
		}

		// Simular error de base de datos
		dbErr := errors.New("database: connection timeout")
		return kerrors.NewKErrorWithCause(
			"Error al guardar usuario",
			500,
			map[string]any{
				"operation": "save_user",
				"user_id":   userID,
				"email":     email,
			},
			dbErr,
		)
	}

	// Casos de prueba
	cases := []struct {
		userID int
		email  string
		desc   string
	}{
		{0, "test@example.com", "ID inválido"},
		{1, "", "Email vacío"},
		{42, "user@example.com", "Error de base de datos"},
	}

	for _, c := range cases {
		fmt.Printf("   Caso: %s\n", c.desc)
		err := validateUser(c.userID, c.email)
		if err != nil {
			fmt.Printf("     ✗ Error: %s\n", err.Message)
			fmt.Printf("       Código: %d\n", err.Code)
			if err.Metadata != nil {
				fmt.Printf("       Metadata: %v\n", err.Metadata)
			}
			if err.Cause != nil {
				fmt.Printf("       Causa: %s\n", err.Cause.Error())
			}
		}
		fmt.Println()
	}
}
