package bsql

import (
	"reflect"
	"strings"
)

// Values return the contents following the sql keyword "VALUES"
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
		return "(" + V(data) + ")"
	}
}

func MapKeyValues(data interface{}) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Map:
		switch value.Type().Key().Kind() {
		case reflect.Slice, reflect.Array:
			var slice []string
			for _, key := range value.MapKeys() {
				slice = append(slice, "("+SliceContents(key)+")")
			}
			return strings.Join(slice, ",")
		default:
			var slice []string
			for _, key := range value.MapKeys() {
				slice = append(slice, V(key.Interface()))
			}
			return "(" + strings.Join(slice, ",") + ")"
		}
	default:
		return Json(data)
	}
}

func SingleColumnValues(data interface{}) string {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		var slice []string
		for i := 0; i < value.Len(); i++ {
			slice = append(slice, "("+V(value.Index(i).Interface())+")")
		}
		return strings.Join(slice, ",")
	default:
		return "(" + V(data) + ")"
	}
}

func SliceContents(value reflect.Value) string {
	var slice []string
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, V(value.Index(i).Interface()))
	}
	return strings.Join(slice, ",")
}
