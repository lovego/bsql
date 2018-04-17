package bsql

import (
	//	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Quote string
func Q(q string) string {
	return "'" + strings.Replace(q, "'", "''", -1) + "'"
}

func StructValues(data interface{}, fields []string) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Slice:
		var slice []string
		for i := 0; i < value.Len(); i++ {
			slice = append(slice, StructValuesIn(value.Index(i), fields))
		}
		return strings.Join(slice, ",")
	case reflect.Struct:
		return StructValuesIn(value, fields)
	default:
		return ""
	}
}

func StructValuesIn(value reflect.Value, fields []string) string {
	var slice []string
	for _, fieldName := range fields {
		slice = append(slice, V(value.FieldByName(fieldName).Interface()))
	}
	return "(" + strings.Join(slice, ",") + ")"
}

func V(i interface{}) string {
	switch v := i.(type) {
	case string:
		return Q(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case bool:
		if v {
			return "t"
		} else {
			return "f"
		}
	case time.Time:
		return v.Format(time.RFC3339Nano)
	case interface {
		String() string
	}:
		return v.String()
	default:
		return ""
	}
}
