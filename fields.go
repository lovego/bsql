package bsql

import (
	"reflect"
	"strings"
)

func Columns2Fields(columns []string) (result []string) {
	for _, column := range columns {
		result = append(result, Column2Field(column))
	}
	return result
}

func Column2Field(column string) string {
	var parts []string
	for _, part := range strings.Split(column, "_") {
		parts = append(parts, strings.Title(part))
	}
	return strings.Join(parts, "")
}

func Field2Column(field string) string {
	return ""
}

func FieldsFromStruct(v interface{}, exclude []string) (result []string) {
	LoopStructFields(reflect.ValueOf(v).Type(), func(name string) {
		if NotIn(name, exclude) {
			result = append(result, name)
		}
	})
	return
}

func LoopStructFields(typ reflect.Type, fn func(name string)) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			LoopStructFields(field.Type, fn)
		} else {
			fn(field.Name)
		}
	}
}

func NotIn(target string, slice []string) bool {
	for _, elem := range slice {
		if elem == target {
			return false
		}
	}
	return true
}
