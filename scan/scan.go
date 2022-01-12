package scan

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/lovego/value"
)

// data must be a sql.Scanner or a non nil pointer.
// If data a pointer, it's indirect once to get the target to store rows from database.
// And it's indirect only once, we cann't indirect all pointer until non pointer type,
// because nil value should be set to the second layer pointer, not the non pointer type.
// If target is a slice, it scan all rows into the slice, otherwise it scan a single row.
func Scan(rows *sql.Rows, data interface{}) error {
	if scanner := trySqlScanner(data); scanner != nil {
		if rows.Next() {
			if err := rows.Scan(scanner); err != nil {
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
	columns, err := ColumnTypes(rows)
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
			if err := ScanRow(rows, columns, elem); err != nil {
				return err
			}
			target.Set(reflect.Append(target, elem))
		}
	default:
		if rows.Next() {
			if err := ScanRow(rows, columns, target); err != nil {
				return err
			}
		}
	}
	return rows.Err()
}

// If target is a struct, it scan all columns into the struct, otherwise it scan a single column.
// No indirect is performed, because nil value should be set to pointer.
func ScanRow(rows *sql.Rows, columns []ColumnType, target reflect.Value) error {
	addr := target.Addr().Interface()
	if scanner := trySqlScanner(addr); scanner != nil {
		return rows.Scan(scanner)
	}
	switch addr.(type) {
	case *time.Time:
		return rows.Scan(scannerOf(target, columns[0]))
	}

	switch target.Kind() {
	case reflect.Struct:
		if err := scan2Struct(rows, columns, target); err != nil {
			return err
		}
	case reflect.Map:
		if err := scan2Map(rows, columns, target); err != nil {
			return err
		}
	default:
		if err := rows.Scan(scannerOf(target, columns[0])); err != nil {
			return err
		}
	}
	return nil
}

func scan2Struct(rows *sql.Rows, columns []ColumnType, target reflect.Value) error {
	var scanners []interface{}
	for _, column := range columns {
		field := value.Settable(target, column.FieldPath)
		if !field.IsValid() {
			return errors.New("bsql: no or multiple field '" +
				strings.Join(column.FieldPath, ".") + "' in struct")
		}
		scanners = append(scanners, scannerOf(field, column))
	}
	if err := rows.Scan(scanners...); err != nil {
		return err
	}
	return nil
}

func scan2Map(rows *sql.Rows, columns []ColumnType, target reflect.Value) error {
	if target.IsNil() {
		target.Set(reflect.MakeMap(target.Type()))
	}

	var scanners []interface{}
	for _, column := range columns {
		scanners = append(scanners, &mapFieldScanner{target, column.FieldName})
	}
	if err := rows.Scan(scanners...); err != nil {
		return err
	}
	return nil
}

type mapFieldScanner struct {
	m     reflect.Value
	field string
}

func (mfs *mapFieldScanner) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		src = string(v)
	}

	mfs.m.SetMapIndex(reflect.ValueOf(mfs.field), reflect.ValueOf(src))
	return nil
}
