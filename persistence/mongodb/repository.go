package mongodb

import (
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongoRepository is a generic MongoDB repository that operates on entities and models.
// E is the domain entity type (implements contracts.IEntity).
// M is the persistence model type (implements contracts.IModel).
// The collection name is taken from M.TableName().
type MongoRepository[E contracts.IEntity, M contracts.IModel] struct {
	Collection *mongo.Collection
	OnMutation models.OnMutationFunc
}

// NewMongoRepository creates a new MongoRepository using the given collection.
func NewMongoRepository[E contracts.IEntity, M contracts.IModel](
	collection *mongo.Collection,
	onMutation models.OnMutationFunc,
) *MongoRepository[E, M] {
	return &MongoRepository[E, M]{
		Collection: collection,
		OnMutation: onMutation,
	}
}
