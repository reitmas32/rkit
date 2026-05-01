package http

import "time"

// shouldLog returns true if logging should be performed.
// Logging is enabled if:
// - Logger is not nil
// - DisableLogging is false
func (c *Client) shouldLog() bool {
	return c.config.Logger != nil && !c.config.DisableLogging
}

// logRequest logs the start of an HTTP request.
func (c *Client) logRequest(method, url string, startTime time.Time) {
	if !c.shouldLog() {
		return
	}
	c.config.Logger.Debug("HTTP request started: method=%s url=%s timestamp=%v\n", method, url, startTime)
}

// logResponse logs the completion of an HTTP request.
func (c *Client) logResponse(method, url string, statusCode int, duration time.Duration, err error) {
	if !c.shouldLog() {
		return
	}

	if err != nil {
		c.config.Logger.Error("HTTP request failed: method=%s url=%s status_code=%d duration=%v error=%v\n", method, url, statusCode, duration, err)
		return
	}

	// Log level based on status code
	if statusCode >= 500 {
		c.config.Logger.Error("HTTP request completed with server error: method=%s url=%s status_code=%d duration=%v\n", method, url, statusCode, duration)
	} else if statusCode >= 400 {
		c.config.Logger.Warn("HTTP request completed with client error: method=%s url=%s status_code=%d duration=%v\n", method, url, statusCode, duration)
	} else {
		c.config.Logger.Info("HTTP request completed successfully: method=%s url=%s status_code=%d duration=%v\n", method, url, statusCode, duration)
	}
}
