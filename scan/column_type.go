package scan

import (
	"database/sql"
	"strings"
)

type ColumnType struct {
	FieldName string
	*sql.ColumnType
}

func ColumnTypes(rows *sql.Rows) ([]ColumnType, error) {
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	var columns []ColumnType
	for _, colType := range columnTypes {
		columns = append(columns, ColumnType{
			FieldName: Column2Field(colType.Name()), ColumnType: colType,
		})
	}
	return columns, nil
}

func Column2Field(column string) string {
	var parts []string
	for _, part := range strings.Split(column, "_") {
		parts = append(parts, strings.Title(part))
	}
	return strings.Join(parts, "")
}
