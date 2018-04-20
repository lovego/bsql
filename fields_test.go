package bsql

import (
	"reflect"
	"testing"
)

func TestColumn2Field(t *testing.T) {
	var inputs = []string{"xiao_mei", "http_status", "you123"}
	var expects = []string{"XiaoMei", "HttpStatus", "You123"}
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
	var inputs = []string{"XiaoMei", "HTTPStatus", "You123"}
	var expects = []string{"xiao_mei", "http_status", "you123"}
	for i := range inputs {
		if got := Field2Column(inputs[i]); expects[i] != got {
			t.Errorf("expect: %v, got: %v", expects[i], got)
		}
	}
	if got := Fields2Columns(inputs); !reflect.DeepEqual(expects, got) {
		t.Errorf("expect: %v, got: %v", expects, got)
	}
}

type TestT struct {
	Name        string
	notExported int
	TestT2
	*TestT3
	TestT4
	testT5
}
type TestT2 struct {
	T2Name string
}
type TestT3 struct {
	T3Name string
}
type TestT4 int
type testT5 string

func TestFieldsFromStruct(t *testing.T) {
	got := FieldsFromStruct(TestT{}, []string{"T2Name"})
	expect := []string{"Name", "T3Name", "TestT4"}
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("unexpected: %v", got)
	}
}
