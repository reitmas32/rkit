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
func NewClient(config Config) *Client {
	timeout := time.Duration(config.Timeout) * time.Second

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
