// Package loguru is a logrus-based implementation of the core/logger ILogger
// contract. It produces structured console logs in the format
// "DATE_TIME | LEVEL | FILE.FUNCTION:LINE | [fields] | MESSAGE", supports
// optional ANSI colors, and exposes logrus hooks so logs can be shipped to
// external systems such as Grafana Loki.
//
// Build a logger with NewLogger(fields) or NewLoggerWithConfig(fields, config),
// where fields implements the Fields interface (see the fields subpackage for
// ready-made HTTP and WebSocket field sets). Attach hooks with AddHook; the
// hooks subpackage provides a buffered Loki hook.
//
//	import "github.com/reitmas32/rkit/observability/logger/loguru"
//
//	log := loguru.NewLogger(nil)
//	log.Info("user %s authenticated", userID)
package loguru
