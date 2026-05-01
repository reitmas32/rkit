package inmemory

import (
	"encoding/json"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/models"
)

func (r *InMemoryMapRepository[E, M]) UpdateByFields(cc *customctx.CustomContext, id string, fields map[string]any) result.Result[E] {
	if id == "" {
		cc.Logger().Error(ErrorItemIDRequired.Error())
		cc.AddError(ErrorItemIDRequired)
		return result.Err[E](ErrorItemIDRequired)
	}

	if len(fields) == 0 {
		cc.Logger().Error(ErrorItemFieldsRequired.Error())
		cc.AddError(ErrorItemFieldsRequired)
		return result.Err[E](ErrorItemFieldsRequired)
	}

	// Buscar el item por ID
	item, ok := r.items[id]
	if !ok {
		cc.Logger().Error(ErrorItemNotFound.Error())
		cc.AddError(ErrorItemNotFound)
		return result.Err[E](ErrorItemNotFound)
	}

	// Convertir el Model a JSON y luego a un mapa
	data, err := json.Marshal(item)
	if err != nil {
		cc.Logger().Error(ErrorConvertModelToJSON.Error())
		cc.AddError(ErrorConvertModelToJSON.WithCause(err))
		return result.Err[E](ErrorConvertModelToJSON.WithCause(err))
	}

	var itemMap map[string]interface{}
	err = json.Unmarshal(data, &itemMap)
	if err != nil {
		cc.Logger().Error(ErrorConvertJSONToMap.Error())
		cc.AddError(ErrorConvertJSONToMap.WithCause(err))
		return result.Err[E](ErrorConvertJSONToMap.WithCause(err))
	}

	// Actualizar el mapa con los nuevos valores
	for key, value := range fields {
		itemMap[key] = value
	}

	// Asegurar que el ID se mantiene
	itemMap["id"] = id

	// Convertir el mapa actualizado de vuelta a Model
	updatedModel, err := contracts.FromJSON[M](itemMap)
	if err != nil {
		cc.Logger().Error(ErrorConvertMapToModel.Error())
		cc.AddError(ErrorConvertMapToModel.WithCause(err))
		return result.Err[E](ErrorConvertMapToModel.WithCause(err))
	}

	// Guardar el item actualizado
	r.items[id] = updatedModel

	// Convertir el modelo actualizado de vuelta a entidad para retornarlo
	entityResult := contracts.ModelToEntity[E, M](updatedModel)
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
	mutationResult := models.NotifyMutation[E, M](cc, r.OnMutation, updatedModel, "update")
	if !mutationResult.IsOk() {
		return result.Err[E](mutationResult.ToKError())
	}

	return result.Ok(entityResult.Value())
}
