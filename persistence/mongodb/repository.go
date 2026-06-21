package mongodb

import (
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/criteria"
	"github.com/reitmas32/rkit/persistence/models"
	"github.com/reitmas32/rkit/persistence/pagination"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongoRepository is a generic MongoDB repository that operates on entities and models.
// E is the domain entity type (implements contracts.IEntity).
// M is the persistence model type (implements contracts.IModel).
// The collection name is taken from M.TableName().
type MongoRepository[E contracts.IEntity, M contracts.IModel] struct {
	Collection *mongo.Collection
	OnMutation models.OnMutationFunc

	// FieldPolicy controls which field names are accepted in criteria filters
	// and sort clauses. The zero value is secure: every field is validated
	// against a strict identifier pattern (criteria.IsValidIdentifier), which
	// keeps client-controlled field names from being used as query operator keys
	// (e.g. "$where"). Set FieldPolicy.Allowed to restrict queries to a known set
	// of fields.
	FieldPolicy criteria.FieldPolicy
}

// checkField validates a single filter/sort field against the FieldPolicy.
func (r *MongoRepository[E, M]) checkField(field string) *kerrors.KError {
	if r.FieldPolicy.Permits(field) {
		return nil
	}
	reason := "field is not a valid identifier"
	if criteria.IsValidIdentifier(field) {
		reason = "field is not in the configured allow-list"
	}
	return ErrorInvalidFieldName.
		WithMetadata("field", field).
		WithMetadata("reason", reason)
}

// validateCriteria checks every filter field and the sort field up front, so a
// malicious or unknown field name is rejected with a clear error before any
// query is built.
func (r *MongoRepository[E, M]) validateCriteria(crit criteria.Criteria, pageable *pagination.Pageable) *kerrors.KError {
	for _, f := range crit.Filters.Get() {
		if kErr := r.checkField(string(f.Field)); kErr != nil {
			return kErr
		}
	}
	if pageable != nil && pageable.Sort != nil && pageable.Sort.IsValid() {
		if kErr := r.checkField(pageable.Sort.Field); kErr != nil {
			return kErr.WithMetadata("clause", "sort")
		}
	}
	return nil
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
