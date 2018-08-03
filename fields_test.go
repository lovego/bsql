package bsql

import (
	"fmt"
)

func ExampleColumn2Field() {
	var inputs = []string{"xiao_mei", "http_status", "you123", "price_p"}
	for i := range inputs {
		fmt.Println(Column2Field(inputs[i]))
	}
	fmt.Println(Columns2Fields(inputs))
	// Output:
	// XiaoMei
	// HttpStatus
	// You123
	// PriceP
	// [XiaoMei HttpStatus You123 PriceP]
}

func ExampleField2Column() {
	var inputs = []string{"XiaoMei", "HTTPStatus", "You123",
		"PriceP", "4sPrice", "Price4s", "goodHTTP", "ILoveGolangAndJSONSoMuch",
	}
	fmt.Println(Fields2Columns(inputs))
	// Output:
	// [xiao_mei http_status you123 price_p 4s_price price4s good_http i_love_golang_and_json_so_much]
}

func ExampleFieldsFromStruct() {
	type TestT2 struct {
		T2Name string
	}
	type TestT3 struct {
		T3Name string
	}
	type TestT4 int
	type testT5 string

	type TestT struct {
		Name        string
		notExported int
		TestT2
		*TestT3
		TestT4
		testT5
	}
	fmt.Println(FieldsFromStruct(TestT{}, []string{"T2Name"}))
	// Output:
	// [Name T3Name TestT4]
}

func ExampleColumnsComments() {
	type Test struct {
		Id          int64  `comment:"主键"`
		Name        string `comment:"名称"`
		notExported int
	}

	fmt.Println(ColumnsComments("tests", Test{}))
	// OutPut:
	// comment on column tests.id is '主键';
	// comment on column tests.name is '名称';
}
