package eventbus

import "context"

type HandlerFunc func(ctx context.Context, event Event) error

type Metadata map[string]string
