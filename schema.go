package bsql

import (
	"fmt"
	"strings"
)

// TableSql add create table sql
// Using name as table name, model struct as table columns and comments.
func TableSql(name string, model interface{}, constraints, extSqls []string) string {
	columns := columnsFromStruct(model)
	columns = append(columns, constraints...)
	for i := range columns {
		columns[i] = `  ` + strings.TrimRight(strings.TrimSpace(columns[i]), `,`)
	}
	columnsStr := strings.Join(columns, ",\n")

	for i := range extSqls {
		extSqls[i] = ensureTail(strings.TrimSpace(extSqls[i]), ';') + "\n"
	}
	extSqlsStr := strings.Join(extSqls, "")

	sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
%s
);
%s%s`, name, columnsStr, extSqlsStr, ColumnsComments(name, model))
	return sql
}

func ensureTail(str string, tail byte) string {
	if str[len(str)-1] == tail {
		return str
	}
	return str + string(tail)
}
