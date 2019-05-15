package bsql

import "fmt"

func ExampleOffsetToLineAndColumn() {
	fmt.Println(OffsetToLineAndColumn("", 0))
	fmt.Println(OffsetToLineAndColumn("", 1))

	fmt.Println(OffsetToLineAndColumn("a", 1))
	fmt.Println(OffsetToLineAndColumn("a\nb", 2))
	fmt.Println(OffsetToLineAndColumn("中文3\n\n\n789", 8))

	fmt.Println(OffsetToLineAndColumn("a\r\n", 3))
	fmt.Println(OffsetToLineAndColumn("a\r\nb", 4))

	// Output:
	// 0 0
	// 0 0
	// 1 1
	// 1 2
	// 4 2
	// 1 3
	// 2 1
}
