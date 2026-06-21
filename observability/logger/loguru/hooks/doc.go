// Package hooks provides logrus hooks for the loguru logger, most notably
// LokiBufferedHook which ships logs to Grafana Loki's push API in batches.
//
// The hook buffers each entry and POSTs a batch once BatchSize is reached
// (default 50); call Flush to send a partial batch, for example at the end of an
// HTTP request or on shutdown. If a POST fails the buffer is retained so no logs
// are lost. Base labels plus any structured fields become Loki stream labels.
//
//	import "github.com/reitmas32/rkit/observability/logger/loguru/hooks"
//
//	hook := hooks.NewLokiBufferedHook("http://localhost:3100/loki/api/v1/push", 50,
//	    map[string]string{"app": "my-service", "environment": "production"})
//	logger.AddHook(hook)
//	defer hook.Flush()
package hooks
