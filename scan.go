package bsql

import (
	"errors"
	"reflect"
)

type Scanner interface {
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}

func Scan(scanner Scanner, data interface{}) error {
	p := reflect.ValueOf(data)
	if p.Kind() != reflect.Ptr {
		return errors.New("data must be a pointer.")
	}
	target := p.Elem()
	switch target.Kind() {
	case reflect.Struct:
		if err := ScanStruct(scanner, target); err != nil {
			return err
		}
	case reflect.Slice:
		if err := ScanSlice(scanner, target, p); err != nil {
			return err
		}
	default:
		if scanner.Next() {
			if err := scanner.Scan(p); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

func ScanStruct(scanner Scanner, target reflect.Value) error {
	columnNames, err := scanner.Columns()
	if err != nil {
		return err
	}
	fieldsAddrs, err := StructFieldsAddrs(target, Columns2Fields(columnNames))
	if err != nil {
		return err
	}
	if scanner.Next() {
		if err := scanner.Scan(fieldsAddrs...); err != nil {
			return err
		}
	}
	return nil
}

func ScanSlice(scanner Scanner, target, p reflect.Value) error {
	columnNames, err := scanner.Columns()
	if err != nil {
		return err
	}
	fieldNames := Columns2Fields(columnNames)
	elemType := target.Type().Elem()
	for scanner.Next() {
		target = reflect.Append(target, reflect.Zero(elemType))
		fieldsAddrs, err := StructFieldsAddrs(target.Index(target.Len()-1), fieldNames)
		if err != nil {
			return err
		}
		if err := scanner.Scan(fieldsAddrs...); err != nil {
			return err
		}
	}
	p.Elem().Set(target)
	return nil
}
