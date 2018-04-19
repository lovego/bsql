package bsql

import (
	"database/sql"
	"errors"
	"fmt"
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

	var got1 []*testUser
	if err := scan(testUsers(), &got1); err != nil {
		t.Fatal(err)
	}
	expect1 := []*testUser{
		&testUser{1, "李雷", "男", testTime}, &testUser{2, "韩梅梅", "女", testTime},
		&testUser{3, "Lili", "女", testTime}, &testUser{4, "Lucy", "女", testTime},
		&testUser{5, "Mr Gao", "男", testTime}, &testUser{6, "Uncle Wang", "男", testTime},
	}
	if !reflect.DeepEqual(got1, expect1) {
		t.Fatalf("unexpected: %+v", got1)
	}
	t.Logf("%+v", got1)
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

type testRows struct {
	columns []string
	rows    [][]interface{}
	i       int
}

func (s *testRows) ColumnTypes() ([]*sql.ColumnType, error) {
	return nil, nil
}

func (s *testRows) Columns() ([]string, error) {
	return s.columns, nil
}
func (s *testRows) Next() bool {
	if s.i < 0 {
		s.i = 0
	} else {
		s.i++
	}
	return s.i < len(s.rows)
}
func (s *testRows) Scan(dests ...interface{}) error {
	if s.i >= len(s.rows) {
		return errors.New("all data has been scanned.")
	}
	row := s.rows[s.i]
	if len(dests) > len(row) {
		return fmt.Errorf("sql: expected most %d destination arguments in Scan, got %d", len(row), len(dests))
	}
	for i, dest := range dests {
		if scanner, ok := dest.(sql.Scanner); ok {
			if err := scanner.Scan(row[i]); err != nil {
				return err
			}
		} else {
			reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(row[i]))
		}
	}
	return nil
}

func (s *testRows) Err() error {
	return nil
}
