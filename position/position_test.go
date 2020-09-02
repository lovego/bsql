package position

import "fmt"

func ExampleGet() {
	fmt.Println(Get([]rune("a"), 0))
	fmt.Println(Get([]rune("a\nb"), 1))
	fmt.Println(Get([]rune("中文2\n\n\n6789"), 7))

	// Output:
	// Line 1: a
	// Char 1: ^
	// Line 1: a
	// Char 2:  ^
	// Line 4: 6789
	// Char 2:  ^
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
