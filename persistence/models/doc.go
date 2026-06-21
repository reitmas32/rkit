// Package models provides the base Entity type embedded by domain entities
// (ID, timestamps, soft-delete flag and a ToJSON helper) and a mutation
// notification mechanism (OnMutationFunc, NotifyMutation) so repositories can
// emit change events for inserts, updates and deletes.
//
//	import "github.com/reitmas32/rkit/persistence/models"
//
//	type UserEntity struct {
//	    models.Entity
//	    Email string
//	}
package models
