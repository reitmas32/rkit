package contracts

import (
	"encoding/json"

	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
)

var (
	ErrorConvertModelToEntity = kerrors.NewKError("Error al convertir el modelo a entidad", 500, nil)
	ErrorConvertEntityToModel = kerrors.NewKError("Error al convertir la entidad a modelo", 500, nil)
)

// --------------------------------
// DOMAIN
// --------------------------------
// IModel
// --------------------------------
// Definimos una interfaz que represente a una entidad.
type IModel interface {
	GetID() string
	TableName() string
}

func ModelToEntity[E IEntity, M IModel](model IModel) result.Result[E] {
	var res map[string]interface{}

	// Convertir el struct a JSON (bytes).
	data, err := json.Marshal(model)
	if err != nil {
		return result.NewErrResult[E](
			ErrorConvertModelToEntity.WithCause(err),
		)
	}

	// Convertir los bytes JSON a un mapa.
	err = json.Unmarshal(data, &res)
	if err != nil {
		return result.NewErrResult[E](
			ErrorConvertEntityToModel.WithCause(err),
		)
	}

	entity, err := FromJSON[E](res)
	if err != nil {
		return result.NewErrResult[E](
			ErrorConvertEntityToModel.WithCause(err),
		)
	}

	return result.Ok(entity)
}
