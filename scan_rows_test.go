package bsql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

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
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(row[i]))
	}
	return nil
}
func (s *testRows) Err() error {
	return nil
}
