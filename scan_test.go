package bsql

import (
	"reflect"
	"testing"
	"time"
)

type testUser struct {
	Id            int
	UserName, Sex string
	CreatedAt     time.Time
}

var testTime = time.Now()

func TestScan2StructSlice(t *testing.T) {
	var got []testUser
	if err := scan(testUsers(), &got); err != nil {
		t.Fatal(err)
	}
	expect := []testUser{
		{1, "李雷", "男", testTime}, {2, "韩梅梅", "女", testTime},
		{3, "Lili", "女", testTime}, {4, "Lucy", "女", testTime},
		{5, "Mr Gao", "男", testTime}, {6, "Uncle Wang", "男", testTime},
	}
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("unexpected: %+v", got)
	}
	t.Logf("%+v", got)
}

func TestScan2Slice(t *testing.T) {
	var got []int
	if err := scan(testUsers(), &got); err != nil {
		t.Fatal(err)
	}
	expect := []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("unexpected: %+v", got)
	}
	t.Logf("%+v", got)
}

func TestScan2Struct(t *testing.T) {
	var got = testUser{}
	if err := scan(testUsers(), &got); err != nil {
		t.Fatal(err)
	}
	expect := testUser{1, "李雷", "男", testTime}
	if got != expect {
		t.Fatalf("unexpected: %+v", got)
	}
	t.Logf("%+v", got)
}

func TestScan2Int(t *testing.T) {
	var got int
	if err := scan(testUsers(), &got); err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Fatalf("unexpected: %+v", got)
	}
	t.Logf("%+v", got)
}

func testUsers() *testRows {
	return &testRows{
		columns: []string{"id", "user_name", "sex", "created_at"},
		rows: [][]interface{}{
			{1, "李雷", "男", testTime}, {2, "韩梅梅", "女", testTime},
			{3, "Lili", "女", testTime}, {4, "Lucy", "女", testTime},
			{5, "Mr Gao", "男", testTime}, {6, "Uncle Wang", "男", testTime},
		},
		i: -1,
	}
}

func TestStructFieldsScanners(t *testing.T) {
	var v struct {
		Id     int64
		Name   string
		Exists bool
	}
	addrs, err := structFieldsScanners(reflect.ValueOf(&v).Elem(), []columnType{
		{FieldName: "Id"}, {FieldName: "Name"}, {FieldName: "Exists"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("unexpected addrs size: %d", len(addrs))
	}
	if p, ok := addrs[0].(*int64); !ok {
		t.Errorf("unexpected type: %T", addrs[0])
	} else if p != &v.Id {
		t.Errorf("unexpected addr: %p, expect: %p", p, &v.Id)
	}
	if p, ok := addrs[1].(*string); !ok {
		t.Errorf("unexpected type: %T", addrs[1])
	} else if p != &v.Name {
		t.Errorf("unexpected addr: %p, expect: %p", p, &v.Name)
	}
	if p, ok := addrs[2].(*bool); !ok {
		t.Errorf("unexpected type: %T", addrs[2])
	} else if p != &v.Exists {
		t.Errorf("unexpected addr: %p, expect: %p", p, &v.Exists)
	}
}
