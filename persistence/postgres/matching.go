package postgres

import (
	"fmt"
	"strings"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/criteria"
	"github.com/reitmas32/rkit/persistence/pagination"
	"gorm.io/gorm"
)

func (r *PostgresRepository[E, M]) Matching(
	cc *customctx.CustomContext,
	crit criteria.Criteria,
	pageable *pagination.Pageable,
) result.Result[[]E] {
	// Si pageable es nil, usar uno por defecto
	if pageable == nil {
		pageable = pagination.NewPageableWithoutSort(0, 10)
	}

	// Validar pageable
	if !pageable.IsValid() {
		cc.Logger().Error(ErrorPageableRequired.Error())
		cc.AddError(ErrorPageableRequired)
		return result.Err[[]E](ErrorPageableRequired)
	}

	// Crear query base
	query := r.Connection.Model(new(M))

	// Aplicar filtros
	query = r.applyFilters(query, crit.Filters)

	// Aplicar ordenamiento si está presente
	if pageable.Sort != nil && pageable.Sort.IsValid() {
		orderClause := fmt.Sprintf("%s %s", pageable.Sort.Field, strings.ToUpper(string(pageable.Sort.Direction)))
		query = query.Order(orderClause)
	}

	// Contar total de elementos que coinciden con los filtros
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[[]E](dbErr)
	}

	// Aplicar paginación
	offset := pageable.Offset()
	limit := pageable.Limit()
	query = query.Offset(offset).Limit(limit)

	// Ejecutar la consulta
	var models []M
	if err := query.Find(&models).Error; err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[[]E](dbErr)
	}

	// Convertir modelos a entidades
	entities := make([]E, 0, len(models))
	for _, model := range models {
		entityResult := contracts.ModelToEntity[E, M](model)
		if !entityResult.IsOk() {
			kErr := entityResult.ToKError()
			if kErr != nil {
				cc.Logger().Error(kErr.Error())
				cc.AddError(kErr)
				return result.Err[[]E](kErr)
			}
			// Si no es un KError, crear uno genérico
			genericErr := kerrors.NewKError("Error converting model to entity", 500, nil).WithCause(entityResult.Error())
			cc.Logger().Error(genericErr.Error())
			cc.AddError(genericErr)
			return result.Err[[]E](genericErr)
		}
		entities = append(entities, entityResult.Value())
	}

	return result.Ok(entities)
}

// MatchingResult contiene los datos y el total de elementos que coinciden con los criterios
type MatchingResult[E contracts.IEntity] struct {
	Data  []E
	Total int64
}

// MatchingWithTotal retorna los datos y el total de elementos que coinciden con los criterios.
// Este método es útil cuando necesitas tanto los datos paginados como el total de elementos
// para construir metadatos de paginación completos.
func (r *PostgresRepository[E, M]) MatchingWithTotal(
	cc *customctx.CustomContext,
	crit criteria.Criteria,
	pageable *pagination.Pageable,
) result.Result[MatchingResult[E]] {
	// Si pageable es nil, usar uno por defecto
	if pageable == nil {
		pageable = pagination.NewPageableWithoutSort(0, 10)
	}

	// Validar pageable
	if !pageable.IsValid() {
		cc.Logger().Error(ErrorPageableRequired.Error())
		cc.AddError(ErrorPageableRequired)
		return result.Err[MatchingResult[E]](ErrorPageableRequired)
	}

	// Crear query base
	query := r.Connection.Model(new(M))

	// Aplicar filtros
	query = r.applyFilters(query, crit.Filters)

	// Contar total de elementos que coinciden con los filtros (antes de aplicar paginación)
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[MatchingResult[E]](dbErr)
	}

	// Aplicar ordenamiento si está presente
	if pageable.Sort != nil && pageable.Sort.IsValid() {
		orderClause := fmt.Sprintf("%s %s", pageable.Sort.Field, strings.ToUpper(string(pageable.Sort.Direction)))
		query = query.Order(orderClause)
	}

	// Aplicar paginación
	offset := pageable.Offset()
	limit := pageable.Limit()
	query = query.Offset(offset).Limit(limit)

	// Ejecutar la consulta
	var models []M
	if err := query.Find(&models).Error; err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[MatchingResult[E]](dbErr)
	}

	// Convertir modelos a entidades
	entities := make([]E, 0, len(models))
	for _, model := range models {
		entityResult := contracts.ModelToEntity[E, M](model)
		if !entityResult.IsOk() {
			kErr := entityResult.ToKError()
			if kErr != nil {
				cc.Logger().Error(kErr.Error())
				cc.AddError(kErr)
				return result.Err[MatchingResult[E]](kErr)
			}
			// Si no es un KError, crear uno genérico
			genericErr := kerrors.NewKError("Error converting model to entity", 500, nil).WithCause(entityResult.Error())
			cc.Logger().Error(genericErr.Error())
			cc.AddError(genericErr)
			return result.Err[MatchingResult[E]](genericErr)
		}
		entities = append(entities, entityResult.Value())
	}

	return result.Ok(MatchingResult[E]{
		Data:  entities,
		Total: totalCount,
	})
}

// applyFilters aplica los filtros a la query de GORM.
func (r *PostgresRepository[E, M]) applyFilters(query *gorm.DB, filters criteria.Filters) *gorm.DB {
	for _, filter := range filters.Get() {
		field := string(filter.Field)
		operator := filter.Operator
		value := filter.Value

		switch operator {
		case criteria.OperatorEqual:
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		case criteria.OperatorNotEqual:
			query = query.Where(fmt.Sprintf("%s <> ?", field), value)
		case criteria.OperatorGreaterThan:
			query = query.Where(fmt.Sprintf("%s > ?", field), value)
		case criteria.OperatorGreaterEqual:
			query = query.Where(fmt.Sprintf("%s >= ?", field), value)
		case criteria.OperatorLessThan:
			query = query.Where(fmt.Sprintf("%s < ?", field), value)
		case criteria.OperatorLessEqual:
			query = query.Where(fmt.Sprintf("%s <= ?", field), value)
		case criteria.OperatorLike:
			// Convertir valor a string y agregar wildcards si es necesario
			likeValue := fmt.Sprintf("%v", value)
			if !strings.Contains(likeValue, "%") {
				likeValue = "%" + likeValue + "%"
			}
			query = query.Where(fmt.Sprintf("%s LIKE ?", field), likeValue)
		case criteria.OperatorNotLike:
			likeValue := fmt.Sprintf("%v", value)
			if !strings.Contains(likeValue, "%") {
				likeValue = "%" + likeValue + "%"
			}
			query = query.Where(fmt.Sprintf("%s NOT LIKE ?", field), likeValue)
		case criteria.OperatorIn:
			// Para IN, el valor debe ser un slice o array
			query = query.Where(fmt.Sprintf("%s IN ?", field), value)
		case criteria.OperatorNotIn:
			query = query.Where(fmt.Sprintf("%s NOT IN ?", field), value)
		}
	}

	return query
}
