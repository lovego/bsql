package scan

import "fmt"

func ExampleColumn2Field() {
	var inputs = []string{"xiao_mei", "http_status", "you123", "price_p"}
	for i := range inputs {
		fmt.Println(Column2Field(inputs[i]))
	}
	// Output:
	// XiaoMei
	// HttpStatus
	// You123
	// PriceP
}
