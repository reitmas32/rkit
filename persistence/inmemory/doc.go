// Package inmemory provides a generic, in-process repository implementation over
// the persistence/contracts abstractions. It stores entities in memory and is
// intended for unit tests and prototyping, where a real database would add
// friction.
//
//	import "github.com/reitmas32/rkit/persistence/inmemory"
//
// For production storage use github.com/reitmas32/rkit/persistence/postgres or
// github.com/reitmas32/rkit/persistence/mongodb.
package inmemory
