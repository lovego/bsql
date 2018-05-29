package bsql

import (
	"database/sql/driver"
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

// Quote string
func Q(q string) string {
	return "'" + strings.Replace(q, "'", "''", -1) + "'"
}

func Array(data interface{}) string {
	v, err := pq.Array(data).Value()
	if err != nil {
		log.Panic("bsql Array: ", err)
	}
	return "'" + strings.Replace(v.(string), "'", "''", -1) + "'"
}

func JsonArray(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		log.Panic("bsql JsonArray: ", err)
	}
	return "'" + strings.Replace(string(b), "'", "''", -1) + "'"
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
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		log.Panic("bsql: data must be struct or struct slice.")
	}
	var slice []string
	for _, fieldName := range fields {
		field := value.FieldByName(fieldName)
		if !field.IsValid() {
			log.Panic("bsql: no field '" + fieldName + "' in struct")
		}
		slice = append(slice, V(field.Interface()))
	}
	return "(" + strings.Join(slice, ",") + ")"
}

func V(i interface{}) string {
	switch v := i.(type) {
	case string:
		return Q(v)
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
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
	case float32:
		return strconv.FormatFloat(float64(v), 'G', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'G', -1, 64)
	case nil:
		return "NULL"
	case []byte:
		return string(v)
	case time.Time:
		return "'" + v.Format(time.RFC3339Nano) + "'"
	case driver.Valuer:
		return valuer(v)
	default:
		buf, err := json.Marshal(v)
		if err != nil {
			log.Panic("bsql json.Marshal: ", err)
		}
		return Q(string(buf))
	}
}

func valuer(v driver.Valuer) string {
	ifc, err := v.Value()
	if err != nil {
		log.Panic("bsql valuer: ", err)
	}
	switch s := ifc.(type) {
	case string:
		if _, err := strconv.ParseFloat(s, 64); err == nil {
			return s
		} else {
			return "'" + s + "'"
		}
	default:
		return V(ifc)
	}
}
