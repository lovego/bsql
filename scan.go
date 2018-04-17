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
	case reflect.Slice:
		if err := Scan2Slice(scanner, target, p); err != nil {
			return err
		}
	case reflect.Struct:
		if err := Scan2Struct(scanner, target); err != nil {
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

func Scan2Slice(scanner Scanner, target, p reflect.Value) error {
	elemType := target.Type().Elem()
	if elemType.Kind() == reflect.Struct {
		return Scan2StructSlice(scanner, target, p)
	}
	for scanner.Next() {
		target = reflect.Append(target, reflect.Zero(elemType))
		if err := scanner.Scan(target.Index(target.Len() - 1).Addr().Interface()); err != nil {
			return err
		}
	}
	p.Elem().Set(target)
	return nil
}
func Scan2StructSlice(scanner Scanner, target, p reflect.Value) error {
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

func Scan2Struct(scanner Scanner, target reflect.Value) error {
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

func StructFieldsAddrs(structValue reflect.Value, fieldNames []string) ([]interface{}, error) {
	var result []interface{}
	for _, fieldName := range fieldNames {
		if field := structValue.FieldByName(fieldName); field.IsValid() {
			result = append(result, field.Addr().Interface())
		} else {
			return nil, errors.New("no field: '" + fieldName + "' in struct.")
		}
	}
	return result, nil
}
