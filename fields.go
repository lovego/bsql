package bsql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lovego/struct_tag"
)

func Column2Field(column string) string {
	var parts []string
	for _, part := range strings.Split(column, "_") {
		parts = append(parts, strings.Title(part))
	}
	return strings.Join(parts, "")
}

/* 单词边界有两种
1. 非大写字符，且下一个是大写字符
2. 大写字符，且下一个是大写字符，且下下一个是非大写字符
*/
func Field2Column(str string) string {
	var slice []string
	start := 0
	for end, char := range str {
		if end+1 < len(str) {
			next := str[end+1]
			if char < 'A' || char > 'Z' {
				if next >= 'A' && next <= 'Z' { // 非大写下一个是大写
					slice = append(slice, str[start:end+1])
					start, end = end+1, end+1
				}
			} else if end+2 < len(str) && (next >= 'A' && next <= 'Z') {
				if next2 := str[end+2]; next2 < 'A' || next2 > 'Z' {
					slice = append(slice, str[start:end+1])
					start, end = end+1, end+1
				}
			}
		} else {
			slice = append(slice, str[start:end+1])
		}
	}
	return strings.ToLower(strings.Join(slice, "_"))
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
	traverseStructFields(reflect.TypeOf(strct), func(field reflect.StructField) {
		if notIn(field.Name, exclude) {
			result = append(result, field.Name)
		}
	})
	return
}

func ColumnsComments(table string, strct interface{}) (result string) {
	traverseStructFields(reflect.TypeOf(strct), func(field reflect.StructField) {
		if comment, _ := struct_tag.Lookup(string(field.Tag), "comment"); comment != "" {
			result += fmt.Sprintf(
				"comment on column %s.%s is %s;\n", table, Field2Column(field.Name), Q(comment),
			)
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
			if value, ok := struct_tag.Lookup(string(field.Tag), `sql`); !ok || value != "-" {
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
