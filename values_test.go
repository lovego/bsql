package bsql

import (
	"fmt"
	"reflect"
)

func ExampleValues() {
	fmt.Println(Values(3))
	fmt.Println(Values("a'bc"))
	fmt.Println(Values([]int{1, 2, 3}))
	fmt.Println(Values([]string{"a", "b", "c"}))
	fmt.Println(Values([][]interface{}{
		{1, "a", true}, {2, "b", true}, {3, "c", false},
	}))
	// Output:
	// (3)
	// ('a''bc')
	// (1,2,3)
	// ('a','b','c')
	// (1,'a',true),(2,'b',true),(3,'c',false)
}

func ExampleValues_map() {
	m := map[string]bool{"('1','2')": true, "('2','1')": true}
	result := Values(map[string]interface{}{"1": nil, "2": nil})
	fmt.Println(m[result])

	m = map[string]bool{"(1,2)": true, "(2,1)": true}
	result = Values(map[int]interface{}{1: nil, 2: nil})
	fmt.Println(m[result])

	m = map[string]bool{"(1,2),(3,4)": true, "(3,4),(1,2)": true}
	result = Values(map[[2]int]interface{}{[2]int{1, 2}: nil, [2]int{3, 4}: nil})
	fmt.Println(m[result])

	// Output:
	// true
	// true
	// true
}

func ExampleSingleColumnValues() {
	fmt.Println(SingleColumnValues(3))
	fmt.Println(SingleColumnValues("a'bc"))
	fmt.Println(SingleColumnValues([]int{1, 2, 3}))
	fmt.Println(SingleColumnValues([]string{"a", "b", "c"}))
	// Output:
	// (3)
	// ('a''bc')
	// (1),(2),(3)
	// ('a'),('b'),('c')
}

func ExampleSliceContents() {
	data1 := []interface{}{"jack", "rose", 1}
	fmt.Println(SliceContents(reflect.ValueOf(data1)))

	type people struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	data2 := []people{
		{Name: "李雷", Age: 20},
		{Name: "韩梅梅", Age: 19},
	}
	fmt.Println(SliceContents(reflect.ValueOf(data2)))
	data3 := []*people{
		{Name: "李雷", Age: 20},
		{Name: "韩梅梅", Age: 19},
	}
	fmt.Println(SliceContents(reflect.ValueOf(data3)))
	// Output:
	// 'jack','rose',1
	// '{"name":"李雷","age":20}','{"name":"韩梅梅","age":19}'
	// '{"name":"李雷","age":20}','{"name":"韩梅梅","age":19}'
}
