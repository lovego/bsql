package bsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/lib/pq"
)

type jsonScanner struct {
	dest interface{}
}
type basicScanner struct {
	dest interface{}
}

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
			// return addr
			return basicScanner{addr}
		}
	}
}

func (js jsonScanner) Scan(src interface{}) error {
	switch buf := src.(type) {
	case string:
		return json.Unmarshal([]byte(buf), js.dest)
	case []byte:
		return json.Unmarshal(buf, js.dest)
	case nil:
		return basicScanner{js.dest}.Scan(src)
	default:
		return fmt.Errorf("bsql jsonScanner unexpected: %T %v", src, src)
	}
}

func (bs basicScanner) Scan(src interface{}) error {
	switch d := bs.dest.(type) {
	case *string:
		return scanString(d, src)
	case *bool:
		return scanBool(d, src)
	case *int:
		return scanInt(d, src)
	case *int8:
		return scanInt8(d, src)
	case *int16:
		return scanInt16(d, src)
	case *int32:
		return scanInt32(d, src)
	case *int64:
		return scanInt64(d, src)
	case *uint:
		return scanUint(d, src)
	case *uint8:
		return scanUint8(d, src)
	case *uint16:
		return scanUint16(d, src)
	case *uint32:
		return scanUint32(d, src)
	case *uint64:
		return scanUint64(d, src)
	case *float32:
		return scanFloat32(d, src)
	case *float64:
		return scanFloat64(d, src)
	case *[]byte:
		return scanBytes(d, src)
	case *time.Time:
		return scanTime(d, src)
	default:
		return fmt.Errorf("bsql: unsupported dest type: %T", bs.dest)
	}
	return nil
}

/*
The src value will be of one of the following types:
   int64
   float64
   bool
   []byte
   string
   time.Time
   nil - for NULL values
*/
func scanBool(d *bool, src interface{}) error {
	switch s := src.(type) {
	case bool:
		*d = s
	case nil:
		*d = false
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to bool", src, src)
	}
	return nil
}

func scanBytes(d *[]byte, src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*d = s
	case string:
		*d = []byte(s)
	case nil:
		*d = nil
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to []byte", src, src)
	}
	return nil
}

func scanString(d *string, src interface{}) error {
	switch s := src.(type) {
	case string:
		*d = s
	case []byte:
		*d = string(s)
	case nil:
		*d = ""
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to string", src, src)
	}
	return nil
}

func scanTime(d *time.Time, src interface{}) error {
	switch s := src.(type) {
	case time.Time:
		*d = s
	case nil:
		*d = time.Time{}
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to time.Time", src, src)
	}
	return nil
}
