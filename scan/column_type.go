package scan

import (
	"database/sql"
	"strings"

	"github.com/lovego/strs"
)

type ColumnType struct {
	FieldPath []string
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
		var path = Column2FieldPath(colType.Name())
		columns = append(columns, ColumnType{
			FieldPath:  path,
			FieldName:  strings.Join(path, "."),
			ColumnType: colType,
		})
	}
	return columns, nil
}

func Column2FieldPath(column string) (path []string) {
	for _, name := range strings.Split(column, ".") {
		path = append(path, strs.SnakeToCamel(name))
	}
	return
}
