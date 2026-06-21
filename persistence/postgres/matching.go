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
	if pageable == nil {
		pageable = pagination.NewPageableWithoutSort(0, 10)
	}

	if !pageable.IsValid() {
		cc.Logger().Error(ErrorPageableRequired.Detail())
		cc.AddError(ErrorPageableRequired)
		return result.Err[[]E](ErrorPageableRequired)
	}

	query := r.Connection.Model(new(M))

	query, fErr := r.applyFilters(query, crit.Filters)
	if fErr != nil {
		cc.Logger().Error(fErr.Detail())
		cc.AddError(fErr)
		return result.Err[[]E](fErr)
	}

	query, sErr := r.applySort(query, pageable)
	if sErr != nil {
		cc.Logger().Error(sErr.Detail())
		cc.AddError(sErr)
		return result.Err[[]E](sErr)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Detail())
		cc.AddError(dbErr)
		return result.Err[[]E](dbErr)
	}

	query = query.Offset(pageable.Offset()).Limit(pageable.Limit())

	var models []M
	if err := query.Find(&models).Error; err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Detail())
		cc.AddError(dbErr)
		return result.Err[[]E](dbErr)
	}

	entities, cErr := r.modelsToEntities(cc, models)
	if cErr != nil {
		return result.Err[[]E](cErr)
	}

	return result.Ok(entities)
}

// MatchingResult contiene los datos y el total de elementos que coinciden con los criterios
type MatchingResult[E contracts.IEntity] struct {
	Data  []E
	Total int64
}

// MatchingWithTotal retorna los datos y el total de elementos que coinciden con los criterios.
// Es útil cuando se necesita tanto la página de datos como el total para construir
// metadatos de paginación completos.
func (r *PostgresRepository[E, M]) MatchingWithTotal(
	cc *customctx.CustomContext,
	crit criteria.Criteria,
	pageable *pagination.Pageable,
) result.Result[MatchingResult[E]] {
	if pageable == nil {
		pageable = pagination.NewPageableWithoutSort(0, 10)
	}

	if !pageable.IsValid() {
		cc.Logger().Error(ErrorPageableRequired.Detail())
		cc.AddError(ErrorPageableRequired)
		return result.Err[MatchingResult[E]](ErrorPageableRequired)
	}

	query := r.Connection.Model(new(M))

	query, fErr := r.applyFilters(query, crit.Filters)
	if fErr != nil {
		cc.Logger().Error(fErr.Detail())
		cc.AddError(fErr)
		return result.Err[MatchingResult[E]](fErr)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Detail())
		cc.AddError(dbErr)
		return result.Err[MatchingResult[E]](dbErr)
	}

	query, sErr := r.applySort(query, pageable)
	if sErr != nil {
		cc.Logger().Error(sErr.Detail())
		cc.AddError(sErr)
		return result.Err[MatchingResult[E]](sErr)
	}

	query = query.Offset(pageable.Offset()).Limit(pageable.Limit())

	var models []M
	if err := query.Find(&models).Error; err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Detail())
		cc.AddError(dbErr)
		return result.Err[MatchingResult[E]](dbErr)
	}

	entities, cErr := r.modelsToEntities(cc, models)
	if cErr != nil {
		return result.Err[MatchingResult[E]](cErr)
	}

	return result.Ok(MatchingResult[E]{Data: entities, Total: totalCount})
}

// applyFilters validates every field against the repository FieldPolicy and
// applies the filters to the GORM query. Field names are checked before being
// used as identifiers, which prevents SQL injection through column names; values
// are always passed as bound parameters.
func (r *PostgresRepository[E, M]) applyFilters(query *gorm.DB, filters criteria.Filters) (*gorm.DB, *kerrors.KError) {
	for _, filter := range filters.Get() {
		field := string(filter.Field)
		if kErr := r.checkField(field); kErr != nil {
			return query, kErr
		}
		value := filter.Value

		switch filter.Operator {
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
			query = query.Where(fmt.Sprintf("%s LIKE ?", field), likePattern(value))
		case criteria.OperatorNotLike:
			query = query.Where(fmt.Sprintf("%s NOT LIKE ?", field), likePattern(value))
		case criteria.OperatorIn:
			query = query.Where(fmt.Sprintf("%s IN ?", field), value)
		case criteria.OperatorNotIn:
			query = query.Where(fmt.Sprintf("%s NOT IN ?", field), value)
		default:
			return query, ErrorInvalidFieldName.
				WithMetadata("field", field).
				WithMetadata("reason", fmt.Sprintf("unsupported operator %q", filter.Operator))
		}
	}

	return query, nil
}

// applySort validates the sort field and applies a safe ORDER BY clause.
// The direction is restricted to ASC/DESC and the field is validated against the
// FieldPolicy, so neither can carry injection.
func (r *PostgresRepository[E, M]) applySort(query *gorm.DB, pageable *pagination.Pageable) (*gorm.DB, *kerrors.KError) {
	if pageable.Sort == nil || !pageable.Sort.IsValid() {
		return query, nil
	}
	if kErr := r.checkField(pageable.Sort.Field); kErr != nil {
		return query, kErr.WithMetadata("clause", "sort")
	}
	direction := "ASC"
	if strings.ToUpper(string(pageable.Sort.Direction)) == "DESC" {
		direction = "DESC"
	}
	return query.Order(fmt.Sprintf("%s %s", pageable.Sort.Field, direction)), nil
}

// likePattern converts a value into a LIKE pattern, wrapping it in wildcards when
// the caller did not provide any.
func likePattern(value any) string {
	s := fmt.Sprintf("%v", value)
	if !strings.Contains(s, "%") {
		s = "%" + s + "%"
	}
	return s
}

// modelsToEntities converts a slice of models to domain entities, returning a
// descriptive error if any conversion fails.
func (r *PostgresRepository[E, M]) modelsToEntities(cc *customctx.CustomContext, models []M) ([]E, *kerrors.KError) {
	entities := make([]E, 0, len(models))
	for _, model := range models {
		entityResult := contracts.ModelToEntity[E, M](model)
		if !entityResult.IsOk() {
			kErr := entityResult.ToKError()
			if kErr == nil {
				kErr = kerrors.NewInternal("error converting model to entity", entityResult.Error())
			}
			cc.Logger().Error(kErr.Detail())
			cc.AddError(kErr)
			return nil, kErr
		}
		entities = append(entities, entityResult.Value())
	}
	return entities, nil
}
