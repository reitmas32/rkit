package pagination

import "github.com/reitmas32/rkit/persistence/criteria"

// CriteriaPageable combina Criteria (filtros) con Pageable (paginación y ordenamiento).
// Permite realizar consultas paginadas con filtros y ordenamiento.
type CriteriaPageable struct {
	// Criteria contiene los filtros a aplicar.
	Criteria criteria.Criteria

	// Pageable contiene los parámetros de paginación y ordenamiento.
	Pageable *Pageable
}

// NewCriteriaPageable crea un nuevo CriteriaPageable.
func NewCriteriaPageable(criteria criteria.Criteria, pageable *Pageable) *CriteriaPageable {
	return &CriteriaPageable{
		Criteria: criteria,
		Pageable: pageable,
	}
}

// NewCriteriaPageableSimple crea un CriteriaPageable con paginación simple (sin ordenamiento).
func NewCriteriaPageableSimple(criteria criteria.Criteria, page, size int) *CriteriaPageable {
	return &CriteriaPageable{
		Criteria: criteria,
		Pageable: NewPageableWithoutSort(page, size),
	}
}

// Offset retorna el offset calculado.
func (cp *CriteriaPageable) Offset() int {
	if cp.Pageable == nil {
		return 0
	}
	return cp.Pageable.Offset()
}

// Limit retorna el límite.
func (cp *CriteriaPageable) Limit() int {
	if cp.Pageable == nil {
		return 10 // Valor por defecto
	}
	return cp.Pageable.Limit()
}

// IsValid verifica si el CriteriaPageable es válido.
func (cp *CriteriaPageable) IsValid() bool {
	return cp.Pageable != nil && cp.Pageable.IsValid()
}
