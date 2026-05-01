package eventbus

import "time"

type Event interface {
	Name() string
	Version() string
	OccurredAt() time.Time
	Payload() any
	Metadata() Metadata
}
