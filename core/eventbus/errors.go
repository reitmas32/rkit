package eventbus

import "github.com/reitmas32/rkit/core/kerrors"

var (
	ErrEventNotFound          = kerrors.NewKError("event not found", 404, nil)
	ErrEventAlreadyExists     = kerrors.NewKError("event already exists", 400, nil)
	ErrEventInvalid           = kerrors.NewKError("event invalid", 400, nil)
	ErrEventFailedToPublish   = kerrors.NewKError("event failed to publish", 500, nil)
	ErrEventFailedToSubscribe = kerrors.NewKError("event failed to subscribe", 500, nil)
)
