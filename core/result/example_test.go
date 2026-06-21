package result_test

import (
	"fmt"

	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
)

func divide(a, b int) result.Result[float64] {
	if b == 0 {
		return result.Err[float64](kerrors.NewKError("division by zero", 400, nil))
	}
	return result.Ok(float64(a) / float64(b))
}

// ExampleResult shows returning and inspecting a Result instead of the
// cascading (value, error) pattern.
func ExampleResult() {
	ok := divide(10, 2)
	fmt.Println(ok.IsOk(), ok.Value())

	bad := divide(1, 0)
	fmt.Println(bad.IsOk(), bad.Error())
	// Output:
	// true 5
	// false division by zero
}
