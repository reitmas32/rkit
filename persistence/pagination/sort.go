package pagination

// SortDirection representa la dirección del ordenamiento.
type SortDirection string

const (
	// SortDirectionASC ordena de forma ascendente.
	SortDirectionASC SortDirection = "ASC"

	// SortDirectionDESC ordena de forma descendente.
	SortDirectionDESC SortDirection = "DESC"
)

// Sort representa información de ordenamiento.
type Sort struct {
	// Field es el campo por el cual ordenar.
	Field string

	// Direction es la dirección del ordenamiento (ASC o DESC).
	Direction SortDirection
}

// NewSort crea un nuevo Sort con validación.
// Si direction no es válido, se establece en ASC por defecto.
func NewSort(field string, direction SortDirection) *Sort {
	if direction != SortDirectionASC && direction != SortDirectionDESC {
		direction = SortDirectionASC
	}
	return &Sort{
		Field:     field,
		Direction: direction,
	}
}

// NewSortASC crea un nuevo Sort con dirección ascendente.
func NewSortASC(field string) *Sort {
	return NewSort(field, SortDirectionASC)
}

// NewSortDESC crea un nuevo Sort con dirección descendente.
func NewSortDESC(field string) *Sort {
	return NewSort(field, SortDirectionDESC)
}

// IsValid verifica si el Sort es válido.
func (s *Sort) IsValid() bool {
	return s != nil && s.Field != "" && (s.Direction == SortDirectionASC || s.Direction == SortDirectionDESC)
}

