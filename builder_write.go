package bsql

import (
	"fmt"
	"log"
	"strings"
)

func (b *Builder) Insert(table string) *Builder {
	b.manipulation = manipulationInsert
	b.table = table
	return b
}

func (b *Builder) Update(table string) *Builder {
	b.manipulation = manipulationUpdate
	b.table = table
	return b
}
func (b *Builder) Delete(table string) *Builder {
	b.manipulation = manipulationDelete
	b.table = table
	return b
}

func (b *Builder) Cols(cols ...string) *Builder {
	b.cols = append(b.cols, cols...)
	return b
}

func (b *Builder) Set(data ...string) *Builder {
	b.updates = append(b.updates, data...)
	return b
}

func (b *Builder) SetMap(data map[string]interface{}) *Builder {
	for k, v := range data {
		b.updates = append(b.updates,
			fmt.Sprintf("%s = %s", k, bsql.V(v)))
	}
	return b
}

func (b *Builder) SetStruct(data interface{}) *Builder {
	b.updateStruct = data
	return b
}

func (b *Builder) Values(values []interface{}) *Builder {
	b.values = values
	return b
}

func (b *Builder) insert() string {
	if len(b.values) == 0 {
		log.Panic("sql builder: inserting values are required")
		return ""
	}
	cols := b.insertCols()
	if len(cols) == 0 {
		log.Panic("sql builder: inserting fields are required")
		return ""
	}
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES %s %s",
		b.tableName(),
		strings.Join(bsql.Fields2Columns(cols), ","),
		bsql.StructValues(b.values, cols),
		b.buildReturning(),
	)
}

func (b *Builder) update() string {
	return strings.Join([]string{
		b.manipulation,
		b.tableName(),
		"SET",
		b.buildUpdates(),
		b.buildWhere(),
		b.buildOrder(),
		b.buildLimit(),
		b.buildReturning(),
	}, " ")
}

func (b *Builder) delete() string {
	where := b.buildWhere()
	if where == "" {
		log.Panic("sql builder: deleting condition are required")
		return ""
	}
	return strings.Join([]string{
		b.manipulation,
		"FROM",
		b.tableName(),
		where,
		b.buildOrder(),
		b.buildLimit(),
		b.buildReturning(),
	}, " ")
}

func (b *Builder) buildUpdates() string {
	if b.updateStruct != nil {
		cols := b.updateCols()
		return fmt.Sprintf("(%s) = %s",
			strings.Join(bsql.Fields2Columns(cols), ","),
			bsql.StructValues([]interface{}{b.updateStruct}, cols))
	}
	if len(b.updates) == 0 {
		log.Panic("sql builder: updating values are required")
		return ""
	}
	return strings.Join(b.updates, ",")
}

func (b *Builder) buildReturning() string {
	if len(b.returning) == 0 {
		return ""
	}
	return strings.Join(b.returning, ",")
}

func (b *Builder) insertCols() []string {
	cols := b.cols
	if len(cols) == 0 {
		strct := b.values[0]
		cols = bsql.FieldsFromStruct(strct,
			[]string{"Id", "UpdatedBy", "UpdatedAt"})
	}
	return cols
}

func (b *Builder) updateCols() []string {
	cols := b.cols
	if len(cols) == 0 {
		cols = bsql.FieldsFromStruct(b.updateStruct,
			[]string{"CreatedBy", "CreatedAt"})
	}
	return cols
}

func (b *Builder) OnConflict(fields string, do string) *Builder {
	if fields == "" {
		b.onConflict = "ON CONFLICT DO " + do
	} else {
		b.onConflict = fmt.Sprintf("ON CONFLICT (%s) DO %s", fields, do)
	}
	return b
}

func (b *Builder) OnConflictDoNothing() *Builder {
	return b.OnConflict("", "NOTHING")
}

func (b *Builder) Returning(fields ...string) *Builder {
	b.returning = append(b.returning, fields...)
	return b
}
