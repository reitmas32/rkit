package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
)

// ErrorInterface defines the interface for custom errors used in responses
type ErrorInterface interface {
	GetCode() int
	GetScope() string
	Error() string
}

// KErrorAdapter adapts kerrors.KError to ErrorInterface
type KErrorAdapter struct {
	*kerrors.KError
	scope string
}

// NewKErrorAdapter creates an adapter for KError
func NewKErrorAdapter(err *kerrors.KError, scope string) *KErrorAdapter {
	return &KErrorAdapter{
		KError: err,
		scope:  scope,
	}
}

// GetCode returns the error code
func (a *KErrorAdapter) GetCode() int {
	return a.KError.Code
}

// GetScope returns the error scope
func (a *KErrorAdapter) GetScope() string {
	return a.scope
}

// HTTPStatus returns the HTTP status code for an error
func HTTPStatus(err ErrorInterface) int {
	code := err.GetCode()
	// Map error codes to HTTP status codes
	if code >= 400 && code < 600 {
		return code
	}
	// Default to 500 for unknown error codes
	return http.StatusInternalServerError
}

// PublicMessageOf returns the public message of an error
func PublicMessageOf(err ErrorInterface) string {
	return err.Error()
}

// Response is a generic HTTP response wrapper.
type Response[T any] struct {
	// Success indicates if the request was successful
	Success bool `json:"success"`

	// StatusCode is the HTTP status code
	StatusCode int `json:"status_code,omitempty"`

	// Data holds the response payload for single items
	Data T `json:"data,omitempty"`

	// Results holds the response payload for collections
	Results []T `json:"results,omitempty"`

	// Error holds error information if Success is false
	Error *ErrorInfo `json:"error,omitempty"`

	// Alert holds user-facing notification
	Alert *Alert `json:"alert,omitempty"`

	// TraceID for request tracing
	TraceID string `json:"trace_id,omitempty"`

	// Meta holds pagination or other metadata
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// PaginationMeta represents pagination metadata in API responses.
type PaginationMeta struct {
	// Page is the current page number (0-indexed)
	Page int `json:"page"`

	// Size is the number of items per page
	Size int `json:"size"`

	// Offset is the number of items to skip
	Offset int `json:"offset"`

	// Total is the total number of items across all pages
	Total int64 `json:"total"`

	// TotalPages is the total number of pages
	TotalPages int `json:"totalPages"`

	// This is the URL for the current page
	This string `json:"this"`

	// Next is the URL for the next page (nil if there is no next page)
	Next *string `json:"next,omitempty"`

	// Prev is the URL for the previous page (nil if there is no previous page)
	Prev *string `json:"prev,omitempty"`
}

// ToMap converts PaginationMeta to a map for inclusion in Response.Meta
func (p PaginationMeta) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"page":       p.Page,
		"size":       p.Size,
		"offset":     p.Offset,
		"total":      p.Total,
		"totalPages": p.TotalPages,
		"this":       p.This,
	}
	if p.Next != nil {
		result["next"] = *p.Next
	}
	if p.Prev != nil {
		result["prev"] = *p.Prev
	}
	return result
}

// BuildPaginationURL constructs a pagination URL with the given parameters.
// It takes a base URL and adds page, size, and offset as query parameters.
func BuildPaginationURL(baseURL string, page, size, offset int) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}
	q := u.Query()
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("size", fmt.Sprintf("%d", size))
	q.Set("offset", fmt.Sprintf("%d", offset))
	u.RawQuery = q.Encode()
	return u.String()
}

// CalculateTotalPages calculates the total number of pages based on total items and page size.
func CalculateTotalPages(total int64, size int) int {
	if size <= 0 {
		return 1
	}
	totalPages := int(math.Ceil(float64(total) / float64(size)))
	if totalPages == 0 {
		return 1
	}
	return totalPages
}

// NewPaginationMeta creates a new PaginationMeta with all URLs calculated.
// It takes the current page, size, offset, total items, and base URL to construct
// the pagination metadata including next, prev, and this URLs.
func NewPaginationMeta(page, size, offset int, total int64, baseURL string) PaginationMeta {
	totalPages := CalculateTotalPages(total, size)

	// Build current page URL
	thisURL := BuildPaginationURL(baseURL, page, size, offset)

	// Build next page URL if there is one
	var nextURL *string
	if page < totalPages-1 {
		nextPage := page + 1
		nextOffset := nextPage * size
		next := BuildPaginationURL(baseURL, nextPage, size, nextOffset)
		nextURL = &next
	}

	// Build previous page URL if there is one
	var prevURL *string
	if page > 0 {
		prevPage := page - 1
		prevOffset := prevPage * size
		prev := BuildPaginationURL(baseURL, prevPage, size, prevOffset)
		prevURL = &prev
	}

	return PaginationMeta{
		Page:       page,
		Size:       size,
		Offset:     offset,
		Total:      total,
		TotalPages: totalPages,
		This:       thisURL,
		Next:       nextURL,
		Prev:       prevURL,
	}
}

// WithPagination adds pagination metadata to a Response's Meta field.
// It creates a new meta map if one doesn't exist, or adds to the existing one.
func WithPagination[T any](resp *Response[T], paginationMeta PaginationMeta) {
	if resp.Meta == nil {
		resp.Meta = make(map[string]interface{})
	}
	resp.Meta["pagination"] = paginationMeta.ToMap()
}

// ErrorInfo represents error information in the response.
type ErrorInfo struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Scope   string                 `json:"scope,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

func (e ErrorInfo) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
		"scope":   e.Scope,
	}
	if e.Meta != nil {
		m["meta"] = e.Meta
	}
	return m
}

func (e ErrorInfo) UnWrap() error {
	return errors.New(e.Message)
}

func NewErrorInfo(code int, message, scope string) *ErrorInfo {
	return &ErrorInfo{
		Code:    code,
		Message: message,
		Scope:   scope,
	}
}

func (e ErrorInfo) WithMeta(meta map[string]any) *ErrorInfo {
	e.Meta = meta
	return &e
}

func (e ErrorInfo) ToKError() *kerrors.KError {
	m := make(map[string]any)
	for k, v := range e.Meta {
		m[k] = v
	}
	if e.Scope != "" {
		m["scope"] = e.Scope
	}
	return kerrors.NewKError(e.Message, e.Code, m)
}

func (e *ErrorInfo) Error() string {
	return e.Message
}

func (e *ErrorInfo) GetCode() int {
	return e.Code
}

func (e *ErrorInfo) GetScope() string {
	return e.Scope
}

// Alert represents a user-facing notification.
type Alert struct {
	Title   string `json:"title,omitempty"`
	Message string `json:"message"`
	Icon    string `json:"icon,omitempty"`
	Type    string `json:"type,omitempty"` // success, error, warning, info
}

// Ok creates a successful response with data.
func Ok[T any](data T) Response[T] {
	return Response[T]{
		Success:    true,
		StatusCode: http.StatusOK,
		Data:       data,
	}
}

// OkWithStatus creates a successful response with custom status.
func OkWithStatus[T any](data T, status int) Response[T] {
	return Response[T]{
		Success:    true,
		StatusCode: status,
		Data:       data,
	}
}

// Created creates a 201 response for newly created resources.
func Created[T any](data T) Response[T] {
	return Response[T]{
		Success:    true,
		StatusCode: http.StatusCreated,
		Data:       data,
	}
}

// List creates a successful response with a collection.
func List[T any](results []T) Response[T] {
	return Response[T]{
		Success:    true,
		StatusCode: http.StatusOK,
		Results:    results,
	}
}

// ListWithMeta creates a list response with metadata (e.g., pagination).
func ListWithMeta[T any](results []T, meta map[string]interface{}) Response[T] {
	return Response[T]{
		Success:    true,
		StatusCode: http.StatusOK,
		Results:    results,
		Meta:       meta,
	}
}

// Fail creates an error response from an error.
func Fail[T any](err ErrorInterface) Response[T] {
	status := HTTPStatus(err)

	return Response[T]{
		Success:    false,
		StatusCode: status,
		Error: &ErrorInfo{
			Code:    err.GetCode(),
			Message: PublicMessageOf(err),
			Scope:   err.GetScope(),
		},
	}
}

// FailWithStatus creates an error response with custom status.
func FailWithStatus[T any](err ErrorInterface, status int) Response[T] {
	return Response[T]{
		Success:    false,
		StatusCode: status,
		Error: &ErrorInfo{
			Code:    err.GetCode(),
			Message: PublicMessageOf(err),
			Scope:   err.GetScope(),
		},
	}
}

// FailFromKError creates an error response from a KError.
func FailFromKError[T any](err *kerrors.KError, scope string) Response[T] {
	adapter := NewKErrorAdapter(err, scope)
	return Fail[T](adapter)
}

// WithTraceID adds a trace ID to the response.
func (r Response[T]) WithTraceID(traceID string) Response[T] {
	r.TraceID = traceID
	return r
}

// WithAlert adds an alert to the response.
func (r Response[T]) WithAlert(alert Alert) Response[T] {
	r.Alert = &alert
	return r
}

// WithMeta adds metadata to the response.
func (r Response[T]) WithMeta(meta map[string]interface{}) Response[T] {
	r.Meta = meta
	return r
}

// removeEmptyFields recursively removes fields that are nil, empty strings, or empty collections
func removeEmptyFields(m map[string]interface{}) {
	for key, value := range m {
		if isEmpty(value) {
			delete(m, key)
			continue
		}

		// Recursively clean nested maps
		if nestedMap, ok := value.(map[string]interface{}); ok {
			removeEmptyFields(nestedMap)
			// If the nested map becomes empty after cleaning, remove it
			if len(nestedMap) == 0 {
				delete(m, key)
				continue
			}
			// Check if all remaining values in the nested map are empty
			allEmpty := true
			for _, v := range nestedMap {
				if !isEmpty(v) {
					allEmpty = false
					break
				}
			}
			if allEmpty {
				delete(m, key)
			}
		}

		// Recursively clean slices of maps
		if slice, ok := value.([]interface{}); ok {
			for i := range slice {
				if nestedMap, ok := slice[i].(map[string]interface{}); ok {
					removeEmptyFields(nestedMap)
				}
			}
		}
	}
}

// isEmpty checks if a value should be considered empty and removed
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	case map[interface{}]interface{}:
		return len(v) == 0
	}

	return false
}

// ToMap converts the response to a map.
// It automatically removes fields that are nil, empty strings, or empty collections.
// If there's an error or success is false, the data field is removed.
func (r Response[T]) ToMap() map[string]interface{} {
	data, err := json.Marshal(r)
	if err != nil {
		return map[string]interface{}{"error": "marshal failed"}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return map[string]interface{}{"error": "unmarshal failed"}
	}

	// If there's an error or success is false, remove the data field immediately
	// This ensures data is not included in error responses
	if r.Error != nil || !r.Success {
		delete(result, "data")
	}

	// Remove empty fields (nil, empty strings, empty slices, empty maps)
	removeEmptyFields(result)

	// Final check: if there's an error or success is false, ensure data is removed
	// (in case it was re-added or not properly removed)
	if r.Error != nil || !r.Success {
		delete(result, "data")
	}

	return result
}

func (r Response[T]) ToMapWithCustomContext(ctx *customctx.CustomContext) map[string]interface{} {

	if ctx == nil {
		return r.ToMap()
	}

	// Try to get trace ID from CustomContext values first
	if traceIDValue := ctx.GetValue("trace-id"); traceIDValue != nil {
		if traceIDStr, ok := traceIDValue.(string); ok {
			r.TraceID = traceIDStr
		}
	}

	// Try to get trace ID from context
	if r.TraceID == "" {
		if traceID := ctx.Context().Value("trace-id"); traceID != nil {
			if traceIDStr, ok := traceID.(string); ok {
				r.TraceID = traceIDStr
			}
		}
	}
	// Also try "fields" key for compatibility
	if r.TraceID == "" {
		if fields := ctx.Context().Value("fields"); fields != nil {
			if fieldsMap, ok := fields.(map[string]interface{}); ok {
				if traceID, exists := fieldsMap["TraceID"]; exists {
					if traceIDStr, ok := traceID.(string); ok {
						r.TraceID = traceIDStr
					}
				}
			}
		}
	}

	res := r.ToMap()

	if len(ctx.Errors()) > 0 {
		// Check environment from context or use a default
		env := ctx.Context().Value("environment")
		if env == nil {
			env = ctx.Context().Value("ENVIRONMENT")
		}

		// Only include errors in non-production environments
		if envStr, ok := env.(string); ok && envStr != "production" {
			// Convert WrapError to a serializable format
			errorsList := make([]map[string]interface{}, len(ctx.Errors()))
			for i, wrapErr := range ctx.Errors() {
				errorsList[i] = map[string]interface{}{
					"call_in": wrapErr.CallIn,
					"error": map[string]interface{}{
						"message":  wrapErr.Error.Message,
						"code":     wrapErr.Error.Code,
						"metadata": wrapErr.Error.Metadata,
					},
				}
			}
			res["errors"] = errorsList
			delete(res, "data")
		}
	}

	if len(r.Results) > 0 {
		delete(res, "data")
	}

	return res
}

// JSON returns the response as JSON bytes.
func (r Response[T]) JSON() ([]byte, error) {
	return json.Marshal(r)
}
