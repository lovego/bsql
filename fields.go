package bsql

import (
	"reflect"
	"strings"

	"github.com/lovego/strs"
	"github.com/lovego/struct_tag"
	"github.com/lovego/structs"
)

func Field2Column(field string) string {
	var path []string
	for _, name := range strings.Split(field, ".") {
		path = append(path, strs.CamelToSnake(name))
	}
	return strings.Join(path, ".")
}

func Fields2Columns(fields []string) (result []string) {
	for _, field := range fields {
		result = append(result, Field2Column(field))
	}
	return
}

func Fields2ColumnsStr(fields []string) string {
	return strings.Join(Fields2Columns(fields), ",")
}

func FieldsToColumns(fields []string, prefix string, exclude []string) (result []string) {
	for _, field := range fields {
		if notIn(field, exclude) {
			result = append(result, prefix+Field2Column(field))
		}
	}
	return
}

func FieldsToColumnsStr(fields []string, prefix string, exclude []string) string {
	return strings.Join(FieldsToColumns(fields, prefix, exclude), ",")
}

func FieldsFromStruct(strct interface{}, exclude []string) (result []string) {
	traverseStructFields(reflect.TypeOf(strct), func(field reflect.StructField) {
		if notIn(field.Name, exclude) {
			result = append(result, field.Name)
		}
	})
	return
}

func traverseStructFields(typ reflect.Type, fn func(field reflect.StructField)) {
	structs.TraverseType(typ, func(field reflect.StructField) bool {
		return struct_tag.Get(string(field.Tag), `sql`) == "-"
	}, fn)
}

func notIn(target string, slice []string) bool {
	for _, elem := range slice {
		if elem == target {
			return false
		}
	}
	return true
}
