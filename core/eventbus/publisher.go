package eventbus

import (
	"time"

	"github.com/reitmas32/rkit/core/customctx"
)

type Publisher interface {
	Publish(ctx *customctx.CustomContext, event Event) error

	PublishWithDelay(ctx *customctx.CustomContext, event Event, delay time.Duration) error
}
