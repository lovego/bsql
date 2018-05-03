package bsql

import (
	"reflect"
	"testing"
)

func TestColumn2Field(t *testing.T) {
	var inputs = []string{"xiao_mei", "http_status", "you123", "price_p"}
	var expects = []string{"XiaoMei", "HttpStatus", "You123", "PriceP"}
	for i := range inputs {
		if got := Column2Field(inputs[i]); !reflect.DeepEqual(expects[i], got) {
			t.Errorf("expect: %v, got: %v", expects, got)
		}
	}
	if got := Columns2Fields(inputs); !reflect.DeepEqual(expects, got) {
		t.Errorf("expect: %v, got: %v", expects, got)
	}
}

func TestField2Column(t *testing.T) {
	var inputs = []string{"XiaoMei", "HTTPStatus", "You123",
		"PriceP", "4sPrice", "Price4s", "goodHTTP", "ILoveGolangAndJSONSoMuch",
	}
	var expects = []string{"xiao_mei", "http_status", "you123",
		"price_p", "4s_price", "price4s", "good_http", "i_love_golang_and_json_so_much",
	}
	for i := range inputs {
		if got := Field2Column(inputs[i]); expects[i] != got {
			t.Errorf("expect: %v, got: %v", expects[i], got)
		}
	}
	if got := Fields2Columns(inputs); !reflect.DeepEqual(expects, got) {
		t.Errorf("expect: %v, got: %v", expects, got)
	} else {
		t.Log(got)
	}
}

func TestFieldsFromStruct(t *testing.T) {
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

	got := FieldsFromStruct(TestT{}, []string{"T2Name"})
	expect := []string{"Name", "T3Name", "TestT4"}
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestColumnsComments(t *testing.T) {
	type Test struct {
		Id          int64  `comment:"主键"`
		Name        string `comment:"名称"`
		notExported int
	}

	got := ColumnsComments("tests", Test{})
	expect := `comment on column tests.id is '主键';
comment on column tests.name is '名称';
`
	if got != expect {
		t.Fatalf("unexpected: %v", got)
	}
}
