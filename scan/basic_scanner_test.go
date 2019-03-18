package scan

import (
	"fmt"
	"reflect"
)

func ExampleGetRealDest() {
	var i int = 99
	dest := getRealDest(reflect.ValueOf(&i))
	fmt.Println(dest, dest.Type())

	var p = &i
	dest = getRealDest(reflect.ValueOf(&p))
	fmt.Println(dest, dest.Type())

	var p2 ***int64
	dest = getRealDest(reflect.ValueOf(&p2))
	fmt.Println(dest, dest.Type())

	// Output:
	// 99 int
	// 99 int
	// 0 int64
}

func ExampleScanInt() {
	var i int
	fmt.Println(scanInt(-100, getRealDest(reflect.ValueOf(&i))), i)
	var u uint
	fmt.Println(scanInt(100, getRealDest(reflect.ValueOf(&u))), u)

	var i64 int64
	fmt.Println(scanInt(-64, getRealDest(reflect.ValueOf(&i64))), i64)
	var u64 uint64
	fmt.Println(scanInt(64, getRealDest(reflect.ValueOf(&u64))), u64)

	var i32 int32
	fmt.Println(scanInt(-32, getRealDest(reflect.ValueOf(&i32))), i32)
	var u32 uint32
	fmt.Println(scanInt(32, getRealDest(reflect.ValueOf(&u32))), u32)

	var i16 int16
	fmt.Println(scanInt(-16, getRealDest(reflect.ValueOf(&i16))), i16)
	var u16 uint16
	fmt.Println(scanInt(16, getRealDest(reflect.ValueOf(&u16))), u16)

	var i8 int8
	fmt.Println(scanInt(-8, getRealDest(reflect.ValueOf(&i8))), i8)
	var u8 uint8
	fmt.Println(scanInt(8, getRealDest(reflect.ValueOf(&u8))), u8)

	// Output:
	// <nil> -100
	// <nil> 100
	// <nil> -64
	// <nil> 64
	// <nil> -32
	// <nil> 32
	// <nil> -16
	// <nil> 16
	// <nil> -8
	// <nil> 8
}

func ExampleScanBytes() {
	var b []byte
	fmt.Println(scanBytes(nil, getRealDest(reflect.ValueOf(&b))), b)
	fmt.Println(scanBytes([]byte{}, getRealDest(reflect.ValueOf(&b))), b)
	fmt.Println(scanBytes([]byte{45, 46, 47}, getRealDest(reflect.ValueOf(&b))), b)
	var f float32
	fmt.Println(scanBytes([]byte("1.23"), getRealDest(reflect.ValueOf(&f))), f)
	// Output:
	// <nil> []
	// <nil> []
	// <nil> [45 46 47]
	// <nil> 1.23
}
