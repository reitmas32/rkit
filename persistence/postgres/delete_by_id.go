package postgres

import (
	"errors"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/models"
	"gorm.io/gorm"
)

func (r *PostgresRepository[E, M]) DeleteById(cc *customctx.CustomContext, id string) result.Result[E] {
	// Validar que el ID no esté vacío
	if id == "" {
		cc.Logger().Error(ErrorItemIDRequired.Error())
		cc.AddError(ErrorItemIDRequired)
		return result.Err[E](ErrorItemIDRequired)
	}

	// Crear una instancia del modelo para buscar y eliminar
	var model M

	// Buscar el item por ID primero para verificar que existe
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

	// Notificar mutación antes de eliminar
	mutationResult := models.NotifyMutation[E, M](cc, r.OnMutation, model, "delete")
	if !mutationResult.IsOk() {
		cc.Logger().Error(mutationResult.ToKError().Error())
	}

	// Eliminar el item usando GORM
	if err := r.Connection.Delete(&model).Error; err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[E](dbErr)
	}

	// Retornar resultado vacío después de eliminar
	return result.Empty[E]()
}
