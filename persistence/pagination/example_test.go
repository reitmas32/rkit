package pagination_test

import (
	"fmt"

	"github.com/reitmas32/rkit/persistence/pagination"
)

// ExamplePageRequest derives SQL-style offset/limit from a (zero-indexed) page
// request.
func ExamplePageRequest() {
	req := pagination.NewPageRequest(2, 20) // third page, 20 items per page
	fmt.Println(req.Offset(), req.Limit())
	// Output: 40 20
}

// ExampleNewPageResult wraps a page of content together with the total counts
// and the computed number of pages.
func ExampleNewPageResult() {
	req := pagination.NewPageRequest(0, 10)
	page := pagination.NewPageResult([]string{"a", "b", "c"}, 25, req)
	fmt.Println(page.TotalPages, page.IsFirst, page.HasNext)
	// Output: 3 true true
}
