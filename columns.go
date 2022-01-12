package bsql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lovego/bsql/scan"
	"github.com/lovego/struct_tag"
)

func Columns2Fields(columns []string) (result []string) {
	for _, column := range columns {
		result = append(result, strings.Join(scan.Column2FieldPath(column), "."))
	}
	return result
}

func ColumnsFromStruct(strct interface{}, exclude []string) (result []string) {
	return Fields2Columns(FieldsFromStruct(strct, exclude))
}

func ColumnsStrFromStruct(strct interface{}, exclude []string) string {
	return Fields2ColumnsStr(FieldsFromStruct(strct, exclude))
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
