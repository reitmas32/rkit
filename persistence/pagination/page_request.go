package pagination

// PageRequest representa una solicitud de paginación.
// Contiene la información necesaria para paginar resultados.
type PageRequest struct {
	// Page es el número de página (empezando desde 0 o 1, dependiendo de la implementación).
	// Por convención, usaremos 0 como primera página.
	Page int

	// Size es el tamaño de la página (número de elementos por página).
	Size int
}

// NewPageRequest crea una nueva PageRequest con validación.
// Si page < 0, se establece en 0.
// Si size <= 0, se establece en un valor por defecto (10).
func NewPageRequest(page, size int) PageRequest {
	if page < 0 {
		page = 0
	}
	if size <= 0 {
		size = 10 // Tamaño por defecto
	}
	return PageRequest{
		Page: page,
		Size: size,
	}
}

// Offset calcula el offset (desplazamiento) para la consulta.
// Retorna el número de elementos a saltar antes de comenzar a retornar resultados.
func (pr PageRequest) Offset() int {
	return pr.Page * pr.Size
}

// Limit retorna el límite de elementos a retornar (equivalente a Size).
func (pr PageRequest) Limit() int {
	return pr.Size
}

// IsValid verifica si la PageRequest es válida.
func (pr PageRequest) IsValid() bool {
	return pr.Page >= 0 && pr.Size > 0
}

