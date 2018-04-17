package bsql

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type testScanner struct {
	columns []string
	rows    [][]interface{}
	i       int
}

func newTestScanner(columns []string, rows [][]interface{}) *testScanner {
	return &testScanner{columns: columns, rows: rows, i: -1}
}

func (s *testScanner) Columns() ([]string, error) {
	return s.columns, nil
}
func (s *testScanner) Next() bool {
	if s.i < 0 {
		s.i = 0
	} else {
		s.i++
	}
	return s.i < len(s.rows)
}
func (s *testScanner) Scan(dests ...interface{}) error {
	if s.i >= len(s.rows) {
		return errors.New("all data has been scanned.")
	}
	row := s.rows[s.i]
	if len(dests) > len(row) {
		return fmt.Errorf("sql: expected most %d destination arguments in Scan, got %d", len(row), len(dests))
	}
	for i, desc := range dests {
		reflect.ValueOf(desc).Elem().Set(reflect.ValueOf(row[i]))
	}
	return nil
}
func (s *testScanner) Err() error {
	return nil
}

type testUser struct {
	Id            int
	UserName, Sex string
	CreatedAt     time.Time
}

var testTime = time.Now()
var testUsers = newTestScanner([]string{"id", "user_name", "sex", "created_at"}, [][]interface{}{
	{1, "李雷", "男", testTime}, {2, "韩梅梅", "女", testTime},
	{3, "Lili", "女", testTime}, {4, "Lucy", "女", testTime},
	{5, "Mr Gao", "男", testTime}, {6, "Uncle Wang", "男", testTime},
})

func TestScanStruct(t *testing.T) {
	var got = testUser{}
	if err := ScanStruct(testUsers, reflect.ValueOf(&got).Elem()); err != nil {
		t.Fatal(err)
	}
	expect := testUser{1, "李雷", "男", testTime}
	if got != expect {
		t.Fatalf("unexpected: %+v", got)
	}
	t.Logf("%+v", got)
}
