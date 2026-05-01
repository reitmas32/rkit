package inmemory

import (
	"encoding/json"
	"fmt"
	"reflect"
	stdSort "sort"
	"strings"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/criteria"
	"github.com/reitmas32/rkit/persistence/pagination"
)

func (r *InMemoryMapRepository[E, M]) Matching(
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

	// Convertir todos los items a slice para poder filtrarlos
	allItems := make([]M, 0, len(r.items))
	for _, item := range r.items {
		allItems = append(allItems, item)
	}

	// Aplicar filtros
	filteredItems, err := r.applyFilters(cc, allItems, crit.Filters)
	if err != nil {
		cc.Logger().Error(ErrorApplyFilters.Error())
		cc.AddError(ErrorApplyFilters.WithCause(err))
		return result.Err[[]E](ErrorApplyFilters.WithCause(err))
	}

	// Aplicar ordenamiento si está presente
	if pageable.Sort != nil && pageable.Sort.IsValid() {
		filteredItems = r.applySort(filteredItems, pageable.Sort)
	}

	// Aplicar paginación
	offset := pageable.Offset()
	limit := pageable.Limit()

	// Calcular índices para la paginación
	start := offset
	end := offset + limit
	if start > len(filteredItems) {
		start = len(filteredItems)
	}
	if end > len(filteredItems) {
		end = len(filteredItems)
	}

	// Extraer la página solicitada
	var pageContent []M
	if start < end {
		pageContent = filteredItems[start:end]
	} else {
		pageContent = []M{}
	}

	// Convertir modelos a entidades
	entities := make([]E, 0, len(pageContent))
	for _, model := range pageContent {
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

// applyFilters aplica los filtros a los items y retorna los que coinciden.
func (r *InMemoryMapRepository[E, M]) applyFilters(
	cc *customctx.CustomContext,
	items []M,
	filters criteria.Filters,
) ([]M, *kerrors.KError) {
	if len(filters.Get()) == 0 {
		return items, nil
	}

	filtered := make([]M, 0)

	for _, item := range items {
		// Convertir el Model a JSON y luego a un mapa
		data, err := json.Marshal(item)
		if err != nil {
			cc.Logger().Error(ErrorConvertModelToJSON.Error())
			cc.AddError(ErrorConvertModelToJSON.WithCause(err))
			return nil, ErrorConvertModelToJSON.WithCause(err)
		}

		var itemMap map[string]interface{}
		err = json.Unmarshal(data, &itemMap)
		if err != nil {
			cc.Logger().Error(ErrorConvertJSONToMap.Error())
			cc.AddError(ErrorConvertJSONToMap.WithCause(err))
			return nil, ErrorConvertJSONToMap.WithCause(err)
		}

		// Verificar si el item coincide con todos los filtros
		matches := true
		for _, filter := range filters.Get() {
			if !r.matchesFilter(itemMap, filter) {
				matches = false
				break
			}
		}

		if matches {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}

// matchesFilter verifica si un item (representado como mapa) coincide con un filtro.
func (r *InMemoryMapRepository[E, M]) matchesFilter(itemMap map[string]interface{}, filter criteria.Filter) bool {
	fieldValue, exists := itemMap[string(filter.Field)]
	if !exists {
		return false
	}

	switch filter.Operator {
	case criteria.OperatorEqual:
		return compareValues(fieldValue, filter.Value) == 0
	case criteria.OperatorNotEqual:
		return compareValues(fieldValue, filter.Value) != 0
	case criteria.OperatorGreaterThan:
		return compareValues(fieldValue, filter.Value) > 0
	case criteria.OperatorGreaterEqual:
		return compareValues(fieldValue, filter.Value) >= 0
	case criteria.OperatorLessThan:
		return compareValues(fieldValue, filter.Value) < 0
	case criteria.OperatorLessEqual:
		return compareValues(fieldValue, filter.Value) <= 0
	case criteria.OperatorLike:
		return matchesLike(fieldValue, filter.Value)
	case criteria.OperatorNotLike:
		return !matchesLike(fieldValue, filter.Value)
	case criteria.OperatorIn:
		return matchesIn(fieldValue, filter.Value)
	case criteria.OperatorNotIn:
		return !matchesIn(fieldValue, filter.Value)
	default:
		return false
	}
}

// compareValues compara dos valores y retorna -1, 0, o 1.
func compareValues(a, b interface{}) int {
	// Convertir a strings para comparación simple
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)

	// Intentar comparación numérica si ambos son números
	if aNum, aOk := toNumber(a); aOk {
		if bNum, bOk := toNumber(b); bOk {
			if aNum < bNum {
				return -1
			}
			if aNum > bNum {
				return 1
			}
			return 0
		}
	}

	// Comparación de strings
	if aStr < bStr {
		return -1
	}
	if aStr > bStr {
		return 1
	}
	return 0
}

// toNumber intenta convertir un valor a float64.
func toNumber(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

// matchesLike verifica si un valor coincide con un patrón LIKE.
func matchesLike(fieldValue, pattern interface{}) bool {
	fieldStr := fmt.Sprintf("%v", fieldValue)
	patternStr := fmt.Sprintf("%v", pattern)

	// Convertir patrón SQL LIKE a regex simple
	// % -> .*, _ -> .
	regexPattern := strings.ReplaceAll(patternStr, "%", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "_", ".")

	// Comparación simple con strings.Contains para patrones básicos
	if strings.Contains(patternStr, "%") {
		// Patrón con wildcards
		if strings.HasPrefix(patternStr, "%") && strings.HasSuffix(patternStr, "%") {
			// %valor%
			middle := strings.Trim(patternStr, "%")
			return strings.Contains(fieldStr, middle)
		}
		if strings.HasPrefix(patternStr, "%") {
			// %valor
			suffix := strings.TrimPrefix(patternStr, "%")
			return strings.HasSuffix(fieldStr, suffix)
		}
		if strings.HasSuffix(patternStr, "%") {
			// valor%
			prefix := strings.TrimSuffix(patternStr, "%")
			return strings.HasPrefix(fieldStr, prefix)
		}
	}

	return fieldStr == patternStr
}

// matchesIn verifica si un valor está en una lista de valores.
func matchesIn(fieldValue, inValue interface{}) bool {
	// Convertir inValue a slice si es posible
	val := reflect.ValueOf(inValue)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return false
	}

	fieldStr := fmt.Sprintf("%v", fieldValue)
	for i := 0; i < val.Len(); i++ {
		elemStr := fmt.Sprintf("%v", val.Index(i).Interface())
		if fieldStr == elemStr {
			return true
		}
	}
	return false
}

// applySort ordena los items según el Sort especificado.
func (r *InMemoryMapRepository[E, M]) applySort(items []M, sort *pagination.Sort) []M {
	if sort == nil || !sort.IsValid() {
		return items
	}

	// Crear una copia para no modificar el slice original
	sorted := make([]M, len(items))
	copy(sorted, items)

	stdSort.Slice(sorted, func(i, j int) bool {
		// Convertir ambos items a mapas
		iData, _ := json.Marshal(sorted[i])
		jData, _ := json.Marshal(sorted[j])

		var iMap, jMap map[string]interface{}
		json.Unmarshal(iData, &iMap)
		json.Unmarshal(jData, &jMap)

		iValue := iMap[sort.Field]
		jValue := jMap[sort.Field]

		comparison := compareValues(iValue, jValue)
		if sort.Direction == pagination.SortDirectionASC {
			return comparison < 0
		}
		return comparison > 0
	})

	return sorted
}
