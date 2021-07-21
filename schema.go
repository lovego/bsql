package bsql

import (
	"fmt"
	"log"
	"strings"
)

type Table struct {
	Name        string
	Desc        string
	Struct      interface{}
	Constraints []string
	Options     []string
	ExtraSqls   []string
}

// Sql add create table sql
// Using Name as table name, Desc as table description, Struct struct as table columns and comments.
func (t Table) Sql() string {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		log.Panic("table name required")
	}
	t.Desc = strings.TrimSpace(t.Desc)
	if t.Desc == "" {
		log.Panic("table desc required")
	}
	if t.Struct == nil {
		log.Panic("table struct required")
	}
	columns := columnsFromStruct(t.Struct)
	columns = append(columns, t.Constraints...)
	for i := range columns {
		columns[i] = `  ` + strings.TrimRight(strings.TrimSpace(columns[i]), `,`)
	}
	columnsStr := strings.Join(columns, ",\n")

	for i := range t.Options {
		t.Options[i] = "\n" + strings.TrimSpace(t.Options[i])
	}
	optionsSql := strings.Join(t.Options, "")

	for i := range t.ExtraSqls {
		t.ExtraSqls[i] = ensureTail(strings.TrimSpace(t.ExtraSqls[i]), ';') + "\n"
	}
	extSqlsStr := strings.Join(t.ExtraSqls, "")

	sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
%s
)%s;
%sCOMMENT ON TABLE %s is %s;
%s`, t.Name, columnsStr, optionsSql, extSqlsStr, t.Name, Q(t.Desc), ColumnsComments(t.Name, t.Struct))
	return sql
}

func ensureTail(str string, tail byte) string {
	if str[len(str)-1] == tail {
		return str
	}
	return str + string(tail)
}
