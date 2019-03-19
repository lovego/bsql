package scan

import (
	"fmt"
	"log"
	"reflect"
)

func ExampleJsonScanner() {
	var m map[string]int
	js := jsonScanner{dest: reflect.ValueOf(&m).Elem()}
	if err := js.Scan(nil); err != nil {
		log.Panic(err)
	}
	fmt.Println(m == nil)

	m = map[string]int{"key": 1}
	if err := js.Scan(nil); err != nil {
		log.Panic(err)
	}
	fmt.Println(m == nil)

	// Output:
	// true
	// true
}

func ExampleJsonScanner_interface() {
	var ifc interface{}
	js := jsonScanner{dest: reflect.ValueOf(&ifc).Elem()}

	if err := js.Scan("123"); err != nil {
		log.Panic(err)
	}
	fmt.Println(ifc, reflect.TypeOf(ifc))

	var i uint
	ifc = &i

	if err := js.Scan("123"); err != nil {
		log.Panic(err)
	}
	value := *ifc.(*uint)
	fmt.Println(value, reflect.TypeOf(value))

	// Output:
	// 123 float64
	// 123 uint
}
