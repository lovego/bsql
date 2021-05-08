package bsql

import "fmt"

func ExampleColumns2Fields() {
	var inputs = []string{"xiao_mei", "http_status", "you123", "price_p"}
	fmt.Println(Columns2Fields(inputs))
	// Output:
	// [XiaoMei HttpStatus You123 PriceP]
}

func ExampleColumnsComments() {
	type Test struct {
		Id          int64  `comment:"主键"`
		Name        string `comment:"名称"`
		notExported int
	}

	fmt.Println(ColumnsComments("tests", Test{}))
	// OutPut:
	// COMMENT ON COLUMN tests.id IS '主键';
	// COMMENT ON COLUMN tests.name IS '名称';
}
