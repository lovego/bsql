package bsql

import (
	"fmt"
	"time"

	"github.com/lovego/date"
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

// special types
func ExampleV_nil() {
	var p *int
	var m map[string]int
	fmt.Println(V(nil), V(p), V(m))
	// Output: NULL NULL 'null'
}

func ExampleV_bytes() {
	var b = []byte("abc")
	fmt.Println(V(b), V(&b))
	// Output: abc abc
}

func ExampleV_time() {
	var t time.Time
	fmt.Println(V(t), V(&t))
	t, err := time.Parse(time.RFC3339Nano, "2019-06-19T13:52:08.123456789+08:00")
	fmt.Println(V(t), V(&t), err)
	// Output:
	// '0001-01-01T00:00:00Z' '0001-01-01T00:00:00Z'
	// '2019-06-19T13:52:08.123456+08:00' '2019-06-19T13:52:08.123456+08:00' <nil>
}

func ExampleV_driverValuer_decimal() {
	var d = decimal.New(1234, -2)
	fmt.Println(V(d), V(&d))
	var p *decimal.Decimal
	fmt.Println(V(p))
	// Output:
	// 12.34 12.34
	// NULL
}

func ExampleV_driverValuer_date() {
	var d = date.Date{}
	fmt.Println(V(d), V(&d))

	var d2, _ = date.New("2019-01-01")
	fmt.Println(V(d2), V(*d2))

	// Output:
	// NULL NULL
	// '2019-01-01' '2019-01-01'
}

// basic types
func ExampleV_string() {
	var s = "string"
	fmt.Println(V(s), V(&s))
	// Output: 'string' 'string'
}

func ExampleV_int() {
	var i = -1234567890
	var p = &i
	fmt.Println(V(i), V(p), V(&p))
	// Output: -1234567890 -1234567890 -1234567890
}

func ExampleV_uint() {
	var i uint = 1234567890
	fmt.Println(V(i), V(&i))
	// Output: 1234567890 1234567890
}

func ExampleV_bool() {
	var t, f = true, false
	fmt.Println(V(t), V(f), V(&t), V(&f))
	// Output: true false true false
}

func ExampleV_float32() {
	var f = 1.234
	fmt.Println(V(f), V(&f))
	// Output: 1.234 1.234
}

func ExampleV_float64() {
	var f float64 = 1.234567
	fmt.Println(V(f), V(&f))
	// Output: 1.234567 1.234567
}

func ExampleV_json() {
	var j = map[int]bool{2: true, 3: false}
	fmt.Println(V(j), V(&j))
	// Output: '{"2":true,"3":false}' '{"2":true,"3":false}'
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

func ExampleJson() {
	fmt.Println(Json([]int{1, 2, 3}))
	fmt.Println(Json([]string{"a", "b", "c"}))
	fmt.Println(Json([][]interface{}{
		{1, "a", true}, {2, "b", true}, {3, "c", false}, {4, "dd'ee", false},
	}))
	// Output:
	// '[1,2,3]'
	// '["a","b","c"]'
	// '[[1,"a",true],[2,"b",true],[3,"c",false],[4,"dd''ee",false]]'
}
