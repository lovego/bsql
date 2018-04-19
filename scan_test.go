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
			{int64(1), "李雷", "男", testTime}, {int64(2), "韩梅梅", "女", testTime},
			{int64(3), "Lili", "女", testTime}, {int64(4), "Lucy", "女", testTime},
			{int64(5), "Mr Gao", "男", testTime}, {int64(6), "Uncle Wang", "男", testTime},
		},
		i: -1,
	}
}
