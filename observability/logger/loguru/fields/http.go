package fields

import (
	"encoding/json"
	"fmt"
)

type HTTPFileds struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	TraceID  string `json:"trace_id"`
	CallerID string `json:"caller_id"`
	ClientIP string `json:"client_ip"`

	Metadata map[string]any `json:"metadata"`
}

func (c *HTTPFileds) Format() string {
	base := fmt.Sprintf(
		"method: %s | path: %s | trace_id: %s | caller_id: %s | client_ip: %s",
		c.Method, c.Path, c.TraceID, c.CallerID, c.ClientIP,
	)

	// Si está vacío: no agregamos metadata
	if len(c.Metadata) == 0 {
		return base
	}

	// metadata como JSON
	metaBytes, err := json.Marshal(c.Metadata)
	if err != nil {
		// fallback: lo imprime "como pueda"
		return base + fmt.Sprintf(" | metadata: %v", c.Metadata)
	}

	return base + fmt.Sprintf(" | metadata: %s", string(metaBytes))
}

func (c *HTTPFileds) ToFields() map[string]any {
	return toFields(c)
}

func (c *HTTPFileds) UpdateAll(fields map[string]any) {
	updateAll(c, fields)
}

func (c *HTTPFileds) UpdateOne(key string, value any) {
	updateOne(c, key, value)
}
