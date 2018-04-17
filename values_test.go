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
