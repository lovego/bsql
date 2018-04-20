package bsql

import (
	"reflect"
	"strings"
)

func CreatingTableColumns(strct interface{}) string {
	var columns []string
	traverseStructFields(reflect.ValueOf(strct).Type(), func(field reflect.StructField) {
		columns = append(columns, Field2Column(field.Name)+" "+getColumnDefinition(field))
	})
	return strings.Join(columns, ",\n")
}

func getColumnDefinition(field reflect.StructField) string {
	if tag, ok := field.Tag.Lookup(`sql`); ok {
		tag = strings.TrimSpace(tag)
		if tag != "" && tag != "-" {
			return tag
		}
	}
	return ""
}

func getColumnType(field reflect.Value) string {
	typ := field.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	}
	return ""

}
