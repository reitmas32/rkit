package main

import (
	"github.com/reitmas32/rkit/observability/logger/loguru"
	"github.com/reitmas32/rkit/observability/logger/loguru/fields"
	"github.com/reitmas32/rkit/observability/logger/loguru/hooks"
)

// Ejemplo de uso del logger con hook de Loki
func main() {
	// Crear campos HTTP
	httpFields := &fields.HTTPFileds{}
	httpFields.UpdateOne("method", "GET")
	httpFields.UpdateOne("path", "/api/v1/users")
	httpFields.UpdateOne("trace_id", "1234567890")
	httpFields.UpdateOne("caller_id", "9876543210")
	httpFields.UpdateOne("client_ip", "192.168.1.100")

	// Crear logger
	logger := loguru.NewLogger(httpFields)

	// Crear hook de Loki con batching
	// BatchSize: 2 significa que se enviarán logs a Loki cuando haya 2 logs acumulados
	lokiHook := hooks.NewLokiBufferedHook(
		"http://localhost:3100/loki/api/v1/push", // URL del endpoint de Loki
		2,                                        // Batch size: envía cuando hay 2 logs
		map[string]string{
			"app":         "base-kit",
			"environment": "development",
			"service":     "api",
		},
	)

	// Agregar hook al logger
	logger.AddHook(lokiHook)

	// Los siguientes logs se acumularán hasta llegar al batch size
	logger.Info("Primer log - será enviado cuando llegue el segundo")
	logger.Info("Segundo log - esto disparará el envío a Loki (batch de 2)")

	// Los siguientes logs se acumularán nuevamente
	logger.Error("Error log - será enviado cuando llegue el siguiente")
	logger.Warn("Warning log - esto disparará otro envío a Loki")

	// Forzar el envío de logs restantes (si hay alguno en el buffer)
	lokiHook.Flush()

	// Ejemplo con campos WebSocket
	wsFields := fields.WSFields{
		ServerID: "server-123",
		ClientID: "client-456",
		TraceID:  "trace-789",
		Path:     "/chat",
		Metadata: map[string]any{
			"room_id": "room-001",
		},
	}

	logger2 := loguru.NewLogger(&wsFields)

	// Hook de Loki para WebSocket (mismo hook, diferentes labels base)
	lokiHook2 := hooks.NewLokiBufferedHook(
		"http://localhost:3100/loki/api/v1/push",
		3, // Batch size: 3 logs
		map[string]string{
			"app":         "base-kit",
			"environment": "development",
			"service":     "websocket",
		},
	)

	logger2.AddHook(lokiHook2)
	logger2.Info("Conexión WebSocket establecida")
	logger2.Info("Mensaje recibido")
	logger2.Info("Mensaje procesado - esto disparará el envío (batch de 3)")

	// Flush final
	lokiHook2.Flush()
}
