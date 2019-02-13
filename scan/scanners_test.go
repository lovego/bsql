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
