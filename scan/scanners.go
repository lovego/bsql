package scan

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/lib/pq"
)

// when JSON/JSONB column use jsonScanner, when ARRAY column use pq.Array.
// otherwise return a raw pointer to use the database/sql's builtin scan logic.
// sql.Rows.Scan can't scan nil to int, so we should avoid it.
func scannerOf(destAddr reflect.Value, column columnType) interface{} {
	addr := destAddr.Interface()
	if _, ok := addr.(sql.Scanner); ok {
		return addr
	}
	var dbType string
	if column.ColumnType != nil {
		dbType = column.ColumnType.DatabaseTypeName()
	}
	switch dbType {
	case "JSONB", "JSON":
		return &jsonScanner{destAddr}
	default:
		if len(dbType) > 0 && dbType[0] == '_' {
			return pq.Array(destAddr.Interface())
		} else {
			return destAddr.Interface()
		}
	}
}

type jsonScanner struct {
	destAddr reflect.Value
}

func (js *jsonScanner) Scan(src interface{}) error {
	switch buf := src.(type) {
	case []byte:
		return json.Unmarshal(buf, js.destAddr.Interface())
	case string:
		return json.Unmarshal([]byte(buf), js.destAddr.Interface())
	case nil:
		// if src is null, should set dest to it's zero value.
		// eg. when dest is int, should set it to 0.
		dest := js.destAddr.Elem()
		dest.Set(reflect.Zero(dest.Type()))
		return nil
	default:
		return fmt.Errorf("bsql jsonScanner unexpected src: %T(%v)", src, src)
	}
}
