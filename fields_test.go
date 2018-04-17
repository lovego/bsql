package bsql

import (
	"reflect"
	"testing"
)

func TestColumn2Field(t *testing.T) {
	var data = []string{"bsql_test", "xiao_mei", "love_go", "you123"}
	var expect = []string{"PsqlTest", "XiaoMei", "LoveGo", "You123"}
	for i := range data {
		if got := Column2Field(data[i]); !reflect.DeepEqual(expect[i], got) {
			t.Errorf("expect %v,got %v", expect, got)
		}
	}
}

func TestColumns2Fields(t *testing.T) {
	var expect = []string{"PsqlTest", "UserName", "Phone"}
	if got := Columns2Fields([]string{"bsql_test", "user_name", "phone"}); !reflect.DeepEqual(expect, got) {
		t.Errorf("expect %v,got %v", expect, got)
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

func TestFieldsFromStruct() {
	got := FieldsFromStruct(T{}, []string{"T2Name"})
	expect := []string{"name", "T3Name"}
	if !refelct.DeepEqual(got, expect) {
		t.Fatalf("unexpected: %v", got)
	}
}
