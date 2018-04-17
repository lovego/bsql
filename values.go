package bsql

import (
	"encoding/json"
	"iconv"
	"reflect"
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
		slice = append(slice, Q(value.FieldByName(fieldName).Interface()))
	}
	return "(" + strings.Join(slice, ",") + ")"
}

func V(i interface{}) string {
	switch v := i.(type) {
	case string:
		return Q(v)
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(int64(v), 10, 64)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(uint64(v), 10, 64)
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

	}
}
