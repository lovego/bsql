package scan

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/lib/pq"
)

// when JSON/JSONB column use jsonScanner, when ARRAY column use pq.Array,
// otherwise use basicScanner.
// Because sql.Rows.Scan's builtin logic can't scan nil to int/string,
// so we always return a sql.Scanner to avoid its builtin logic.
func scannerOf(dest reflect.Value, column columnType) interface{} {
	addr := dest.Addr().Interface()
	if _, ok := addr.(sql.Scanner); ok {
		return addr
	}
	var dbType string
	if column.ColumnType != nil {
		dbType = column.ColumnType.DatabaseTypeName()
	}
	switch dbType {
	case "JSONB", "JSON":
		return &jsonScanner{dest}
	default:
		if len(dbType) > 0 && dbType[0] == '_' {
			return pq.Array(dest.Addr().Interface())
		} else {
			return &basicScanner{dest}
		}
	}
}

type jsonScanner struct {
	dest reflect.Value
}

func (js *jsonScanner) Scan(src interface{}) error {
	switch buf := src.(type) {
	case []byte:
		return json.Unmarshal(buf, js.dest.Addr().Interface())
	case string:
		return json.Unmarshal([]byte(buf), js.dest.Addr().Interface())
	case nil:
		// if src is null, should set dest to it's zero value.
		// eg. when dest is int, should set it to 0.
		js.dest.Set(reflect.Zero(js.dest.Type()))
		return nil
	default:
		return fmt.Errorf("bsql jsonScanner unexpected src: %T(%v)", src, src)
	}
}
