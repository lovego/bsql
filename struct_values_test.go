package bsql

import (
	"fmt"
	"reflect"
)

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

func ExampleStructValues_map() {
	type student struct {
		Id        int
		Name, Sex string
	}
	m := map[string]bool{"(1,'李雷'),(2,'韩梅梅')": true, "(2,'韩梅梅'),(1,'李雷')": true}
	result := StructValues(map[student]int{
		student{1, "李雷", "男"}:  1,
		student{2, "韩梅梅", "女"}: 2,
	}, []string{"Id", "Name"})

	fmt.Println(m[result])

	// Output:
	// true
}

func ExampleStructFields() {
	type student struct {
		Id        int
		Name, Sex string
	}
	data1 := student{1, "李雷", "男"}
	fmt.Println(StructFields(data1, []string{"Id", "Name"}))

	data2 := &student{1, "李雷", "男"}
	fmt.Println(StructFields(data2, []string{"Id", "Name"}))

	var data3 interface{} = student{1, "李雷", "男"}
	fmt.Println(StructFields(data3, []string{"Id", "Name"}))
	// Output:
	// 1,'李雷'
	// 1,'李雷'
	// 1,'李雷'
}

func ExampleGetValue() {
	type T2 struct {
		Name string
	}
	type T struct {
		T2
	}
	v := reflect.ValueOf(T{T2{"name"}})
	fmt.Println(getValue(v, "Name").Interface())
	fmt.Println(getValue(v, "T2.Name").Interface())
	// Output:
	// name
	// name
}

func ExampleStructFiledValues() {
	type student struct {
		Id        int
		Name, Sex string
	}
	data1 := []student{
		{1, "李雷", "男"}, {2, "韩梅梅", "女"},
		{3, "Lili", "女"}, {4, "Lucy", "女"},
	}
	fmt.Println(StructFieldValues(data1, "Id"))
	fmt.Println(StructFieldValues(data1, "Name"))

	// Output:
	// (1,2,3,4)
	// ('李雷','韩梅梅','Lili','Lucy')
}
