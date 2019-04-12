package bsql

import "fmt"

func ExampleOffsetToLineAndColumn() {
	fmt.Println(OffsetToLineAndColumn("", 0))
	fmt.Println(OffsetToLineAndColumn("", 1))

	fmt.Println(OffsetToLineAndColumn("a", 1))
	fmt.Println(OffsetToLineAndColumn("a\n", 2))
	fmt.Println(OffsetToLineAndColumn("a\r\n", 3))
	fmt.Println(OffsetToLineAndColumn("a\r\nb", 4))

	// Output:
	// 0 0
	// 0 0
	// 1 1
	// 1 2
	// 1 3
	// 2 1
}
