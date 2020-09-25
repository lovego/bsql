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
