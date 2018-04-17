package psql

import (
	"errors"
	"strings"
)

func Q(q string) string {
	return "'" + strings.Replace(q, "'", "''", -1) + "'"
}

func Columns2Fields(columns []string) (result []string) {
	for _, column := range columns {
		result = append(Column2Field(column))
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

func StructFieldsAddrs(structValue reflect.Value, fieldNames []string) (result []interface{}, err error) {
	for _, fieldName := range fieldNames {
		if field := structValue.FieldByName(fieldName); field.IsValid() {
			result = append(result, field.Addr().Interface())
		} else {
			return nil, errors.New("no field: '" + field + "' in struct.")
		}
	}
	return
}
