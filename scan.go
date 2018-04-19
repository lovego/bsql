package bsql

import (
	"database/sql"
	"errors"
	"reflect"
)

type rowsType interface {
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}

func scan(rows rowsType, data interface{}) error {
	if _, ok := data.(sql.Scanner); ok {
		if rows.Next() {
			if err := rows.Scan(data); err != nil {
				return err
			}
		}
		return rows.Err()
	}

	p := reflect.ValueOf(data)
	if p.Kind() != reflect.Ptr {
		return errors.New("data must be a pointer.")
	}
	columns, err := getColumns(rows)
	if err != nil {
		return err
	}
	if len(columns) == 0 {
		return errors.New("no columns.")
	}
	target := p.Elem()
	switch target.Kind() {
	case reflect.Slice, reflect.Array:
		if err := scan2Slice(rows, columns, target, p); err != nil {
			return err
		}
	case reflect.Struct:
		if rows.Next() {
			if err := scan2Struct(rows, columns, target); err != nil {
				return err
			}
		}
	default:
		if rows.Next() {
			if err := rows.Scan(scannerOf(p, columns[0])); err != nil {
				return err
			}
		}
	}
	return rows.Err()
}

func scan2Slice(rows rowsType, columns []columnType, targets, p reflect.Value) error {
	elemType := targets.Type().Elem()
	var isPtr bool
	if elemType.Kind() == reflect.Ptr {
		elemType, isPtr = elemType.Elem(), true
	}
	for rows.Next() {
		ts := reflect.MakeSlice(reflect.SliceOf(elemType), 1, 1)
		ts.Index(0).Set(reflect.Zero(elemType))
		target := ts.Index(0)
		if elemType.Kind() == reflect.Struct {
			if err := scan2Struct(rows, columns, target); err != nil {
				return err
			}
		} else if err := rows.Scan(scannerOf(target.Addr(), columns[0])); err != nil {
			return err
		}
		if isPtr {
			target = target.Addr()
		}
		targets = reflect.Append(targets, target)
	}
	p.Elem().Set(targets)
	return nil
}

func scan2Struct(rows rowsType, columns []columnType, target reflect.Value) error {
	var scanners []interface{}
	for _, column := range columns {
		field := target.FieldByName(column.FieldName)
		if !field.IsValid() {
			return errors.New("no field: '" + column.FieldName + "' in struct.")
		}
		scanners = append(scanners, scannerOf(field.Addr(), column))
	}
	if err := rows.Scan(scanners...); err != nil {
		return err
	}
	return nil
}
