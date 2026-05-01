package main

import (
	"github.com/reitmas32/rkit/observability/logger/loguru"
	"github.com/reitmas32/rkit/observability/logger/loguru/fields"
	"github.com/reitmas32/rkit/observability/logger/loguru/hooks"
	"github.com/sirupsen/logrus"
)

// Ejemplo avanzado con diferentes configuraciones y opciones
func main() {
	// 1. Logger básico con campos HTTP y colores
	httpFields := &fields.HTTPFileds{}
	httpFields.UpdateOne("method", "GET")
	httpFields.UpdateOne("path", "/api/v1/advanced")
	httpFields.UpdateOne("trace_id", "advanced-trace-001")
	httpFields.UpdateOne("caller_id", "advanced-caller-001")
	httpFields.UpdateOne("client_ip", "10.0.0.1")

	config := loguru.Config{
		Level:     logrus.DebugLevel,
		Colorable: true,
	}
	logger := loguru.NewLoggerWithConfig(httpFields, config)
	logger.Debug("Mensaje de debug con colores")
	logger.Info("Mensaje de info con colores")
	logger.Warn("Mensaje de warn con colores")
	logger.Error("Mensaje de error con colores")

	// 2. Logger con hook de Loki (configuración avanzada)
	lokiHook := hooks.NewLokiBufferedHook(
		"http://localhost:3100/loki/api/v1/push",
		5, // Batch size más grande para producción
		map[string]string{
			"app":         "base-kit",
			"environment": "production",
			"service":     "api",
			"version":     "1.0.0",
		},
	)

	logger.AddHook(lokiHook)

	// Simular múltiples logs
	for i := 1; i <= 7; i++ {
		logger.Info("Log número %d - se enviará cuando llegue al batch size", i)
	}

	// Flush manual para enviar logs restantes
	lokiHook.Flush()

	// 3. Cambiar nivel de log dinámicamente
	logger.SetConfig(loguru.Config{
		Level:     logrus.WarnLevel, // Solo warnings y errores
		Colorable: true,
	})
	logger.Debug("Este debug NO se mostrará")
	logger.Info("Este info NO se mostrará")
	logger.Warn("Este warn SÍ se mostrará")
	logger.Error("Este error SÍ se mostrará")

	// 4. Logger con campos WebSocket y diferentes configuraciones
	wsFields := fields.WSFields{
		ServerID: "advanced-server-001",
		ClientID: "advanced-client-001",
		TraceID:  "advanced-trace-002",
		Path:     "/advanced/websocket",
		Metadata: map[string]any{
			"session_id":      "session-001",
			"connection_type": "persistent",
			"qos":             "high",
		},
	}

	logger2 := loguru.NewLoggerWithConfig(&wsFields, loguru.Config{
		Level:     logrus.InfoLevel,
		Colorable: false, // Sin colores para producción
	})

	lokiHook2 := hooks.NewLokiBufferedHook(
		"http://localhost:3100/loki/api/v1/push",
		10, // Batch size grande para WebSocket
		map[string]string{
			"app":         "base-kit",
			"environment": "production",
			"service":     "websocket",
			"version":     "1.0.0",
		},
	)

	logger2.AddHook(lokiHook2)
	logger2.Info("Conexión WebSocket establecida")
	logger2.Info("Mensaje de alta prioridad recibido")
	logger2.Warn("Advertencia: conexión inestable detectada")
	logger2.Error("Error crítico en conexión WebSocket")

	// Flush final
	lokiHook2.Flush()
}
