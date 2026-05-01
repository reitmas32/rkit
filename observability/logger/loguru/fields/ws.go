package fields

import (
	"encoding/json"
	"fmt"
)

type WSFields struct {
	ServerID string `json:"server_id"`
	ClientID string `json:"client_id"`
	TraceID  string `json:"trace_id"`
	Path     string `json:"path"`

	Metadata map[string]any `json:"metadata"`
}

func (c *WSFields) Format() string {
	base := fmt.Sprintf(
		"server_id: %s | client_id: %s | trace_id: %s | path: %s",
		c.ServerID, c.ClientID, c.TraceID, c.Path,
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

func (c *WSFields) ToFields() map[string]any {
	return toFields(c)
}

func (c *WSFields) UpdateAll(fields map[string]any) {
	updateAll(c, fields)
}

func (c *WSFields) UpdateOne(key string, value any) {
	updateOne(c, key, value)
}
