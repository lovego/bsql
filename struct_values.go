package bsql

import (
	"log"
	"reflect"
	"strings"
)

func StructValues(data interface{}, fields []string) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		var slice []string
		for i := 0; i < value.Len(); i++ {
			slice = append(slice, "("+StructFieldsReflect(value.Index(i), fields)+")")
		}
		return strings.Join(slice, ",")
	default:
		return "(" + StructFieldsReflect(value, fields) + ")"
	}
}

func StructFields(value interface{}, fields []string) string {
	return StructFieldsReflect(reflect.ValueOf(value), fields)
}

func StructFieldsReflect(value reflect.Value, fields []string) string {
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
	return strings.Join(slice, ",")
}

func StructValuesWithType(data interface{}, fields []string) string {
	value := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		var slice []string
		for i := 0; i < value.Len(); i++ {
			slice = append(slice, StructFieldsWithType(value.Index(i), typ.Elem(), fields))
		}
		return strings.Join(slice, ",")
	default:
		return StructFieldsWithType(value, typ, fields)
	}
}

func StructFieldsWithType(value reflect.Value, typ reflect.Type, fields []string) string {
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		value = value.Elem()
	}
	if k := typ.Kind(); k == reflect.Ptr || k == reflect.Interface {
		typ = typ.Elem()
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
		var fieldType string
		if fieldName != `Id` {
			structField, ok := typ.FieldByName(fieldName)
			if !ok {
				log.Panic("bsql: no field '" + fieldName + "' in struct")
			}
			fieldType = "::" + getColumnType(structField)
		}
		slice = append(slice, V(field.Interface())+fieldType)
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
