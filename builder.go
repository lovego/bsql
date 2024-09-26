package bsql

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lovego/date"
)

type Builder struct {
	manipulation     string
	table            string
	tableAlias       string
	join             []string
	groupBy          []string
	orderBy          []string
	having           string
	limit            int
	offset           int
	isForUpdate      bool
	fields           []string
	cols             []string
	returning        []string
	onConflict       string
	values           []interface{}
	updates          []string
	updateStruct     interface{}
	conditionBuilder ConditionBuilder
}

func NewBuilder() *Builder {
	return &Builder{conditionBuilder: ConditionBuilder{}}
}

func (b *Builder) Clone() *Builder {
	values := make([]interface{}, len(b.values))
	copy(values, b.values)
	return &Builder{
		manipulation: b.manipulation,
		table:        b.table,
		tableAlias:   b.tableAlias,
		join:         copyStringSlice(b.join),
		groupBy:      copyStringSlice(b.groupBy),
		orderBy:      copyStringSlice(b.orderBy),
		having:       b.having,
		limit:        b.limit,
		offset:       b.offset,
		isForUpdate:  b.isForUpdate,
		fields:       copyStringSlice(b.fields),
		cols:         copyStringSlice(b.cols),
		returning:    copyStringSlice(b.returning),
		onConflict:   b.onConflict,
		values:       values,
		updates:      copyStringSlice(b.updates),
		updateStruct: b.updateStruct,
		conditionBuilder: ConditionBuilder{
			wheres: copyStringSlice(b.conditionBuilder.wheres),
		},
	}
}

func (b *Builder) Clear() {
	b.manipulation = ""
	b.table = ""
	b.tableAlias = ""
	b.join = nil
	b.groupBy = nil
	b.orderBy = nil
	b.having = ""
	b.limit = 0
	b.offset = 0
	b.isForUpdate = false
	b.fields = nil
	b.cols = nil
	b.returning = nil
	b.onConflict = ""
	b.values = nil
	b.updates = nil
	b.updateStruct = nil
	b.conditionBuilder.Clear()
}

func (b *Builder) Alias(alias string) *Builder {
	b.tableAlias = alias
	return b
}

func (b *Builder) OrderBy(order string) *Builder {
	b.orderBy = append(b.orderBy, order)
	return b
}

func (b *Builder) Limit(limit int) *Builder {
	b.limit = limit
	return b
}

func (b *Builder) Offset(offset int) *Builder {
	b.offset = offset
	return b
}

// 添加自定义sql策略 Strategy接口形式
func (b *Builder) Strategies(strategies ...Strategy) *Builder {
	for _, strategy := range strategies {
		strategy.Execute(b)
	}
	return b
}

// 添加自定义sql策略 回调函数形式
func (b *Builder) StrategyFuncs(
	strategyFuncs ...func(b *Builder)) *Builder {
	for _, strategyFunc := range strategyFuncs {
		strategyFunc(b)
	}
	return b
}

func (b *Builder) Build() string {
	if b.table == "" {
		log.Panic("sql builder: table is required")
		return ""
	}
	switch b.manipulation {
	case manipulationSelect:
		return b.query()
	case manipulationInsert:
		return b.insert()
	case manipulationUpdate:
		return b.update()
	case manipulationDelete:
		return b.delete()
	default:
		log.Panic("sql builder: wrong manipulation")
		return ""
	}
}

func (b *Builder) buildWhere() string {
	condition := b.conditionBuilder.Build()
	if condition != "" {
		condition = "WHERE " + condition
	}
	return condition
}

func (b *Builder) tableName() string {
	table := b.table
	if b.tableAlias != "" {
		table += " AS " + b.tableAlias
	}
	return table
}

func (b *Builder) buildOrder() string {
	if len(b.orderBy) == 0 {
		return ""
	}
	return "ORDER BY " + strings.Join(b.orderBy, ",")
}

func (b *Builder) buildLimit() string {
	if b.limit <= 0 {
		return ""
	}
	sql := fmt.Sprintf("LIMIT %d", b.limit)
	if b.offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", b.offset)
	}
	return sql
}

func (b *Builder) Where(strs ...string) *Builder {
	b.conditionBuilder.Where(strs...)
	return b
}

func (b *Builder) Or(strs ...string) *Builder {
	b.conditionBuilder.Or(strs...)
	return b
}

func (b *Builder) Equal(dbField string, value interface{}) *Builder {
	b.conditionBuilder.Equal(dbField, value)
	return b
}

func (b *Builder) TryEqual(dbField string, value interface{}) *Builder {
	b.conditionBuilder.Equal(dbField, value)
	return b
}

func (b *Builder) Like(dbField, value string) *Builder {
	b.conditionBuilder.Like(dbField, value)
	return b
}

func (b *Builder) TryLike(dbField string, value string) *Builder {
	b.conditionBuilder.TryLike(dbField, value)
	return b
}

func (b *Builder) MultiLike(dbFields []string, value string) *Builder {
	b.conditionBuilder.MultiLike(dbFields, value)
	return b
}

func (b *Builder) TryMultiLike(dbFields []string, value string) *Builder {
	b.TryMultiLike(dbFields, value)
	return b
}

func (b *Builder) Between(
	dbField string, start, end interface{}) *Builder {
	b.conditionBuilder.Between(dbField, start, end)
	return b
}

func (b *Builder) In(dbField string, values interface{}) *Builder {
	b.conditionBuilder.In(dbField, values)
	return b
}

func (b *Builder) TryIn(dbField string, values interface{}) *Builder {
	b.conditionBuilder.TryIn(dbField, values)
	return b
}

func (b *Builder) NotIn(dbField string, values interface{}) *Builder {
	b.conditionBuilder.NotIn(dbField, values)
	return b
}

func (b *Builder) Any(dbField string, values interface{}) *Builder {
	b.conditionBuilder.Any(dbField, values)
	return b
}

func (b *Builder) TryAny(dbField string, values interface{}) *Builder {
	b.conditionBuilder.TryAny(dbField, values)
	return b
}

func (b *Builder) TryTimeRange(dbField string, startTime, endTime time.Time) *Builder {
	b.conditionBuilder.TryTimeRange(dbField, startTime, endTime)
	return b
}

func (b *Builder) TryDateRange(dbField string, startDate, endDate date.Date) *Builder {
	b.conditionBuilder.TryDateRange(dbField, startDate, endDate)
	return b
}

func copyStringSlice(src []string) []string {
	res := make([]string, len(src))
	copy(res, src)
	return res
}

const (
	manipulationInsert = "INSERT"
	manipulationDelete = "DELETE"
	manipulationUpdate = "UPDATE"
	manipulationSelect = "SELECT"

	TimeLayout = "2006-01-02 15:04:05"
)
