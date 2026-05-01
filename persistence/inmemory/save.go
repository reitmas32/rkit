package inmemory

import (
	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/models"
)

func (r *InMemoryMapRepository[E, M]) Save(cc *customctx.CustomContext, item E) result.Result[E] {
	// Validar que el item tenga ID
	if item.GetID() == "" {
		cc.Logger().Error(ErrorItemIDRequired.Error())
		cc.AddError(ErrorItemIDRequired)
		return result.Err[E](ErrorItemIDRequired)
	}

	// Convertir entidad a modelo para guardarlo
	modelResult := contracts.EntityToModel[E, M](item)
	if !modelResult.IsOk() {
		kErr := modelResult.ToKError()
		if kErr != nil {
			cc.Logger().Error(kErr.Error())
			cc.AddError(kErr)
			return result.Err[E](kErr)
		}
		// Si no es un KError, crear uno genérico
		genericErr := kerrors.NewKError("Error converting entity to model", 500, nil).WithCause(modelResult.Error())
		cc.Logger().Error(genericErr.Error())
		cc.AddError(genericErr)
		return result.Err[E](genericErr)
	}

	model := modelResult.Value()

	// Guardar el modelo
	r.items[model.GetID()] = model

	// Convertir el modelo guardado de vuelta a entidad para retornarlo
	entityResult := contracts.ModelToEntity[E, M](model)
	if !entityResult.IsOk() {
		kErr := entityResult.ToKError()
		if kErr != nil {
			cc.Logger().Error(kErr.Error())
			cc.AddError(kErr)
			return result.Err[E](kErr)
		}
		// Si no es un KError, crear uno genérico
		genericErr := kerrors.NewKError("Error converting model to entity", 500, nil).WithCause(entityResult.Error())
		cc.Logger().Error(genericErr.Error())
		cc.AddError(genericErr)
		return result.Err[E](genericErr)
	}

	// Notificar mutación si está configurado
	mutationResult := models.NotifyMutation[E, M](cc, r.OnMutation, model, "save")
	if !mutationResult.IsOk() {
		cc.Logger().Error(mutationResult.ToKError().Error())
	}

	return result.Ok(entityResult.Value())
}
