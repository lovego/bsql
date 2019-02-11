package scan

import (
	"fmt"
	"reflect"
	"strconv"
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
func (bs *basicScanner) Scan(src interface{}) error {
	switch buf := src.(type) {
	case int64:
		// setInt(src, getRealDest(bs.dest))
	case float64:
	case bool:
	case []byte:
	case string:
	case time.Time:
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
func setInt(src int64, dest reflect.Value) error {
	switch dest.Kind() {
	case reflect.Int64, reflect.Uint64:
		dest.SetInt(src)
	case reflect.Int32:
	case reflect.Uint32:
		dest.SetInt(src)

		if i := int(s); int64(i) == s {
			*d = i
		} else {
			return fmt.Errorf("bsql: cannot assign %T(%v) to int: value out of range", src, src)
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign int(%v) to %T", src, src)
	}
	return nil
}

func scanInt8(d *int8, src interface{}) error {
	switch s := src.(type) {
	case int64:
		if i := int8(s); int64(i) == s {
			*d = i
		} else {
			return fmt.Errorf("bsql: cannot assign %T(%v) to int8: value out of range", src, src)
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to int8", src, src)
	}
	return nil
}

func scanInt16(d *int16, src interface{}) error {
	switch s := src.(type) {
	case int64:
		if i := int16(s); int64(i) == s {
			*d = i
		} else {
			return fmt.Errorf("bsql: cannot assign %T(%v) to int16: value out of range", src, src)
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to int16", src, src)
	}
	return nil
}

func scanInt32(d *int32, src interface{}) error {
	switch s := src.(type) {
	case int64:
		if i := int32(s); int64(i) == s {
			*d = i
		} else {
			return fmt.Errorf("bsql: cannot assign %T(%v) to int32: value out of range", src, src)
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to int32", src, src)
	}
	return nil
}

func scanInt64(d *int64, src interface{}) error {
	switch s := src.(type) {
	case int64:
		*d = s
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to int64", src, src)
	}
	return nil
}

func scanUint(d *uint, src interface{}) error {
	switch s := src.(type) {
	case int64:
		if i := uint(s); int64(i) == s {
			*d = i
		} else {
			return fmt.Errorf("bsql: cannot assign %T(%v) to uint: value out of range", src, src)
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to uint", src, src)
	}
	return nil
}

func scanUint8(d *uint8, src interface{}) error {
	switch s := src.(type) {
	case int64:
		if i := uint8(s); int64(i) == s {
			*d = i
		} else {
			return fmt.Errorf("bsql: cannot assign %T(%v) to uint8: value out of range", src, src)
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to uint8", src, src)
	}
	return nil
}

func scanUint16(d *uint16, src interface{}) error {
	switch s := src.(type) {
	case int64:
		if i := uint16(s); int64(i) == s {
			*d = i
		} else {
			return fmt.Errorf("bsql: cannot assign %T(%v) to uint16: value out of range", src, src)
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to uint16", src, src)
	}
	return nil
}

func scanUint32(d *uint32, src interface{}) error {
	switch s := src.(type) {
	case int64:
		if i := uint32(s); int64(i) == s {
			*d = i
		} else {
			return fmt.Errorf("bsql: cannot assign %T(%v) to uint32: value out of range", src, src)
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to uint32", src, src)
	}
	return nil
}

func scanUint64(d *uint64, src interface{}) error {
	switch s := src.(type) {
	case int64:
		if i := uint64(s); int64(i) == s {
			*d = i
		} else {
			return fmt.Errorf("bsql: cannot assign %T(%v) to uint64: value out of range", src, src)
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to uint64", src, src)
	}
	return nil
}

func scanFloat32(d *float32, src interface{}) error {
	switch s := src.(type) {
	case float64:
		if f := float32(s); float64(f) == s {
			*d = f
		} else {
			return fmt.Errorf("bsql: cannot assign %T(%v) to float32: value out of range", src, src)
		}
	case []byte:
		if f, err := strconv.ParseFloat(string(s), 32); err != nil {
			return err
		} else {
			*d = float32(f)
		}
	case string:
		if f, err := strconv.ParseFloat(s, 32); err != nil {
			return err
		} else {
			*d = float32(f)
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to float32", src, src)
	}
	return nil
}

func scanFloat64(d *float64, src interface{}) error {
	switch s := src.(type) {
	case float64:
		*d = s
	case []byte:
		if f, err := strconv.ParseFloat(string(s), 64); err != nil {
			return err
		} else {
			*d = f
		}
	case string:
		if f, err := strconv.ParseFloat(s, 64); err != nil {
			return err
		} else {
			*d = f
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to float64", src, src)
	}
	return nil
}

func scanBool(d *bool, src interface{}) error {
	switch s := src.(type) {
	case bool:
		*d = s
	case nil:
		*d = false
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to bool", src, src)
	}
	return nil
}

func scanBytes(d *[]byte, src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*d = s
	case string:
		*d = []byte(s)
	case nil:
		*d = nil
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to []byte", src, src)
	}
	return nil
}

func scanString(d *string, src interface{}) error {
	switch s := src.(type) {
	case string:
		*d = s
	case []byte:
		*d = string(s)
	case nil:
		*d = ""
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to string", src, src)
	}
	return nil
}

func scanTime(d *time.Time, src interface{}) error {
	switch s := src.(type) {
	case time.Time:
		*d = s
	case nil:
		*d = time.Time{}
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to time.Time", src, src)
	}
	return nil
}
