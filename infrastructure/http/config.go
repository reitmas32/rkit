package http

import (
	"net/http"

	"github.com/reitmas32/rkit/core/logger"
)

// Config holds the configuration for the HTTP client.
type Config struct {
	// Timeout is the default per-request timeout in seconds.
	//   >  0  use this many seconds.
	//   == 0  use DefaultTimeoutSeconds (secure default; avoids an unbounded client).
	//   <  0  explicit opt-in to no timeout.
	Timeout int

	// MaxResponseBytes caps how many bytes are read from a response body,
	// protecting against memory exhaustion from oversized or malicious responses.
	//   >  0  reject bodies larger than this many bytes (ErrResponseTooLarge).
	//   == 0  use DefaultMaxResponseBytes.
	//   <  0  explicit opt-in to unlimited (not recommended).
	MaxResponseBytes int64

	// BaseURL is an optional base URL that will be prepended to all requests
	BaseURL string

	// DefaultHeaders are headers that will be added to all requests
	DefaultHeaders map[string]string

	// Transport is the HTTP transport to use (optional, uses http.DefaultTransport if nil)
	Transport http.RoundTripper

	// CheckRedirect specifies the policy for handling redirects
	CheckRedirect func(req *http.Request, via []*http.Request) error

	// Logger is an optional logger for HTTP request/response logging.
	// If nil, no logs will be generated.
	Logger logger.ILogger

	// DisableLogging explicitly disables logging even if Logger is set.
	// This allows you to disable logs without removing the logger from config.
	DisableLogging bool

	// MaxRetries is the maximum number of retry attempts for failed requests.
	// If 0, no retries will be performed (default).
	// Retries are only performed for network errors and retryable status codes.
	MaxRetries int

	// RetryDelay is the delay between retry attempts in milliseconds.
	// If 0, a default delay of 100ms will be used.
	RetryDelay int

	// RetryableStatusCodes is a list of HTTP status codes that should trigger a retry.
	// By default, retries are performed for:
	// - Network errors (timeouts, connection errors)
	// - Status codes: 429 (Too Many Requests), 500, 502, 503, 504
	// If nil, default retryable codes will be used.
	RetryableStatusCodes []int

	// RetryableMethods specifies which HTTP methods should be retried.
	// By default, only idempotent methods (GET, HEAD, OPTIONS, DELETE) are retried.
	// If you want to retry non-idempotent methods (POST, PUT, PATCH), add them here.
	RetryableMethods []string
}

const (
	// DefaultTimeoutSeconds is the per-request timeout applied when Config.Timeout is 0.
	DefaultTimeoutSeconds = 30
	// DefaultMaxResponseBytes is the response body cap applied when
	// Config.MaxResponseBytes is 0 (10 MiB).
	DefaultMaxResponseBytes int64 = 10 << 20
)

// DefaultConfig returns a default configuration.
func DefaultConfig() Config {
	return Config{
		Timeout:              DefaultTimeoutSeconds,
		MaxResponseBytes:     DefaultMaxResponseBytes,
		BaseURL:              "",
		DefaultHeaders:       make(map[string]string),
		Transport:            nil, // Uses http.DefaultTransport
		CheckRedirect:        nil, // Uses http.DefaultClient's redirect policy
		Logger:               nil, // No logger by default
		DisableLogging:       false,
		MaxRetries:           0,                                            // No retries by default
		RetryDelay:           100,                                          // 100ms default delay
		RetryableStatusCodes: []int{429, 500, 502, 503, 504},               // Default retryable codes
		RetryableMethods:     []string{"GET", "HEAD", "OPTIONS", "DELETE"}, // Default idempotent methods
	}
}
