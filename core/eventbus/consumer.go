package eventbus

import (
	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/result"
)

type Consumer interface {
	Consume(ctx *customctx.CustomContext, event Event) result.Result[DeliveryChannel]
}
