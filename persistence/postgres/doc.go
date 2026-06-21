// Package postgres provides a generic repository implementation over the
// persistence/contracts abstractions backed by PostgreSQL via GORM
// (gorm.io/gorm). It offers the standard operations — Save, GetByID,
// UpdateByFields, DeleteByID and criteria-based matching — typed by entity and
// model so domain code stays storage-agnostic.
//
//	import "github.com/reitmas32/rkit/persistence/postgres"
//
//	repo := &postgres.PostgresRepository[UserEntity, UserModel]{Connection: db}
package postgres
