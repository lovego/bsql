package bsql

import (
	"reflect"
	"testing"
)

func TestQ(t *testing.T) {
	var data = []string{"b'sql", "xiaomei", "lovego", "小明", "18300004444"}
	var expect = []string{"'b''sql'", "'xiaomei'", "'lovego'", "'小明'", "'18300004444'"}
	for i := range data {
		if got := Q(data[i]); !reflect.DeepEqual(expect[i], got) {
			t.Errorf("expect %v,got %v", expect, got)
		}
	}
}

func TestColumn2Field(t *testing.T) {
	var data = []string{"bsql_test", "xiao_mei", "love_go", "you123"}
	var expect = []string{"BsqlTest", "XiaoMei", "LoveGo", "You123"}
	for i := range data {
		if got := Column2Field(data[i]); !reflect.DeepEqual(expect[i], got) {
			t.Errorf("expect %v,got %v", expect, got)
		}
	}
}

func TestColumn2Fields(t *testing.T) {
	var expect = []string{"BsqlTest", "UserName", "Phone"}
	if got := Columns2Fields([]string{"bsql_test", "user_name", "phone"}); !reflect.DeepEqual(expect, got) {
		t.Errorf("expect %v,got %v", expect, got)
	}
}

func TestStructFieldsAddrs(t *testing.T) {
	var bsqlTest struct {
		Id     int64
		Name   string
		Exists bool
	}
	var addrs []interface{}
	if addrs, err := StructFieldsAddrs(reflect.ValueOf(&bsqlTest).Elem(),
		[]string{"Id", "Name", "Exists"}); err != nil || len(addrs) != 3 {
		t.Fatalf("unexpted addrs size: %d", len(addrs))
	}
	for i := range addrs {
		if _, ok := addrs[i].(*int64); ok {
			continue
		}
		if _, ok := addrs[i].(*string); ok {
			continue
		}
		if _, ok := addrs[i].(*bool); ok {
			continue
		}
		t.Fatalf("指针类型错误", addrs[i])
	}
}
