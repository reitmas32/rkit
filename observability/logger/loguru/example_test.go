package loguru_test

import (
	"github.com/reitmas32/rkit/observability/logger/loguru"
	"github.com/reitmas32/rkit/observability/logger/loguru/fields"
	"github.com/reitmas32/rkit/observability/logger/loguru/hooks"
)

// ExampleNewLogger logs a message enriched with structured HTTP fields. The
// fields are rendered in the console line and attached to the entry so hooks can
// turn them into labels.
func ExampleNewLogger() {
	f := &fields.HTTPFileds{}
	f.UpdateOne("method", "GET")
	f.UpdateOne("path", "/api/v1/users")
	f.UpdateOne("trace_id", "req-123")

	log := loguru.NewLogger(f)
	log.Info("request received")
}

// ExampleLogger_AddHook ships logs to Grafana Loki in batches of 50 and flushes
// any partial batch at the end of the unit of work.
func ExampleLogger_AddHook() {
	log := loguru.NewLogger(nil)

	hook := hooks.NewLokiBufferedHook(
		"http://localhost:3100/loki/api/v1/push",
		50,
		map[string]string{"app": "my-service", "environment": "production"},
	)
	log.AddHook(hook)
	defer hook.Flush()

	log.Info("user authenticated")
}
