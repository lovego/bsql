package bsql

import (
	"fmt"
	"reflect"
	"testing"
)

func TestQ(t *testing.T) {
	var data = []string{"bsql", "xiaomei", "lovego", "小明", "18300004444"}
	var expect = []string{"'bsql'", "'xiaomei'", "'lovego'", "'小明'", "'18300004444'"}
	for i := range data {
		if got := Q(data[i]); !reflect.DeepEqual(expect[i], got) {
			t.Errorf("expect %v,got %v", expect, got)
		}
	}
}

func TestColumn2Field(t *testing.T) {
	var expect = []string{"PsqlTest", "XiaoMei", "LoveGo", "You123"}
	var data = []string{"psql_test", "xiao_mei", "love_go", "you123"}
	for i := range data {
		if got := Column2Field(data[i]); !reflect.DeepEqual(expect[i], got) {
			t.Errorf("expect %v,got %v", expect, got)
		}

	}
}

func TestColumn2Fields(t *testing.T) {
	var expect = []string{"PsqlTest", "UserName", "Phone"}
	if got := Columns2Fields([]string{"psql_test", "user_name", "phone"}); !reflect.DeepEqual(expect, got) {
		t.Errorf("expect %v,got %v", expect, got)
	}
}

func TestStructFieldsAddrs(t *testing.T) {
	var psqlTest struct {
		Id       int64
		Name     string
		IsExists bool
	}
	psqlTest.Id = int64(1)
	psqlTest.Name = "psql"
	psqlTest.IsExists = true
	// var expect = []interface{}{psqlTest.Id, psqlTest.Name, psqlTest.IsExists}
	// if got, err := StructFieldsAddrs(reflect.ValueOf(psqlTest),
	// 	[]string{"Id", "Name", "IsExists"}); err != nil || !reflect.DeepEqual(expect, got) {
	// 	t.Errorf("expect %v,got %v,%t", expect, got, err)
	// }
	fmt.Println(StructFieldsAddrs(reflect.ValueOf(&psqlTest),
		[]string{"Id", "Name", "IsExists"}))
}
