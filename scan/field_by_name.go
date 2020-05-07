package scan

import "reflect"

func FieldByName(v reflect.Value, name string) reflect.Value {
	if f, ok := v.Type().FieldByName(name); ok {
		return FieldByIndex(v, f.Index)
	}
	return reflect.Value{}
}

func FieldByIndex(v reflect.Value, index []int) reflect.Value {
	if len(index) == 1 {
		return v.Field(index[0])
	}
	for i, x := range index {
		if i > 0 {
			if v.Kind() == reflect.Ptr && v.Type().Elem().Kind() == reflect.Struct {
				if v.IsNil() {
					v.Set(reflect.New(v.Type().Elem()))
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}
	return v
}
