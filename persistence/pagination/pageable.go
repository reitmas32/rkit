package pagination

// Pageable es una interfaz que combina paginación con criterios de búsqueda.
// Permite especificar tanto los parámetros de paginación como los filtros.
type Pageable struct {
	// PageRequest contiene los parámetros de paginación.
	PageRequest PageRequest

	// Sort contiene información de ordenamiento (opcional).
	// Puede ser nil si no se requiere ordenamiento.
	Sort *Sort
}

// NewPageable crea un nuevo Pageable con PageRequest y Sort opcional.
func NewPageable(pageRequest PageRequest, sort *Sort) *Pageable {
	return &Pageable{
		PageRequest: pageRequest,
		Sort:        sort,
	}
}

// NewPageableWithoutSort crea un nuevo Pageable sin ordenamiento.
func NewPageableWithoutSort(page, size int) *Pageable {
	return &Pageable{
		PageRequest: NewPageRequest(page, size),
		Sort:        nil,
	}
}

// Offset retorna el offset calculado de la PageRequest.
func (p *Pageable) Offset() int {
	return p.PageRequest.Offset()
}

// Limit retorna el límite de la PageRequest.
func (p *Pageable) Limit() int {
	return p.PageRequest.Limit()
}

// IsValid verifica si el Pageable es válido.
func (p *Pageable) IsValid() bool {
	return p.PageRequest.IsValid()
}

