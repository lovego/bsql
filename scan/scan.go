package scan

import (
	"database/sql"
	"errors"
	"reflect"
	"time"
)

// data must be a sql.Scanner or a non nil pointer.
// If data a pointer, it's indirect once to get the target to store rows from database.
// And it's indirect only once, we cann't indirect all pointer until non pointer type,
// because nil value should be set to the second layer pointer, not the non pointer type.
// If target is a slice, it scan all rows into the slice, otherwise it scan a single row.
func Scan(rows *sql.Rows, data interface{}) error {
	if _, ok := data.(sql.Scanner); ok {
		if rows.Next() {
			if err := rows.Scan(data); err != nil {
				return err
			}
		}
		return rows.Err()
	}

	ptr := reflect.ValueOf(data)
	if ptr.Kind() != reflect.Ptr {
		return errors.New("bsql: data must be a pointer.")
	}
	if ptr.IsNil() {
		return errors.New("bsql: data is a nil pointer.")
	}
	columns, err := getColumns(rows)
	if err != nil {
		return err
	}
	if len(columns) == 0 {
		return errors.New("bsql: no columns.")
	}

	target := ptr.Elem()
	switch target.Kind() {
	case reflect.Slice:
		typ := target.Type().Elem()
		for rows.Next() {
			elem := reflect.New(typ).Elem()
			if err := scanSingleRow(rows, columns, elem); err != nil {
				return err
			}
			target.Set(reflect.Append(target, elem))
		}
	default:
		if rows.Next() {
			if err := scanSingleRow(rows, columns, target); err != nil {
				return err
			}
		}
	}
	return rows.Err()
}

// If target is a struct, it scan all columns into the struct, otherwise it scan a single column.
// No indirect is performed, because nil value should be set to pointer.
func scanSingleRow(rows *sql.Rows, columns []columnType, target reflect.Value) error {
	addr := target.Addr().Interface()
	switch addr.(type) {
	case sql.Scanner:
		return rows.Scan(addr)
	case *time.Time:
		return rows.Scan(scannerOf(target, columns[0]))
	}

	switch target.Kind() {
	case reflect.Struct:
		if err := scan2Struct(rows, columns, target); err != nil {
			return err
		}
	default:
		if err := rows.Scan(scannerOf(target, columns[0])); err != nil {
			return err
		}
	}
	return nil
}

func scan2Struct(rows *sql.Rows, columns []columnType, target reflect.Value) error {
	var scanners []interface{}
	for _, column := range columns {
		field := FieldByName(target, column.FieldName)
		if !field.IsValid() {
			return errors.New("bsql: no or multiple field '" + column.FieldName + "' in struct")
		}
		scanners = append(scanners, scannerOf(field, column))
	}
	if err := rows.Scan(scanners...); err != nil {
		return err
	}
	return nil
}
