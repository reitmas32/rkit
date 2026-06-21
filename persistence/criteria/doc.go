// Package criteria provides a composable, ORM-agnostic way to express query
// filters. A Filter pairs a FilterField with an Operator (OperatorEqual,
// OperatorGreaterThan, OperatorLike, OperatorIn, ...) and a value; Filters groups
// several of them, and Criteria carries the filter set into a repository query.
//
// Building queries from criteria keeps domain and application code free of
// SQL/driver details while repository implementations translate the criteria to
// their backend.
//
//	import "github.com/reitmas32/rkit/persistence/criteria"
//
//	filters := criteria.NewFilters([]criteria.Filter{
//	    {Field: "status", Operator: criteria.OperatorEqual, Value: "active"},
//	    {Field: "age", Operator: criteria.OperatorGreaterThan, Value: 18},
//	})
package criteria
