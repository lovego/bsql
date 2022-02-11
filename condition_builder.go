package bsql

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/lovego/date"
)

type ConditionBuilder struct {
	wheres []string
}

// 生成最终的sql
func (b *ConditionBuilder) Build() string {
	return strings.TrimSpace(strings.Join(b.wheres, " AND "))
}

// 清空
func (b *ConditionBuilder) Clear() {
	b.wheres = nil
}

// 添加多个查询AND条件
func (b *ConditionBuilder) Where(strs ...string) *ConditionBuilder {
	for _, str := range strs {
		if str != "" {
			// tip: 括号包裹条件，防止条件之间相互影响优先级
			b.wheres = append(b.wheres, "("+str+")")
		}
	}
	return b
}

//添加多个OR条件
func (b *ConditionBuilder) Or(strs ...string) *ConditionBuilder {
	var cons []string
	for _, str := range strs {
		if str != "" {
			cons = append(cons, "("+str+")")
		}
	}
	if len(cons) > 0 {
		b.Where(strings.Join(cons, " OR "))
	}
	return b
}

// 添加相等条件
func (b *ConditionBuilder) Equal(dbField string, value interface{}) *ConditionBuilder {
	return b.Where(fmt.Sprintf("%s = %s", dbField, bsql.V(value)))
}

// 添加相等条件，value为零值时跳过
func (b *ConditionBuilder) TryEqual(dbField string, value interface{}) *ConditionBuilder {
	if value == nil ||
		value == "" ||
		value == 0 ||
		value == false {
		return b
	}
	v := reflect.ValueOf(value)
	kind := v.Kind()
	if kind == reflect.Array || kind == reflect.Slice {
		if v.Len() == 0 {
			return b
		}
	}
	return b.Equal(dbField, value)
}

// 添加LIKE条件，左右模糊匹配，
// 如果需要单边模糊匹配，请使用Where
func (b *ConditionBuilder) Like(dbField, value string) *ConditionBuilder {
	return b.Where(fmt.Sprintf("%s LIKE %s", dbField, bsql.Q("%"+value+"%")))
}

// 添加LIKE条件，左右模糊匹配，value为零值时跳过
func (b *ConditionBuilder) TryLike(dbField string, value string) *ConditionBuilder {
	if value := strings.TrimSpace(value); value != "" {
		return b.Like(dbField, value)
	}
	return b
}

// 添加多个LIKE条件
func (b *ConditionBuilder) MultiLike(dbFields []string, value string) *ConditionBuilder {
	v := bsql.Q("%" + value + "%")
	var cons []string
	for _, field := range dbFields {
		cons = append(cons, fmt.Sprintf("%s LIKE %s", field, v))
	}
	return b.Or(cons...)
}

// 添加多个LIKE条件，value为零值时跳过
func (b *ConditionBuilder) TryMultiLike(dbFields []string, value string) *ConditionBuilder {
	if v := strings.TrimSpace(value); v != "" {
		return b.MultiLike(dbFields, v)
	}
	return b
}

// 添加BETWEEN条件
func (b *ConditionBuilder) Between(
	dbField string, start, end interface{}) *ConditionBuilder {
	return b.Where(fmt.Sprintf("%s BETWEEN %s AND %s",
		dbField, bsql.V(start), bsql.V(end)))
}

// 添加IN条件
func (b *ConditionBuilder) In(dbField string, values interface{}) *ConditionBuilder {
	if condition := buildInCondition(dbField, values); condition != "" {
		return b.Where(condition)
	}
	return b.Where("1=0")
}

// 添加IN条件，value为零值时跳过
func (b *ConditionBuilder) TryIn(dbField string, values interface{}) *ConditionBuilder {
	if condition := buildInCondition(dbField, values); condition != "" {
		return b.Where(condition)
	}
	return b
}

// 添加NOT IN条件
func (b *ConditionBuilder) NotIn(dbField string, values interface{}) *ConditionBuilder {
	if condition := buildNotInCondition(dbField, values); condition != "" {
		return b.Where(condition)
	}
	return b
}

// 添加Any条件
// values 可传类型：
// 		string: 子查询sql
// 		array/slice: 结果集，效果同In
func (b *ConditionBuilder) Any(dbField string, values interface{}) *ConditionBuilder {
	if condition := buildAnyCondition(dbField, values); condition != "" {
		return b.Where(condition)
	}
	return b.Where("1=0")
}

// 添加IN条件，value为零值时跳过
func (b *ConditionBuilder) TryAny(dbField string, values interface{}) *ConditionBuilder {
	if condition := buildAnyCondition(dbField, values); condition != "" {
		return b.Where(condition)
	}
	return b
}

// 添加时间范围条件，value为零值时跳过
func (b *ConditionBuilder) TryTimeRange(
	dbField string, startTime, endTime time.Time) *ConditionBuilder {
	if !startTime.IsZero() && !endTime.IsZero() {
		return b.Between(dbField, startTime, endTime)
	}
	if !startTime.IsZero() {
		return b.Where(fmt.Sprintf("%s >= %s", dbField, bsql.V(startTime)))
	}
	if !endTime.IsZero() {
		return b.Where(fmt.Sprintf("%s <= %s", dbField, bsql.V(endTime)))
	}
	return b
}

// 添加日期范围条件，value为零值时跳过
func (b *ConditionBuilder) TryDateRange(
	dbField string, startDate, endDate date.Date) *ConditionBuilder {
	startTime := startDate.Time
	endTime := endDate.Time
	if !startDate.IsZero() {
		startTime, _ = time.Parse(TimeLayout, startDate.String()+" 00:00:00")
	}
	if !endDate.IsZero() {
		endTime, _ = time.Parse(TimeLayout, endDate.String()+" 23:59:59")
	}
	return b.TryTimeRange(dbField, startTime, endTime)
}

func buildInCondition(field string, values interface{}) string {
	if v := sliceValue(values); v != "" {
		return fmt.Sprintf("%s IN (%s)", field, v)
	}
	return ""
}
func buildNotInCondition(field string, values interface{}) string {
	if v := sliceValue(values); v != "" {
		return fmt.Sprintf("%s NOT IN (%s)", field, v)
	}
	return ""
}

func buildAnyCondition(field string, values interface{}) string {
	switch values.(type) {
	case string:
		if values == "" {
			return ""
		}
		return fmt.Sprintf("%s = ANY(%s)", field, values)
	default:
		if v := sliceValue(values); v != "" {
			return fmt.Sprintf("%s = ANY(ARRAY[%s])", field, v)
		}
		return ""
	}
}

func sliceValue(values interface{}) string {
	if values == nil {
		return ""
	}
	v := reflect.ValueOf(values)
	kind := v.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		return ""
	}
	vLen := v.Len()
	if vLen == 0 {
		return ""
	}
	var s []string
	for i := 0; i < vLen; i++ {
		s = append(s, bsql.V(v.Index(i).Interface()))
	}
	return strings.Join(s, ",")
}
