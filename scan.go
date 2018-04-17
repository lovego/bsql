package bsql

import (
	"database/sql"
	"errors"
	"reflect"
)

func Scan(rows *sql.Rows, data interface{}) error {
	p := reflect.ValueOf(data)
	if p.Kind() != reflect.Ptr {
		return errors.New("data must be a pointer.")
	}
	target := p.Elem()
	switch target.Kind() {
	case reflect.Struct:
		if err := ScanStruct(rows, target); err != nil {
			return err
		}
	case reflect.Slice:
		if err := ScanSlice(rows, target, p); err != nil {
			return err
		}
	default:
		if rows.Next() {
			if err := rows.Scan(p); err != nil {
				return err
			}
		}
	}
	return rows.Err()
}

func ScanStruct(rows *sql.Rows, target reflect.Value) error {
	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}
	fieldsAddrs, err := StructFieldsAddrs(target, Columns2Fields(columnNames))
	if err != nil {
		return err
	}
	if rows.Next() {
		if err := rows.Scan(fieldsAddrs...); err != nil {
			return err
		}
	}
	return nil
}

func ScanSlice(rows *sql.Rows, target, p reflect.Value) error {
	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}
	fieldNames := Columns2Fields(columnNames)
	elemType := target.Type().Elem()
	for rows.Next() {
		elemValue := reflect.Zero(elemType)
		if fieldsAddrs, err := StructFieldsAddrs(elemValue, fieldNames); err != nil {
			return err
		} else if err := rows.Scan(fieldsAddrs...); err != nil {
			return err
		}
		target = reflect.Append(target, elemValue)
	}
	p.Elem().Set(target)
	return nil
}
