package postgres

import (
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/models"
	"gorm.io/gorm"
)

type PostgresRepository[E contracts.IEntity, M contracts.IModel] struct {
	Connection *gorm.DB
	OnMutation models.OnMutationFunc
}
