package loguru

type Fields interface {
	ToFields() map[string]any
	UpdateAll(fields map[string]any)

	UpdateOne(key string, value any)

	Format() string
}
