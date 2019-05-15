package bsql

import (
	"fmt"
	"strings"
)

// TableSql add create table sql
// Using name as table name, model struct as table columns and comments.
func TableSql(name string, model interface{}, constraints, extSqls []string) string {
	lines := columnsFromStruct(model)
	lines = append(lines, constraints...)

	for i := range lines {
		lines[i] = `  ` + strings.TrimRight(strings.TrimSpace(lines[i]), `,`)
	}
	linesStr := strings.Join(lines, ",\n")

	for i := range extSqls {
		extSqls[i] = ensureTail(strings.TrimSpace(extSqls[i]), ';')
	}
	extSql := strings.Join(extSqls, "\n")

	comments := ColumnsComments(name, model)

	sql := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
%s
);
%s
%s`, name, linesStr, extSql, comments)
	return sql
}

func ensureTail(str string, tail byte) string {
	if str[len(str)-1] == tail {
		return str
	}
	return str + string(tail)
}
