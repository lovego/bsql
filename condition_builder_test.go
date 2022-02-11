package bsql

import (
	"fmt"
	"testing"
)

func TestConditionBuilder(t *testing.T) {
	ExampleCondition()
	ExampleCondition2()
}

func ExampleCondition() {
	builder := ConditionBuilder{}

	// Where
	fmt.Println("Where:")
	builder.Where("a = 2", "b = 3")
	fmt.Println(builder.Build()) // print: (a = 2) AND (b = 3)

	// Clear
	fmt.Println("Clear:")
	builder.Clear()
	fmt.Println(builder.Build()) //  print ""

	// Or
	fmt.Println("Where:")
	builder.Or("a = 2", "b = 3")
	fmt.Println(builder.Build()) // print: ((a = 2) OR (b = 3))
	builder.Clear()

	// Equal
	fmt.Println("Equal:")
	builder.Equal("a", 3).Equal("b", "bbb")
	fmt.Println(builder.Build()) //  print: (a = 3) AND (b = 'bbb')
	builder.Clear()

	// In
	fmt.Println("In:")
	builder.In("a", []int{1, 2, 3}).In("b", []string{"b", "bb", "bbb"})
	fmt.Println(builder.Build()) //  print: (a IN (1,2,3)) AND (b IN ('b','bb','bbb'))
	builder.Clear()

	// Like
	fmt.Println("Like:")
	builder.Like("a", "AA")
	fmt.Println(builder.Build()) //  print: (a LIKE '%AA%')
	builder.Clear()

	// MultiLike
	fmt.Println("MultiLike:")
	builder.MultiLike([]string{"a", "b"}, "AA")
	fmt.Println(builder.Build()) //  print: ((a LIKE '%AA%') OR (b LIKE '%AA%'))
	builder.Clear()
}

func ExampleCondition2() {
	builder := ConditionBuilder{}

	// Between
	fmt.Println("Between:")
	builder.Between("a", 1, 100)
	fmt.Println(builder.Build()) //  print: (a BETWEEN 1 AND 100)
	builder.Clear()

	// Any
	fmt.Println("Any:")
	builder.Any("a", []int{1, 2, 3}).
		Any("b", "select id from account_sets.account_sets where id < 10")
	fmt.Println(builder.Build())
	/*  print:
	 	(a = ANY(ARRAY[1,2,3])) AND
		(b = ANY(select id from account_sets.account_sets where id < 10))
	*/
	builder.Clear()
}
