// Package logger defines the ILogger contract used across rkit together with a
// minimal SimpleLogger implementation and log-level helpers (Level, ParseLevel,
// IsLoggable). Depending on ILogger lets components stay decoupled from any
// concrete logging backend.
//
// For a structured, production-grade logger with caller info and a Grafana Loki
// hook, see github.com/reitmas32/rkit/observability/logger/loguru.
//
//	import "github.com/reitmas32/rkit/core/logger"
//
//	log := logger.NewSimpleLogger("info")
//	log.Info("server started on port %d", 8080)
package logger
