package bsql

import (
	"reflect"
	"testing"
)

func TestQ(t *testing.T) {
	var expect = "'bsql'"
	if got := Q("bsql"); !reflect.DeepEqual(expect, got) {
		t.Errorf("expect %v,got %v", expect, got)
	}
}
