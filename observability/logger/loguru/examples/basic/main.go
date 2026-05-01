package main

import (
	"github.com/reitmas32/rkit/observability/logger/loguru"
	"github.com/sirupsen/logrus"
)

// Ejemplo básico de uso del logger sin campos personalizados
func main() {
	// Crear logger básico (sin fields)
	logger := loguru.NewLogger(nil)

	logger.Debug("Este es un mensaje de debug")
	logger.Info("Este es un mensaje de información")
	logger.Warn("Este es un mensaje de advertencia")
	logger.Error("Este es un mensaje de error")

	// Logger con configuración personalizada (nivel y colores)
	config := loguru.Config{
		Level:     logrus.InfoLevel,
		Colorable: true,
	}
	logger2 := loguru.NewLoggerWithConfig(nil, config)
	logger2.Debug("Este debug no se mostrará (nivel Info)")
	logger2.Info("Este info se mostrará con colores")
	logger2.Warn("Este warn se mostrará con colores")
	logger2.Error("Este error se mostrará con colores")

	// Cambiar configuración en tiempo de ejecución
	logger3 := loguru.NewLogger(nil)
	logger3.Info("Sin colores")
	logger3.SetConfig(loguru.Config{
		Level:     logrus.DebugLevel,
		Colorable: true,
	})
	logger3.Debug("Ahora con colores y nivel Debug")
	logger3.Info("Este info también se mostrará con colores")
}
