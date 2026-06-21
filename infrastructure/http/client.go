package http

import (
	"net/http"
	"time"
)

// Client is an HTTP client implementation using the standard library.
type Client struct {
	httpClient *http.Client
	config     Config
}

// NewClient creates a new HTTP client with the given configuration.
//
// A zero Config.Timeout is replaced with DefaultTimeoutSeconds so a client built
// from a bare Config{} is never left without a timeout. A negative Timeout is an
// explicit opt-in to no timeout. Likewise, a zero MaxResponseBytes is replaced
// with DefaultMaxResponseBytes.
func NewClient(config Config) *Client {
	var timeout time.Duration
	switch {
	case config.Timeout > 0:
		timeout = time.Duration(config.Timeout) * time.Second
	case config.Timeout == 0:
		timeout = DefaultTimeoutSeconds * time.Second
	default: // < 0: explicit opt-in to no timeout
		timeout = 0
	}

	if config.MaxResponseBytes == 0 {
		config.MaxResponseBytes = DefaultMaxResponseBytes
	}

	client := &http.Client{
		Timeout:       timeout,
		Transport:     config.Transport,
		CheckRedirect: config.CheckRedirect,
	}

	return &Client{
		httpClient: client,
		config:     config,
	}
}
