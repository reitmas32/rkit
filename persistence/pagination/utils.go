package pagination

// ValidatePageRequest valida y normaliza una PageRequest.
// Retorna una PageRequest válida con valores por defecto si es necesario.
func ValidatePageRequest(page, size int) PageRequest {
	return NewPageRequest(page, size)
}

// CalculateOffset calcula el offset basado en page y size.
func CalculateOffset(page, size int) int {
	if page < 0 {
		page = 0
	}
	if size <= 0 {
		size = 10
	}
	return page * size
}

// CalculateTotalPages calcula el número total de páginas dado el total de elementos y el tamaño de página.
func CalculateTotalPages(totalElements, size int64) int {
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

// IsValidPage verifica si un número de página es válido dado el total de páginas.
func IsValidPage(page, totalPages int) bool {
	return page >= 0 && page < totalPages
}
