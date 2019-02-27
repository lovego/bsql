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

// Q quote a string, removing all zero byte('\000') in it.
func Q(s string) string {
	s = strings.Replace(s, "'", "''", -1)
	s = strings.Replace(s, "\000", "", -1)
	return "'" + s + "'"
}

func V(i interface{}) string {
	// special types
	switch v := i.(type) {
	case []byte:
		return string(v)
	case time.Time:
		return "'" + v.Format(time.RFC3339Nano) + "'"
	case driver.Valuer:
		return valuer(v)
	case nil:
		return "NULL"
	}

	// basic types: use kind to handle type redefine
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.String:
		return Q(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		if v.Bool() {
			return "true"
		} else {
			return "false"
		}
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'G', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'G', -1, 64)
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return "NULL"
		} else {
			return V(v.Elem().Interface())
		}
	}

	// other types: use json
	b, err := json.Marshal(i)
	if err != nil {
		log.Panic("bsql json.Marshal: ", err)
	}
	return Q(string(b))
}

// Array return data in postgres array form.
func Array(data interface{}) string {
	v, err := pq.Array(data).Value()
	if err != nil {
		log.Panic("bsql Array: ", err)
	}
	if v == nil {
		return "'{}'"
	}
	return Q(v.(string))
}

func Json(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		log.Panic("bsql JsonArray: ", err)
	}
	return Q(string(b))
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
			return Q(s)
		}
	default:
		return V(ifc)
	}
}
