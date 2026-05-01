package fields

import "encoding/json"

func toFields[T any](cfg *T) map[string]any {
	b, err := json.Marshal(cfg)
	if err != nil {
		return map[string]any{}
	}

	out := map[string]any{}
	if err := json.Unmarshal(b, &out); err != nil {
		return map[string]any{}
	}
	return out
}

func updateAll[T any](cfg *T, fields map[string]any) {
	b, err := json.Marshal(fields)
	if err != nil {
		return
	}
	_ = json.Unmarshal(b, cfg)
}

func updateOne[T any](cfg *T, key string, value any) {
	updateAll(cfg, map[string]any{
		key: value,
	})
}
