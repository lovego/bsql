package scan

import (
	"fmt"
	"reflect"
)

type A struct {
	*B
}

type B struct {
	C string
}

func ExampleFieldByName() {
	var a = A{}
	FieldByName(reflect.ValueOf(&a).Elem(), "C").Set(reflect.ValueOf("ok"))
	fmt.Println(a.C)
	// Output:
	// ok
}

func ExampleValue_FieldByName() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	fmt.Println(reflect.ValueOf(&A{}).Elem().FieldByName("C"))

	// Output:
	// reflect: indirection through nil pointer to embedded struct
}
