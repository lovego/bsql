package bsql

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestQ(t *testing.T) {
	var data = []string{"bsql", "xi'ao'mei", "love'go"}
	var expect = []string{"'bsql'", "'xi''ao''mei'", "'love''go'"}
	for i := range data {
		if got := Q(data[i]); !reflect.DeepEqual(expect[i], got) {
			t.Errorf("expect: %v, got: %v", expect, got)
		}
	}
}

func TestValues(t *testing.T) {
	if got := Values(3); got != "(3)" {
		t.Errorf("unexpected: %s", got)
	}

	if got := Values("a'bc"); got != "('a''bc')" {
		t.Errorf("unexpected: %s", got)
	}

	if got := Values([]int{1, 2, 3}); got != "(1,2,3)" {
		t.Errorf("unexpected: %s", got)
	}

	if got := Values([]string{"a", "b", "c"}); got != "('a','b','c')" {
		t.Errorf("unexpected: %s", got)
	}

	if got := Values([][]interface{}{
		{1, "a", true}, {2, "b", true}, {3, "c", false},
	}); got != "(1,'a',t),(2,'b',t),(3,'c',f)" {
		t.Errorf("unexpected: %s", got)
	}
}

func TestStructValues(t *testing.T) {
	data := []struct {
		Id        int
		Name, Sex string
	}{
		{1, "李雷", "男"}, {2, "韩梅梅", "女"},
		{3, "Lili", "女"}, {4, "Lucy", "女"},
	}
	got := StructValues(data, []string{"Id", "Name"})
	if got != "(1,'李雷'),(2,'韩梅梅'),(3,'Lili'),(4,'Lucy')" {
		t.Errorf("unexpected: %s", got)
	}
}

func TestV(t *testing.T) {
	if got := V(decimal.New(1230, -2)); got != "12.3" {
		t.Errorf("unexpected: %s", got)
	}

	if got := V(map[int]bool{2: true, 3: false}); got != `'{"2":true,"3":false}'` {
		t.Errorf("unexpected: %s", got)
	}
}
