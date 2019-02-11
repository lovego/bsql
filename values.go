package bsql

import (
	"log"
	"reflect"
	"strings"
)

// Values return the contents following the sql keyword "values"
func Values(data interface{}) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		switch value.Type().Elem().Kind() {
		case reflect.Slice, reflect.Array:
			var slice []string
			for i := 0; i < value.Len(); i++ {
				slice = append(slice, "("+SliceContents(value.Index(i))+")")
			}
			return strings.Join(slice, ",")
		default:
			return "(" + SliceContents(value) + ")"
		}
	default:
		return "(" + V(value.Interface()) + ")"
	}
}

func SliceContents(value reflect.Value) string {
	var slice []string
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, V(value.Index(i).Interface()))
	}
	return strings.Join(slice, ",")
}

func StructValues(data interface{}, fields []string) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		var slice []string
		for i := 0; i < value.Len(); i++ {
			slice = append(slice, StructValuesIn(value.Index(i), fields))
		}
		return strings.Join(slice, ",")
	default:
		return StructValuesIn(value, fields)
	}
}

func StructValuesIn(value reflect.Value, fields []string) string {
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		log.Panic("bsql: data must be struct or struct slice.")
	}
	var slice []string
	for _, fieldName := range fields {
		field := structField(value, fieldName)
		if !field.IsValid() {
			log.Panic("bsql: no field '" + fieldName + "' in struct")
		}
		slice = append(slice, V(field.Interface()))
	}
	return "(" + strings.Join(slice, ",") + ")"
}

func structField(strct reflect.Value, fieldName string) reflect.Value {
	if strings.IndexByte(fieldName, '.') <= 0 {
		return strct.FieldByName(fieldName)
	}
	for _, name := range strings.Split(fieldName, ".") {
		strct = strct.FieldByName(name)
		if !strct.IsValid() {
			return strct
		}
	}
	return strct
}
