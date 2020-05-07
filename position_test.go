package bsql

import "fmt"

func ExampleGetPosition() {
	fmt.Println(GetPosition([]rune("a"), 0))
	fmt.Println(GetPosition([]rune("a\nb"), 1))
	fmt.Println(GetPosition([]rune("中文3\n\n\n6789"), 7))

	// Output:
	// Line 1: a
	//         ^
	// Line 1: a
	//          ^
	// Line 4: 6789
	//          ^
}

func ExampleOffsetToLineAndColumn() {
	fmt.Println(OffsetToLineAndColumn([]rune(""), -1))
	fmt.Println(OffsetToLineAndColumn([]rune(""), 0))

	fmt.Println(OffsetToLineAndColumn([]rune("a"), 0))
	fmt.Println(OffsetToLineAndColumn([]rune("a\nb"), 1))
	fmt.Println(OffsetToLineAndColumn([]rune("中文3\n\n\n6789"), 7))

	fmt.Println(OffsetToLineAndColumn([]rune("a\r\n"), 2))
	fmt.Println(OffsetToLineAndColumn([]rune("a\r\nb"), 3))

	// Output:
	// 0 0
	// 0 0
	// 1 1
	// 1 2
	// 4 2
	// 1 3
	// 2 1
}
