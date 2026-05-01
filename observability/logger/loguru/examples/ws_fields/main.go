package main

import (
	"github.com/reitmas32/rkit/observability/logger/loguru"
	"github.com/reitmas32/rkit/observability/logger/loguru/fields"
)

// Ejemplo de uso del logger con campos WebSocket
func main() {
	// Crear campos WebSocket
	wsFields := fields.WSFields{
		ServerID: "e8a18817-fde9-4594-a88f-a72379a8a3c2",
		ClientID: "client-123",
		TraceID:  "cfad8504-4d47-428d-90e2-b7d17213f831",
		Path:     "/chat",
		Metadata: map[string]any{
			"operator_id":     "74906178-2564-4150-9008-61466217737b",
			"organization_id": "e8a18817-fde9-4594-a88f-a72379a8a3c2",
			"room_id":         "room-456",
		},
	}

	// Logger con campos WebSocket
	logger := loguru.NewLogger(&wsFields)
	logger.Info("Conexión WebSocket establecida")
	logger.Info("Mensaje recibido del cliente")
	logger.Warn("Cliente con latencia alta detectado")
	logger.Error("Error al procesar mensaje WebSocket")

	// Actualizar campos dinámicamente
	wsFields2 := fields.WSFields{}
	wsFields2.UpdateOne("server_id", "new-server-id")
	wsFields2.UpdateOne("client_id", "new-client-id")
	wsFields2.UpdateOne("trace_id", "new-trace-id")
	wsFields2.UpdateOne("path", "/notifications")

	metadata2 := map[string]any{
		"session_id": "session-789",
		"user_type":  "premium",
	}
	wsFields2.UpdateOne("metadata", metadata2)

	logger2 := loguru.NewLogger(&wsFields2)
	logger2.Info("Nueva conexión WebSocket")
	logger2.Info("Mensaje broadcast enviado")
}
