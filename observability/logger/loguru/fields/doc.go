// Package fields provides ready-made implementations of the loguru Fields
// interface for common contexts: HTTPFileds for HTTP requests (method, path,
// trace_id, caller_id, client_ip, metadata) and WSFields for WebSocket
// connections (server_id, client_id, trace_id, path, metadata).
//
// Field values are both formatted into the console log line and attached to the
// underlying logrus entry, so hooks (for example the Loki hook) can turn them
// into labels.
//
//	import "github.com/reitmas32/rkit/observability/logger/loguru/fields"
package fields
