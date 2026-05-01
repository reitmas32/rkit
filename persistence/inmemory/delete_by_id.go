package inmemory

import (
	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/models"
)

func (r *InMemoryMapRepository[E, M]) DeleteById(cc *customctx.CustomContext, id string) result.Result[E] {
	if id == "" {
		cc.Logger().Error(ErrorItemIDRequired.Error())
		cc.AddError(ErrorItemIDRequired)
		return result.Err[E](ErrorItemIDRequired)
	}

	// Obtener el modelo antes de eliminarlo para notificar la mutación
	model, ok := r.items[id]
	if !ok {
		cc.Logger().Error(ErrorItemNotFound.Error())
		cc.AddError(ErrorItemNotFound)
		return result.Err[E](ErrorItemNotFound)
	}

	// Notificar mutación antes de eliminar
	mutationResult := models.NotifyMutation[E, M](cc, r.OnMutation, model, "delete")
	if !mutationResult.IsOk() {
		return result.Err[E](mutationResult.ToKError())
	}

	// Eliminar el item
	delete(r.items, id)
	return result.Empty[E]()
}
