package bsql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	for rows.Next() {
		targets = reflect.Append(targets, reflect.Zero(elemType))
		target := targets.Index(targets.Len() - 1)
		if elemType.Kind() == reflect.Struct {
			if err := scan2Struct(rows, columns, target); err != nil {
				return err
			}
		} else if err := rows.Scan(scannerOf(target.Addr(), columns[0])); err != nil {
			return err
		}
	}
	p.Elem().Set(targets)
	return nil
}

func scan2Struct(rows rowsType, columns []columnType, target reflect.Value) error {
	scanners, err := structFieldsScanners(target, columns)
	if err != nil {
		return err
	}
	if err := rows.Scan(scanners...); err != nil {
		return err
	}
	return nil
}

func structFieldsScanners(structValue reflect.Value, columns []columnType) ([]interface{}, error) {
	var result []interface{}
	for _, column := range columns {
		field := structValue.FieldByName(column.FieldName)
		if !field.IsValid() {
			return nil, errors.New("no field: '" + column.FieldName + "' in struct.")
		}
		result = append(result, scannerOf(field.Addr(), column))
	}
	return result, nil
}

func scannerOf(addr reflect.Value, column columnType) interface{} {
	var dbType string
	if column.ColumnType != nil {
		dbType = column.ColumnType.DatabaseTypeName()
	}
	switch dbType {
	case "JSONB", "JSON":
		return jsonScanner{addr.Interface()}
	// case "ARRAY":
	default:
		return addr.Interface()
	}
}

type jsonScanner struct {
	data interface{}
}

func (s jsonScanner) Scan(src interface{}) error {
	switch buf := src.(type) {
	case string:
		return json.Unmarshal([]byte(buf), s.data)
	case []byte:
		return json.Unmarshal(buf, s.data)
	default:
		return fmt.Errorf("bsql jsonScanner unexpected: %T %v", src, src)
	}
}
