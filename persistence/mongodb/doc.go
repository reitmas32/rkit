// Package mongodb provides a generic repository implementation over the
// persistence/contracts abstractions backed by MongoDB
// (go.mongodb.org/mongo-driver/v2). It offers the standard operations — Save,
// GetByID, UpdateByFields, DeleteByID and criteria-based matching — so domain
// code can target MongoDB through the same contracts as the other backends.
//
//	import "github.com/reitmas32/rkit/persistence/mongodb"
package mongodb
