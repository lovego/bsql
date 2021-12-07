package bsql

import (
	"fmt"
)

func ExamplePrettyPrint() {
	var v = struct {
		A int
		B string
	}{A: 1, B: "line1\nline2"}

	fmt.Print(PrettyPrint(v))
	// Output:
	// {
	//   A: 1
	//   B: line1
	// line2
	// }
}
