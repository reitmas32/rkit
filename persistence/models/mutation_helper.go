package models

import (
	"encoding/json"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
)

// NotifyMutation convierte un modelo a entidad y notifica la mutación a través de OnMutationFunc.
// Esta función es reutilizable para operaciones de save, update y delete.
func NotifyMutation[E contracts.IEntity, M contracts.IModel](
	cc *customctx.CustomContext,
	onMutation OnMutationFunc,
	model M,
	operation string,
) result.Result[bool] {
	if onMutation == nil {
		return result.Ok(true)
	}

	// Convertir modelo a entidad
	resultEntity := contracts.ModelToEntity[E, M](model)
	if !resultEntity.IsOk() {
		kErr := resultEntity.ToKError()
		if kErr != nil {
			cc.Logger().Error(kErr.Error())
			cc.AddError(kErr)
			return result.Err[bool](kErr)
		}
		// Si no es un KError, crear uno genérico
		genericErr := kerrors.NewKError("Error converting model to entity", 500, nil).WithCause(resultEntity.Error())
		cc.Logger().Error(genericErr.Error())
		cc.AddError(genericErr)
		return result.Err[bool](genericErr)
	}

	// Convertir entidad a JSON
	entityJSON := contracts.ToJSON[E](resultEntity.Value())

	// Convertir JSON a mapa
	var entityMap map[string]interface{}
	err := json.Unmarshal(entityJSON, &entityMap)
	if err != nil {
		cc.Logger().Error(err.Error())
		unmarshalErr := kerrors.NewKError("Error unmarshalling entity to map", 500, nil).WithCause(err)
		cc.AddError(unmarshalErr)
		return result.Err[bool](unmarshalErr)
	}

	// Obtener el nombre de la tabla del modelo
	// Intentar obtener TableName del modelo si está disponible
	tableName := ""
	if tableNameProvider, ok := any(model).(interface{ TableName() string }); ok {
		tableName = tableNameProvider.TableName()
	}

	// Notificar la mutación
	err = onMutation(cc, tableName, operation, entityMap)
	if err != nil {
		cc.Logger().Error(err.Error())
		mutationErr := kerrors.NewKError("Error on mutation", 500, nil).WithCause(err)
		cc.AddError(mutationErr)
		return result.Err[bool](mutationErr)
	}

	return result.Ok(true)
}
