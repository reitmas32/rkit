package models

import (
	"github.com/reitmas32/rkit/core/customctx"
)

type OnMutationFunc func(cc *customctx.CustomContext, table string, operation string, model map[string]interface{}) error
