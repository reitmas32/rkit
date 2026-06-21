// Package types provides base domain value types, most notably DomainID: a
// self-describing identifier that encodes a module code and an entity code
// alongside a UUID. This yields stable, traceable IDs whose origin can be read
// back from the value itself.
//
// Construct IDs with NewDomainID, parse them with ParseDomainID, and inspect
// their parts via Module, Entity and UUID.
//
//	import "github.com/reitmas32/rkit/core/types"
//
//	id, err := types.NewDomainID("02", "01") // module 02, entity 01
package types
