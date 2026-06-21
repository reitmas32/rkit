package postgres

import (
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/criteria"
	"github.com/reitmas32/rkit/persistence/models"
	"gorm.io/gorm"
)

type PostgresRepository[E contracts.IEntity, M contracts.IModel] struct {
	Connection *gorm.DB
	OnMutation models.OnMutationFunc

	// FieldPolicy controls which field names are accepted in criteria filters
	// and sort clauses. The zero value is secure: every field is validated
	// against a strict identifier pattern (criteria.IsValidIdentifier), which
	// prevents SQL injection through column names. Set FieldPolicy.Allowed to
	// restrict queries to a known column allow-list — recommended whenever field
	// names can originate from client input.
	FieldPolicy criteria.FieldPolicy
}

// checkField validates a filter/sort field against the repository's FieldPolicy
// and returns a descriptive error (with the offending field and the reason in
// metadata) when it is rejected.
func (r *PostgresRepository[E, M]) checkField(field string) *kerrors.KError {
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
