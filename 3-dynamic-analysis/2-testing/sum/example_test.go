package sum_test

import (
	"fmt"

	"github.com/campoy/go-tooling-workshop/3-dynamic-analysis/2-testing/sum"
)

func ExampleSum() {
	fmt.Println(sum.All(1, 2, 3, 4, 5))
	// Output:
	// 15
}
