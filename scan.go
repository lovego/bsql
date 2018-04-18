package bsql

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type Rows interface {
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}

func Scan(rows Rows, data interface{}) error {
	p := reflect.ValueOf(data)
	if p.Kind() != reflect.Ptr {
		return errors.New("data must be a pointer.")
	}
	target := p.Elem()
	switch target.Kind() {
	case reflect.Slice, reflect.Array:
		if err := Scan2Slice(rows, target, p); err != nil {
			return err
		}
	case reflect.Struct:
		if err := Scan2Struct(rows, target); err != nil {
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

func Scan2Slice(rows Rows, target, p reflect.Value) error {
	elemType := target.Type().Elem()
	if elemType.Kind() == reflect.Struct {
		return Scan2StructSlice(rows, target, p)
	}
	for rows.Next() {
		target = reflect.Append(target, reflect.Zero(elemType))
		if err := rows.Scan(target.Index(target.Len() - 1).Addr().Interface()); err != nil {
			return err
		}
	}
	p.Elem().Set(target)
	return nil
}
func Scan2StructSlice(rows Rows, target, p reflect.Value) error {
	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}
	fieldNames := Columns2Fields(columnNames)
	elemType := target.Type().Elem()
	for rows.Next() {
		target = reflect.Append(target, reflect.Zero(elemType))
		fieldsAddrs, err := StructFieldsScanners(target.Index(target.Len()-1), fieldNames)
		if err != nil {
			return err
		}
		if err := rows.Scan(fieldsAddrs...); err != nil {
			return err
		}
	}
	p.Elem().Set(target)
	return nil
}

func Scan2Struct(rows Rows, target reflect.Value) error {
	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}
	fieldsAddrs, err := StructFieldsScanners(target, Columns2Fields(columnNames))
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

func StructFieldsScanners(structValue reflect.Value, fieldNames []string) ([]interface{}, error) {
	var result []interface{}
	for _, fieldName := range fieldNames {
		field := structValue.FieldByName(fieldName)
		if !field.IsValid() {
			return nil, errors.New("no field: '" + fieldName + "' in struct.")
		}
		result = append(result, scannerOf(field))
	}
	return result, nil
}

func scannerOf(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
		return jsonScanner{value.Addr().Interface()}
	default:
		return value.Addr().Interface()
	}
}

type jsonScanner struct {
	p interface{}
}

func (js jsonScanner) Scan(src interface{}) error {
	switch buf := src.(type) {
	case string:
		return json.Unmarshal([]byte(buf), js.p)
	case []byte:
		return json.Unmarshal(buf, js.p)
	default:
		return fmt.Errorf("bsql unexpected: %T %v", src, src)
	}
}
