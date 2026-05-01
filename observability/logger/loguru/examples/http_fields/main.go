package main

import (
	"github.com/reitmas32/rkit/observability/logger/loguru"
	"github.com/reitmas32/rkit/observability/logger/loguru/fields"
	"github.com/sirupsen/logrus"
)

// Ejemplo de uso del logger con campos HTTP
func main() {
	// Crear campos HTTP
	httpFields := &fields.HTTPFileds{}
	httpFields.UpdateOne("method", "GET")
	httpFields.UpdateOne("path", "/api/v1/users")
	httpFields.UpdateOne("trace_id", "1234567890")
	httpFields.UpdateOne("caller_id", "9876543210")
	httpFields.UpdateOne("client_ip", "192.168.1.100")

	metadata := map[string]any{
		"user_id":   "user-123",
		"user_name": "John Doe",
		"role":      "admin",
	}
	httpFields.UpdateOne("metadata", metadata)

	// Logger con campos HTTP
	logger := loguru.NewLogger(httpFields)
	logger.Info("Usuario autenticado correctamente")
	logger.Info("Solicitud procesada exitosamente")

	// Logger con configuración y campos HTTP
	config := loguru.Config{
		Level:     logrus.InfoLevel,
		Colorable: true,
	}
	logger2 := loguru.NewLoggerWithConfig(httpFields, config)
	logger2.Info("Inicio de solicitud HTTP")
	logger2.Debug("Debug: detalles de la solicitud")
	logger2.Warn("Advertencia: tasa de solicitudes alta")
	logger2.Error("Error al procesar solicitud")

	// Actualizar campos dinámicamente
	httpFields2 := &fields.HTTPFileds{}
	httpFields2.UpdateOne("method", "POST")
	httpFields2.UpdateOne("path", "/api/v1/users")
	httpFields2.UpdateOne("trace_id", "new-trace-id")
	httpFields2.UpdateOne("caller_id", "new-caller-id")
	httpFields2.UpdateOne("client_ip", "10.0.0.1")

	logger3 := loguru.NewLogger(httpFields2)
	logger3.Info("Nueva solicitud POST recibida")
}
