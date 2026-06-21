// Package pagination provides page-oriented query types. PageRequest captures
// the requested page and size (with Offset/Limit helpers); Pageable adds an
// optional Sort (field + SortDirection); and PageResult[T] wraps a page of
// content together with total counts and the computed number of pages.
// CriteriaPageable combines a criteria.Criteria with a Pageable for repository
// queries.
//
//	import "github.com/reitmas32/rkit/persistence/pagination"
//
//	req := pagination.NewPageRequest(0, 20)
//	page := pagination.NewPageResult(items, total, req)
package pagination
