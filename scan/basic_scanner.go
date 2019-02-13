package scan

import (
	"fmt"
	"reflect"
	"time"
)

type basicScanner struct {
	dest reflect.Value
}

/*
The src value will be of one of the following types:
   int64
   float64
   bool
   []byte
   string
   time.Time
   nil - for NULL values
*/
func (bs *basicScanner) Scan(srcIfc interface{}) error {
	switch src := srcIfc.(type) {
	case int64:
		return scanInt(src, getRealDest(bs.dest))
	case float64:
		return scanFloat(src, getRealDest(bs.dest))
	case bool:
		return scanBool(src, getRealDest(bs.dest))
	case []byte:
		return scanBytes(src, getRealDest(bs.dest))
	case string:
		return scanString(src, getRealDest(bs.dest))
	case time.Time:
		return scanTime(src, getRealDest(bs.dest))
	case nil:
		// if src is null, should set dest to it's zero value.
		// eg. when dest is int, should set it to 0.
		bs.dest.Set(reflect.Zero(bs.dest.Type()))
		return nil
	default:
		return fmt.Errorf("bsql basicScanner unexpected src: %T(%v)", src, src)
	}
}

// the preceding steps ensured that dest is valid
func getRealDest(dest reflect.Value) reflect.Value {
	for dest.Kind() == reflect.Ptr {
		if dest.IsNil() {
			dest.Set(reflect.New(dest.Type().Elem()))
		}
		dest = dest.Elem()
	}
	return dest
}

// use reflect.Value's  SetXXX methods instead of pointers and type switch,
// because we should set by kind, not type.
func scanInt(src int64, dest reflect.Value) error {
	const (
		minInt32  = -1 << 31
		maxInt32  = 1<<31 - 1
		maxUint32 = 1<<32 - 1

		minInt16  = -1 << 15
		maxInt16  = 1<<15 - 1
		maxUint16 = 1<<16 - 1

		minInt8  = -1 << 7
		maxInt8  = 1<<7 - 1
		maxUint8 = 1<<8 - 1
	)
	switch dest.Kind() {
	case reflect.Int64, reflect.Int: // ignore 32 bit machine
		dest.SetInt(src)
	case reflect.Uint64, reflect.Uint:
		dest.SetUint(uint64(src))
	case reflect.Int32:
		if src < minInt32 || src > maxInt32 {
			return errorValueOutOfRange(src, dest)
		}
		dest.SetInt(src)
	case reflect.Uint32:
		if src < 0 || src > maxUint32 {
			return errorValueOutOfRange(src, dest)
		}
		dest.SetUint(uint64(src))
	case reflect.Int16:
		if src < minInt16 || src > maxInt16 {
			return errorValueOutOfRange(src, dest)
		}
		dest.SetInt(src)
	case reflect.Uint16:
		if src < 0 || src > maxUint16 {
			return errorValueOutOfRange(src, dest)
		}
		dest.SetUint(uint64(src))
	case reflect.Int8:
		if src < minInt8 || src > maxInt8 {
			return errorValueOutOfRange(src, dest)
		}
		dest.SetInt(src)
	case reflect.Uint8:
		if src < 0 || src > maxUint8 {
			return errorValueOutOfRange(src, dest)
		}
		dest.SetUint(uint64(src))
	default:
		return errorCannotAssign(src, dest)
	}
	return nil
}

func scanFloat(src float64, dest reflect.Value) error {
	switch dest.Kind() {
	case reflect.Float64, reflect.Float32:
		dest.SetFloat(src)
	default:
		return errorCannotAssign(src, dest)
	}
	return nil
}

func scanBool(src bool, dest reflect.Value) error {
	switch dest.Kind() {
	case reflect.Bool:
		dest.SetBool(src)
	default:
		return errorCannotAssign(src, dest)
	}
	return nil
}

func scanBytes(src []byte, dest reflect.Value) error {
	switch dest.Kind() {
	case reflect.Slice:
		if dest.Type().Elem().Kind() != reflect.Uint8 {
			return errorCannotAssign(src, dest)
		}
		dest.SetBytes(src)
	case reflect.String:
		dest.SetString(string(src))
	default:
		return errorCannotAssign(src, dest)
	}
	return nil
}

func scanString(src string, dest reflect.Value) error {
	switch dest.Kind() {
	case reflect.String:
		dest.SetString(src)
	case reflect.Slice:
		if dest.Type().Elem().Kind() != reflect.Uint8 {
			return errorCannotAssign(src, dest)
		}
		dest.SetBytes([]byte(src))
	default:
		return errorCannotAssign(src, dest)
	}
	return nil
}

func scanTime(src time.Time, dest reflect.Value) error {
	addr := dest.Addr().Interface()
	if ptr, ok := addr.(*time.Time); ok {
		*ptr = src
	} else {
		return errorCannotAssign(src, dest)
	}
	return nil
}

func errorValueOutOfRange(src interface{}, dest reflect.Value) error {
	return fmt.Errorf("bsql: cannot assign %T(%v) to %v: value out of range", src, src, dest.Type())
}

func errorCannotAssign(src interface{}, dest reflect.Value) error {
	return fmt.Errorf("bsql: cannot assign %T(%v) to %v", src, src, dest.Type())
}
