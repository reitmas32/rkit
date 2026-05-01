package inmemory

import (
	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
)

func (r *InMemoryMapRepository[E, M]) GetById(cc *customctx.CustomContext, id string) result.Result[E] {
	if id == "" {
		cc.Logger().Error(ErrorItemIDRequired.Error())
		cc.AddError(ErrorItemIDRequired)
		return result.Err[E](ErrorItemIDRequired)
	}

	item, ok := r.items[id]
	if !ok {
		cc.Logger().Error(ErrorItemNotFound.Error())
		cc.AddError(ErrorItemNotFound)
		return result.Err[E](ErrorItemNotFound)
	}

	result := contracts.ModelToEntity[E, M](item)
	if result.ToKError() != nil {
		cc.Logger().Error(result.ToKError().Error())
		cc.AddError(result.ToKError())
	}

	return result
}
