package bsql

import (
	"reflect"
	"regexp"
	"strings"
)

func Column2Field(column string) string {
	var parts []string
	for _, part := range strings.Split(column, "_") {
		parts = append(parts, strings.Title(part))
	}
	return strings.Join(parts, "")
}

var camel = regexp.MustCompile("(^[^A-Z0-9]*|[A-Z0-9]*)([A-Z0-9][^A-Z]+|$)")

func Field2Column(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "_"))
}

func Columns2Fields(columns []string) (result []string) {
	for _, column := range columns {
		result = append(result, Column2Field(column))
	}
	return result
}

func Fields2Columns(fields []string) (result []string) {
	for _, field := range fields {
		result = append(result, Field2Column(field))
	}
	return
}

func Fields2ColumnsStr(fields []string) string {
	var result []string
	for _, field := range fields {
		result = append(result, Field2Column(field))
	}
	return strings.Join(result, ",")
}

func FieldsFromStruct(strct interface{}, exclude []string) (result []string) {
	traverseStructFields(reflect.ValueOf(strct).Type(), func(field reflect.StructField) {
		if notIn(field.Name, exclude) {
			result = append(result, field.Name)
		}
	})
	return
}

func traverseStructFields(typ reflect.Type, fn func(field reflect.StructField)) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// exported field has an empty PkgPath
		if (!field.Anonymous || !traverseStructFields(field.Type, fn)) && field.PkgPath == "" {
			if value, ok := field.Tag.Lookup(`sql`); !ok || value != "-" {
				fn(field)
			}
		}
	}
	return true
}

func notIn(target string, slice []string) bool {
	for _, elem := range slice {
		if elem == target {
			return false
		}
	}
	return true
}
