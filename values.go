package bsql

import (
	"database/sql/driver"
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Quote string
func Q(q string) string {
	return "'" + strings.Replace(q, "'", "''", -1) + "'"
}

func Values(data interface{}) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		switch value.Type().Elem().Kind() {
		case reflect.Slice, reflect.Array:
			var slice []string
			for i := 0; i < value.Len(); i++ {
				slice = append(slice, "("+SliceContents(value.Index(i))+")")
			}
			return strings.Join(slice, ",")
		default:
			return "(" + SliceContents(value) + ")"
		}
	default:
		return "(" + V(value.Interface()) + ")"
	}
}

func Array(data interface{}) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		return "{" + SliceContents(value) + "}"
	default:
		return "{" + V(value.Interface()) + "}"
	}
}

func SliceContents(value reflect.Value) string {
	var slice []string
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, V(value.Index(i).Interface()))
	}
	return strings.Join(slice, ",")
}

func StructValues(data interface{}, fields []string) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		var slice []string
		for i := 0; i < value.Len(); i++ {
			slice = append(slice, StructValuesIn(value.Index(i), fields))
		}
		return strings.Join(slice, ",")
	default:
		return StructValuesIn(value, fields)
	}
}

func StructValuesIn(value reflect.Value, fields []string) string {
	if value.Kind() != reflect.Struct {
		log.Panic("bsql: data must be struct or struct slice.")
	}
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
	case []byte:
		return string(v)
	case time.Time:
		return v.Format(time.RFC3339Nano)
	case driver.Valuer:
		ifc, err := v.Value()
		if err != nil {
			log.Panic("bsql json.Marshal: ", err)
		}
		if s, ok := ifc.(string); ok {
			return s
		} else {
			return V(ifc)
		}
	default:
		buf, err := json.Marshal(v)
		if err != nil {
			log.Panic("bsql json.Marshal: ", err)
		}
		return Q(string(buf))
	}
}
