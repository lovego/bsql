package bsql

import (
	"fmt"
)

func ExampleField2Column() {
	var inputs = []string{"XiaoMei", "HTTPStatus", "You123", "ILoveGolangAndJSONSoMuch"}
	for i := range inputs {
		fmt.Println(Field2Column(inputs[i]))
	}
	// Output:
	// xiao_mei
	// http_status
	// you123
	// i_love_golang_and_json_so_much
}

func ExampleFields2Columns() {
	var inputs = []string{"XiaoMei", "HTTPStatus", "You123",
		"PriceP", "4sPrice", "Price4s", "goodHTTP", "ILoveGolangAndJSONSoMuch",
	}
	fmt.Println(Fields2Columns(inputs))
	// Output:
	// [xiao_mei http_status you123 price_p 4s_price price4s good_http i_love_golang_and_json_so_much]
}

func ExampleFields2ColumnsStr() {
	var inputs = []string{"XiaoMei", "HttpStatus", "You123", "PriceP"}
	fmt.Println(Fields2ColumnsStr(inputs))
	// Output:
	// xiao_mei,http_status,you123,price_p
}

func ExampleFieldsToColumnsStr() {
	var inputs = []string{"XiaoMei", "HttpStatus", "You123", "PriceP"}
	fmt.Println(FieldsToColumnsStr(inputs, "t.", []string{"PriceP"}))
	// Output:
	// t.xiao_mei,t.http_status,t.you123
}

func ExampleFieldsFromStruct() {
	type TestT2 struct {
		T2Name string
	}
	type testT3 struct {
		T3Name string
	}
	type TestT4 int
	type testT5 string

	type TestT struct {
		Name        string
		notExported int
		TestT2
		*testT3
		TestT4
		testT5
	}
	fmt.Println(FieldsFromStruct(TestT{}, []string{"T2Name"}))
	// Output:
	// [Name T3Name TestT4]
}
