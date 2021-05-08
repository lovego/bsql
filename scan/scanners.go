package scan

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

// when JSON/JSONB column use jsonScanner, when ARRAY column use pq.Array,
// otherwise use basicScanner.
// Because sql.Rows.Scan's builtin logic can't scan nil to int/string,
// so we always return a sql.Scanner to avoid its builtin logic.
func scannerOf(dest reflect.Value, column ColumnType) sql.Scanner {
	addr := dest.Addr()
	if scanner := trySqlScanner(addr.Interface()); scanner != nil {
		return scanner
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
			return pq.Array(addr.Interface())
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
	case nil:
		// if src is null, should set dest to it's zero value.
		// eg. when dest is int, should set it to 0.
		js.dest.Set(reflect.Zero(js.dest.Type()))
		return nil
	case []byte:
		return json.Unmarshal(buf, getJsonDest(js.dest))
	case string:
		return json.Unmarshal([]byte(buf), getJsonDest(js.dest))
	default:
		return fmt.Errorf("bsql jsonScanner unexpected src: %T(%v)", src, src)
	}
}

func getJsonDest(dest reflect.Value) interface{} {
	if dest.Kind() == reflect.Interface && !dest.IsNil() {
		return dest.Elem().Interface()
	}
	return dest.Addr().Interface()
}

func trySqlScanner(ptr interface{}) sql.Scanner {
	if scanner, ok := ptr.(sql.Scanner); ok {
		switch v := ptr.(type) {
		case *decimal.Decimal:
			return decimalScanner{v}
		}
		return scanner
	}
	return nil
}

type decimalScanner struct {
	d *decimal.Decimal
}

func (ds decimalScanner) Scan(src interface{}) error {
	if src == nil {
		*ds.d = decimal.Zero
		return nil
	}
	return ds.d.Scan(src)
}
