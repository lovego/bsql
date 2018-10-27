package bsql

import (
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
)

func ExampleQ() {
	var data = []string{"bsql", "xi'ao'mei", "love'go", "a\000\000b\000c"}
	for i := range data {
		fmt.Println(Q(data[i]))
	}
	// Output:
	// 'bsql'
	// 'xi''ao''mei'
	// 'love''go'
	// 'abc'
}

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

func ExampleArray() {
	fmt.Println(Array([]int{1, 2, 3}))
	fmt.Println(Array([]string{"a", "b", "c"}))
	fmt.Println(Array([][]interface{}{
		{1, "a", true}, {2, "b", true}, {3, "c", false}, {4, "dd'ee", false},
	}))
	// Output:
	// '{1,2,3}'
	// '{"a","b","c"}'
	// '{{1,"a",true},{2,"b",true},{3,"c",false},{4,"dd''ee",false}}'
}

func ExampleJsonArray() {
	fmt.Println(JsonArray([]int{1, 2, 3}))
	fmt.Println(JsonArray([]string{"a", "b", "c"}))
	fmt.Println(JsonArray([][]interface{}{
		{1, "a", true}, {2, "b", true}, {3, "c", false}, {4, "dd'ee", false},
	}))
	// Output:
	// '[1,2,3]'
	// '["a","b","c"]'
	// '[[1,"a",true],[2,"b",true],[3,"c",false],[4,"dd''ee",false]]'
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

func ExampleStructValues() {
	type student struct {
		Id        int
		Name, Sex string
	}
	data1 := []student{
		{1, "李雷", "男"}, {2, "韩梅梅", "女"},
		{3, "Lili", "女"}, {4, "Lucy", "女"},
	}
	fmt.Println(StructValues(data1, []string{"Id", "Name"}))

	data2 := []*student{
		{1, "李雷", "男"}, {2, "韩梅梅", "女"},
		{3, "Lili", "女"}, {4, "Lucy", "女"},
	}
	fmt.Println(StructValues(data2, []string{"Id", "Name"}))

	data3 := []interface{}{
		student{1, "李雷", "男"}, student{2, "韩梅梅", "女"},
		student{3, "Lili", "女"}, student{4, "Lucy", "女"},
	}
	fmt.Println(StructValues(data3, []string{"Id", "Name"}))
	// Output:
	// (1,'李雷'),(2,'韩梅梅'),(3,'Lili'),(4,'Lucy')
	// (1,'李雷'),(2,'韩梅梅'),(3,'Lili'),(4,'Lucy')
	// (1,'李雷'),(2,'韩梅梅'),(3,'Lili'),(4,'Lucy')
}

func ExampleStructValuesIn() {
	type student struct {
		Id        int
		Name, Sex string
	}
	data1 := student{1, "李雷", "男"}
	fmt.Println(StructValuesIn(reflect.ValueOf(data1), []string{"Id", "Name"}))

	data2 := &student{1, "李雷", "男"}
	fmt.Println(StructValuesIn(reflect.ValueOf(data2), []string{"Id", "Name"}))

	var data3 interface{} = student{1, "李雷", "男"}
	fmt.Println(StructValuesIn(reflect.ValueOf(data3), []string{"Id", "Name"}))
	// Output:
	// (1,'李雷')
	// (1,'李雷')
	// (1,'李雷')
}

func ExampleV() {
	fmt.Println(V(1.23))
	fmt.Println(V(float64(1.23)))
	fmt.Println(V(decimal.New(1230, -2)))
	fmt.Println(V(map[int]bool{2: true, 3: false}))
	// Output:
	// 1.23
	// 1.23
	// 12.3
	// '{"2":true,"3":false}'
}

func ExampleStructField() {
	type T2 struct {
		Name string
	}
	type T struct {
		T2
	}
	v := reflect.ValueOf(T{T2{"name"}})
	fmt.Println(structField(v, "Name").Interface())
	fmt.Println(structField(v, "T2.Name").Interface())
	// Output:
	// name
	// name
}
