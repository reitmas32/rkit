package inmemory

import (
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/models"
)

// --------------------------------
// INFRASTRUCTURE
// --------------------------------
// InMemoryRepository
// --------------------------------

type InMemoryMapRepository[E contracts.IEntity, M contracts.IModel] struct {
	items      map[string]M
	OnMutation models.OnMutationFunc
}

func NewInMemoryMapRepository[E contracts.IEntity, M contracts.IModel](onMutation models.OnMutationFunc) *InMemoryMapRepository[E, M] {
	return &InMemoryMapRepository[E, M]{items: make(map[string]M), OnMutation: onMutation}
}
