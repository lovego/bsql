package bsql

import (
	"fmt"
	"strconv"
)

func scanInt(d *int, src interface{}) error {
	switch s := src.(type) {
	case int64:
		if i := int(s); int64(i) == s {
			*d = i
		} else {
			return fmt.Errorf("bsql: cannot assign %T(%v) to int: value out of range", src, src)
		}
	case nil:
		*d = 0
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to int", src, src)
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
