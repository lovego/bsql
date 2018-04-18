package bsql

import (
	"database/sql"
)

type columnType struct {
	FieldName string
	*sql.ColumnType
}

func getColumns(rows rowsType) ([]columnType, error) {
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	if len(columnTypes) == 0 { // for unit test
		return getColumnsFromNames(rows)
	}
	var columns []columnType
	for _, colType := range columnTypes {
		columns = append(columns, columnType{
			FieldName: Column2Field(colType.Name()), ColumnType: colType,
		})
	}
	return columns, nil
}

func getColumnsFromNames(rows rowsType) ([]columnType, error) {
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var columns []columnType
	for _, columnName := range columnNames {
		columns = append(columns, columnType{FieldName: Column2Field(columnName)})
	}
	return columns, nil
}
