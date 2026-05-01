package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/reitmas32/rkit/core/customctx"
	corehttp "github.com/reitmas32/rkit/core/http"
)

// Do performs an HTTP request with the given context and request options.
// It will retry the request if retries are configured and the error/status code is retryable.
func (c *Client) Do(ctx *customctx.CustomContext, req *corehttp.Request) (*corehttp.Response, error) {
	startTime := time.Now()
	fullURL := c.buildURL(req.URL)

	// Read body once if it exists, so we can reuse it for retries
	var bodyBytes []byte
	var bodyReader io.Reader
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			duration := time.Since(startTime)
			c.logResponse(req.Method, fullURL, 0, duration, fmt.Errorf("failed to read request body: %w", err))
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		// Close original body if it's a Closer
		if closer, ok := req.Body.(io.Closer); ok {
			closer.Close()
		}
	}

	var lastErr error
	var lastResp *corehttp.Response
	maxAttempts := c.config.MaxRetries + 1 // +1 for initial attempt

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Log request start (only on first attempt or retries)
		if attempt == 0 {
			c.logRequest(req.Method, fullURL, startTime)
		}

		// Create body reader for this attempt
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		// Create http.Request with full URL
		httpReq, err := http.NewRequestWithContext(ctx, req.Method, fullURL, bodyReader)
		if err != nil {
			duration := time.Since(startTime)
			c.logResponse(req.Method, fullURL, 0, duration, fmt.Errorf("failed to create HTTP request: %w", err))
			return nil, fmt.Errorf("failed to create HTTP request: %w", err)
		}

		// Set headers
		if req.ContentType != "" {
			httpReq.Header.Set("Content-Type", req.ContentType)
		}
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}

		// Add default headers
		for k, v := range c.config.DefaultHeaders {
			if _, exists := req.Headers[k]; !exists {
				httpReq.Header.Set(k, v)
			}
		}

		// Handle timeout if specified in request
		if req.Timeout != nil {
			timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(*req.Timeout)*time.Second)
			defer cancel()
			httpReq = httpReq.WithContext(timeoutCtx)
		}

		// Perform request
		httpResp, err := c.httpClient.Do(httpReq)
		duration := time.Since(startTime)

		if err != nil {
			lastErr = err
			// Check if we should retry
			if c.shouldRetry(req.Method, err, 0) && attempt < maxAttempts-1 {
				delay := c.getRetryDelay(attempt)
				c.logRetry(req.Method, fullURL, attempt, delay, fmt.Sprintf("network error: %v", err))
				time.Sleep(delay)
				continue
			}
			// No more retries or not retryable
			c.logResponse(req.Method, fullURL, 0, duration, fmt.Errorf("HTTP request failed: %w", err))
			return nil, fmt.Errorf("HTTP request failed: %w", err)
		}

		// Convert response
		resp := &corehttp.Response{
			StatusCode:    httpResp.StatusCode,
			Status:        httpResp.Status,
			Headers:       make(map[string]string),
			Body:          httpResp.Body,
			ContentLength: httpResp.ContentLength,
		}

		// Copy headers
		for k, v := range httpResp.Header {
			if len(v) > 0 {
				resp.Headers[k] = v[0]
			}
		}

		// Check if status code is retryable
		if c.shouldRetry(req.Method, nil, httpResp.StatusCode) && attempt < maxAttempts-1 {
			// Close the response body before retrying
			httpResp.Body.Close()
			delay := c.getRetryDelay(attempt)
			c.logRetry(req.Method, fullURL, attempt, delay, fmt.Sprintf("status code: %d", httpResp.StatusCode))
			time.Sleep(delay)
			lastResp = resp
			continue
		}

		// Success or non-retryable error
		c.logResponse(req.Method, fullURL, httpResp.StatusCode, duration, nil)
		return resp, nil
	}

	// All retries exhausted
	if lastErr != nil {
		duration := time.Since(startTime)
		c.logResponse(req.Method, fullURL, 0, duration, fmt.Errorf("HTTP request failed after %d attempts: %w", maxAttempts, lastErr))
		return nil, fmt.Errorf("HTTP request failed after %d attempts: %w", maxAttempts, lastErr)
	}

	// Last response (with retryable status code but no more retries)
	if lastResp != nil {
		duration := time.Since(startTime)
		c.logResponse(req.Method, fullURL, lastResp.StatusCode, duration, nil)
		return lastResp, nil
	}

	// Should not reach here
	return nil, fmt.Errorf("unexpected error in retry logic")
}
