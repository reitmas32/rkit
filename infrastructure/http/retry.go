package http

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"
)

// isRetryableError checks if an error is retryable.
// Network errors (timeouts, connection errors) are retryable.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		// Network errors (timeouts, connection errors) are retryable
		return true
	}

	// Context deadline exceeded is retryable
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	return false
}

// isRetryableStatusCode checks if a status code should trigger a retry.
func (c *Client) isRetryableStatusCode(statusCode int) bool {
	if c.config.RetryableStatusCodes == nil {
		// Use default retryable codes
		defaultCodes := []int{429, 500, 502, 503, 504}
		for _, code := range defaultCodes {
			if statusCode == code {
				return true
			}
		}
		return false
	}

	for _, code := range c.config.RetryableStatusCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}

// isRetryableMethod checks if an HTTP method is retryable.
func (c *Client) isRetryableMethod(method string) bool {
	// Default idempotent methods
	defaultMethods := []string{"GET", "HEAD", "OPTIONS", "DELETE"}
	
	// Check if method is in default retryable methods
	for _, m := range defaultMethods {
		if strings.EqualFold(method, m) {
			return true
		}
	}

	// Check if method is in custom retryable methods
	if c.config.RetryableMethods != nil {
		for _, m := range c.config.RetryableMethods {
			if strings.EqualFold(method, m) {
				return true
			}
		}
	}

	return false
}

// shouldRetry determines if a request should be retried based on error and status code.
func (c *Client) shouldRetry(method string, err error, statusCode int) bool {
	// No retries configured
	if c.config.MaxRetries <= 0 {
		return false
	}

	// Check if method is retryable
	if !c.isRetryableMethod(method) {
		return false
	}

	// Retry on network errors
	if err != nil && isRetryableError(err) {
		return true
	}

	// Retry on retryable status codes
	if statusCode > 0 && c.isRetryableStatusCode(statusCode) {
		return true
	}

	return false
}

// getRetryDelay returns the delay to wait before retrying.
func (c *Client) getRetryDelay(attempt int) time.Duration {
	delay := c.config.RetryDelay
	if delay <= 0 {
		delay = 100 // Default 100ms
	}

	// Simple exponential backoff: delay * 2^attempt
	// Cap at 5 seconds
	exponentialDelay := time.Duration(delay) * time.Millisecond * time.Duration(1<<uint(attempt))
	if exponentialDelay > 5*time.Second {
		exponentialDelay = 5 * time.Second
	}

	return exponentialDelay
}

// logRetry logs a retry attempt.
func (c *Client) logRetry(method, url string, attempt int, delay time.Duration, reason string) {
	if !c.shouldLog() {
		return
	}
	c.config.Logger.Warn("HTTP request retry: method=%s url=%s attempt=%d/%d delay=%v reason=%s\n",
		method, url, attempt+1, c.config.MaxRetries+1, delay, reason)
}
