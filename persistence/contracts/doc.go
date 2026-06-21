// Package contracts defines the persistence abstractions shared by all
// repository implementations: IEntity (a domain entity) and IModel (its database
// representation), plus generic conversion helpers EntityToModel, ModelToEntity,
// ToJSON and FromJSON.
//
// Repositories under github.com/reitmas32/rkit/persistence (inmemory, postgres,
// mongodb) are written against these contracts so domain code never depends on a
// specific storage engine.
//
//	import "github.com/reitmas32/rkit/persistence/contracts"
package contracts
