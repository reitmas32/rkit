package pagination

// PageResult representa el resultado de una consulta paginada.
// Contiene los elementos de la página actual junto con información de paginación.
type PageResult[T any] struct {
	// Content contiene los elementos de la página actual.
	Content []T

	// TotalElements es el número total de elementos que coinciden con los criterios.
	TotalElements int64

	// TotalPages es el número total de páginas disponibles.
	TotalPages int

	// Page es el número de la página actual (0-indexed).
	Page int

	// Size es el tamaño de la página.
	Size int

	// IsFirst indica si esta es la primera página.
	IsFirst bool

	// IsLast indica si esta es la última página.
	IsLast bool

	// HasNext indica si hay una página siguiente.
	HasNext bool

	// HasPrevious indica si hay una página anterior.
	HasPrevious bool
}

// NewPageResult crea un nuevo PageResult con los valores calculados automáticamente.
func NewPageResult[T any](content []T, totalElements int64, pageRequest PageRequest) PageResult[T] {
	totalPages := calculateTotalPages(totalElements, int64(pageRequest.Size))
	page := pageRequest.Page

	return PageResult[T]{
		Content:       content,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          page,
		Size:          pageRequest.Size,
		IsFirst:       page == 0,
		IsLast:        page >= totalPages-1 || totalPages == 0,
		HasNext:       page < totalPages-1,
		HasPrevious:   page > 0,
	}
}

// EmptyPageResult crea un PageResult vacío para una PageRequest dada.
func EmptyPageResult[T any](pageRequest PageRequest) PageResult[T] {
	return PageResult[T]{
		Content:       []T{},
		TotalElements: 0,
		TotalPages:    0,
		Page:          pageRequest.Page,
		Size:          pageRequest.Size,
		IsFirst:       true,
		IsLast:        true,
		HasNext:       false,
		HasPrevious:   false,
	}
}

// calculateTotalPages calcula el número total de páginas.
func calculateTotalPages(totalElements, size int64) int {
	if size == 0 {
		return 0
	}
	if totalElements == 0 {
		return 0
	}
	pages := int(totalElements / size)
	if totalElements%size > 0 {
		pages++
	}
	return pages
}

