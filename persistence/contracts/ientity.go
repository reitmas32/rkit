package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
)

var (
	ErrorConvertEntityToJSON = kerrors.NewKError("Error al convertir la entidad a JSON", 500, nil)
	ErrorConvertJSONToMap    = kerrors.NewKError("Error al convertir los bytes JSON a un mapa", 500, nil)
	ErrorConvertMapToModel   = kerrors.NewKError("Error al convertir el mapa a modelo", 500, nil)
)

// --------------------------------
// DOMAIN
// --------------------------------
// IEntity
//--------------------------------

// Definimos una interfaz que represente a una entidad.
type IEntity interface {
	GetID() string
}

func ToJSON[E IEntity](entity E) []byte {
	jsonData, err := json.MarshalIndent(entity, "", "  ")
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return nil
	}

	return jsonData
}

// Función genérica que opera sobre tipos que cumplen con IEntity.
func FromJSON[E IEntity](m map[string]interface{}) (E, error) {
	var entity E

	// Convertir el mapa a bytes JSON.
	bytes, err := json.Marshal(m)
	if err != nil {
		return entity, err
	}

	// Deserializar los bytes JSON en la entidad.
	err = json.Unmarshal(bytes, &entity)
	return entity, err
}

func EntityToModel[E IEntity, M IModel](entity IEntity) result.Result[M] {
	var res map[string]interface{}

	// Convertir la entidad a JSON (bytes).
	data, err := json.Marshal(entity)
	if err != nil {
		return result.NewErrResult[M](
			ErrorConvertEntityToJSON.WithCause(err),
		)
	}

	// Convertir los bytes JSON a un mapa.
	err = json.Unmarshal(data, &res)
	if err != nil {
		return result.NewErrResult[M](
			ErrorConvertJSONToMap.WithCause(err),
		)
	}

	// Convertir el mapa a modelo.
	model, err := FromJSON[M](res)
	if err != nil {
		return result.NewErrResult[M](
			ErrorConvertMapToModel.WithCause(err),
		)
	}

	return result.NewOkResult[M](model)
}
