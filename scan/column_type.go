package scan

import (
	"database/sql"
	"strings"
)

type columnType struct {
	FieldName string
	*sql.ColumnType
}

func getColumns(rows *sql.Rows) ([]columnType, error) {
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	var columns []columnType
	for _, colType := range columnTypes {
		columns = append(columns, columnType{
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
