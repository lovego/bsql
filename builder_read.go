package bsql

import (
	"fmt"
	"strings"
)

func (b *Builder) Select(table string) *Builder {
	b.manipulation = manipulationSelect
	b.table = table
	return b
}

func (b *Builder) SelectSubQuery(subQueryBuilder *Builder) *Builder {
	b.manipulation = manipulationSelect
	sub := subQueryBuilder.Build()
	b.table = "(" + sub + ")"
	return b
}

func (b *Builder) query() string {
	return strings.Join([]string{
		b.selectFields(),
		"FROM",
		b.tableName(),
		b.buildJoin(),
		b.buildWhere(),
		b.buildGroup(),
		b.buildOrder(),
		b.buildHaving(),
		b.buildLimit(),
		b.buildForUpdate(),
	}, " ")
}

func (b *Builder) Join(joinType, table, as, on string) *Builder {
	b.join = append(b.join,
		fmt.Sprintf("%s JOIN %s AS %s ON %s", joinType, table, as, on))
	return b
}

func (b *Builder) LeftJoin(table, as, on string) *Builder {
	return b.Join("LEFT", table, as, on)
}

func (b *Builder) RightJoin(table, as, on string) *Builder {
	return b.Join("RIGHT", table, as, on)
}

func (b *Builder) InnerJoin(table, as, on string) *Builder {
	return b.Join("INNER", table, as, on)
}

func (b *Builder) GroupBy(group string) *Builder {
	b.groupBy = append(b.groupBy, group)
	return b
}

func (b *Builder) Having(having string) *Builder {
	b.having = having
	return b
}

func (b *Builder) Fields(fields ...string) *Builder {
	b.fields = append(b.fields, fields...)
	return b
}

func (b *Builder) ForUpdate() *Builder {
	b.isForUpdate = true
	return b
}
func (b *Builder) buildForUpdate() string {
	if b.isForUpdate {
		return "FOR UPDATE"
	}
	return ""
}

func (b *Builder) buildJoin() string {
	if len(b.join) == 0 {
		return ""
	}
	return strings.Join(b.join, " ")
}

func (b *Builder) buildGroup() string {
	if len(b.groupBy) == 0 {
		return ""
	}
	return "GROUP BY " + strings.Join(b.groupBy, ",")
}
func (b *Builder) buildHaving() string {
	if b.having == "" {
		return ""
	}
	return "HAVING " + b.having
}

func (b *Builder) selectFields() string {
	fields := "*"
	if len(b.fields) > 0 {
		fields = strings.Join(b.fields, ",")
	}
	return fmt.Sprintf("%s %s", b.manipulation, fields)
}
