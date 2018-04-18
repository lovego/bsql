package bsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/lib/pq"
)

func scannerOf(addrValue reflect.Value, column columnType) interface{} {
	addr := addrValue.Interface()
	if _, ok := addr.(sql.Scanner); ok {
		return addr
	}
	var dbType string
	if column.ColumnType != nil {
		dbType = column.ColumnType.DatabaseTypeName()
	}
	switch dbType {
	case "JSONB", "JSON":
		return jsonScanner{addr}
	default:
		if len(dbType) > 0 && dbType[0] == '_' {
			return pq.Array(addr)
		} else {
			return addr
		}
	}
}

type basicScanner struct {
	data interface{}
}

func (s basicScanner) Scan(src interface{}) error {
	return nil
}

type jsonScanner struct {
	data interface{}
}

func (s jsonScanner) Scan(src interface{}) error {
	switch buf := src.(type) {
	case string:
		return json.Unmarshal([]byte(buf), s.data)
	case []byte:
		return json.Unmarshal(buf, s.data)
	default:
		return fmt.Errorf("bsql jsonScanner unexpected: %T %v", src, src)
	}
}
