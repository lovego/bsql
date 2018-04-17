package psql

func Scan(rows *sql.Rows, data interface{}) error {
	p := reflect.ValueOf(data)
	if p.Kind() != reflect.Ptr {
		return errors.New("data must be a pointer.")
	}
	target := p.Elem()
	switch target.Kind() {
	case reflect.Struct:
		fieldNames := Columns2Fields(rows.Columns())
		if rows.Next() {
			if err := rows.Scan(StructFieldsAddrs(target, fieldNames)...); err != nil {
				return err
			}
		}
	case reflect.Slice:
		elemType := target.Type().Elem()
		fieldNames := Columns2Fields(rows.Columns())
		for rows.Next() {
			elemValue := reflect.Zero(elemType)
			if err := rows.Scan(StructFieldsAddrs(elemValue, fieldNames)...); err != nil {
				return err
			}
			target = reflect.Append(target, elemValue)
		}
		p.Elem().Set(target)
	default:
		if rows.Next() {
			if err := rows.Scan(p); err != nil {
				return err
			}
		}
	}
	return rows.Err()
}
