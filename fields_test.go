package bsql

import (
	"reflect"
	"strings"
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
	expect := strings.Join(expects, ",")
	if got := Fields2Columns(inputs); !reflect.DeepEqual(expect, got) {
		t.Errorf("expect: %v, got: %v", expect, got)
	}
}

type T struct {
	Name string
	T2
}
type T2 struct {
	T2Name string
	T3
}
type T3 struct {
	T3Name string
}

func TestFieldsFromStruct(t *testing.T) {
	got := FieldsFromStruct(T{}, []string{"T2Name"})
	expect := []string{"Name", "T3Name"}
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("unexpected: %v", got)
	}
}
