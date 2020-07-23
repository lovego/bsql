package bsql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lovego/bsql/scan"
	"github.com/lovego/struct_tag"
	"github.com/lovego/structs"
)

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
		result = append(result, scan.Column2Field(column))
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

func ColumnsComments(table string, strct interface{}) (result string) {
	traverseStructFields(reflect.TypeOf(strct), func(field reflect.StructField) {
		comment, _ := struct_tag.Lookup(string(field.Tag), "c")
		if comment == "" {
			comment, _ = struct_tag.Lookup(string(field.Tag), "comment")
		}
		if comment != "" {
			result += fmt.Sprintf(
				"COMMENT ON COLUMN %s.%s IS %s;\n", table, Field2Column(field.Name), Q(comment),
			)
		}
	})
	return
}

func traverseStructFields(typ reflect.Type, fn func(field reflect.StructField)) {
	structs.TraverseType(typ, func(field reflect.StructField) {
		if value, ok := struct_tag.Lookup(string(field.Tag), `sql`); !ok || value != "-" {
			fn(field)
		}
	})
}

func notIn(target string, slice []string) bool {
	for _, elem := range slice {
		if elem == target {
			return false
		}
	}
	return true
}
