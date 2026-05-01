package http

import (
	"net/url"
	"strings"
)

// buildURL builds the full URL by prepending the base URL if configured.
func (c *Client) buildURL(urlPath string) string {
	if c.config.BaseURL == "" {
		return urlPath
	}

	// Parse base URL
	baseURL, err := url.Parse(c.config.BaseURL)
	if err != nil {
		// If base URL is invalid, just concatenate
		return strings.TrimSuffix(c.config.BaseURL, "/") + "/" + strings.TrimPrefix(urlPath, "/")
	}

	// Parse path URL
	pathURL, err := url.Parse(urlPath)
	if err != nil {
		// If path is invalid, just concatenate
		return strings.TrimSuffix(c.config.BaseURL, "/") + "/" + strings.TrimPrefix(urlPath, "/")
	}

	// Resolve reference
	resolved := baseURL.ResolveReference(pathURL)
	return resolved.String()
}
