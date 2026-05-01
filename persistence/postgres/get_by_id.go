package postgres

import (
	"errors"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
	"gorm.io/gorm"
)

func (r *PostgresRepository[E, M]) GetById(cc *customctx.CustomContext, id string) result.Result[E] {
	// Validar que el ID no esté vacío
	if id == "" {
		cc.Logger().Error(ErrorItemIDRequired.Error())
		cc.AddError(ErrorItemIDRequired)
		return result.Err[E](ErrorItemIDRequired)
	}

	// Crear una instancia del modelo para buscar
	var model M

	// Buscar el item por ID usando GORM
	if err := r.Connection.Where("id = ?", id).First(&model).Error; err != nil {
		// Si el error es "record not found", retornar ErrorItemNotFound
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cc.Logger().Error(ErrorItemNotFound.Error())
			cc.AddError(ErrorItemNotFound)
			return result.Err[E](ErrorItemNotFound)
		}

		// Otros errores de base de datos
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[E](dbErr)
	}

	// Convertir el modelo a entidad para retornarlo
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

	return result.Ok(entityResult.Value())
}
