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
	Increment   int
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
	var incrementByCol string
	for i := range columns {
		if strings.Contains(columns[i], "PRIMARY KEY") {
			incrementByCol = strings.Split(columns[i], " ")[0]
		}
		columns[i] = `  ` + strings.TrimRight(strings.TrimSpace(columns[i]), `,`)
	}
	columnsStr := strings.Join(columns, ",\n")
	var incrementBy string
	if incrementByCol != "" && t.Increment > 0 {
		incrementBy = fmt.Sprintf(
			"ALTER sequence %s_%s_seq INCREMENT %d;\n", t.Name, incrementByCol, t.Increment,
		)
	}

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
%s%sCOMMENT ON TABLE %s is %s;
%s`, t.Name, columnsStr, optionsSql,
		incrementBy, extSqlsStr, t.Name, Q(t.Desc),
		ColumnsComments(t.Name, t.Struct),
	)
	return sql
}

func ensureTail(str string, tail byte) string {
	if str[len(str)-1] == tail {
		return str
	}
	return str + string(tail)
}
