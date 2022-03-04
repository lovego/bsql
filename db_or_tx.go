package bsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/lovego/bsql/position"
)

type DbOrTx interface {
	Query(data interface{}, sql string, args ...interface{}) error
	QueryCtx(ctx context.Context, opName string, data interface{}, sql string, args ...interface{}) error
	QueryT(duration time.Duration, data interface{}, sql string, args ...interface{}) error
	Exec(sql string, args ...interface{}) (sql.Result, error)
	ExecT(duration time.Duration, sql string, args ...interface{}) (sql.Result, error)
}

func IsNil(dbOrTx DbOrTx) bool {
	if dbOrTx == nil {
		return true
	}
	if tx, ok := dbOrTx.(*Tx); ok && tx == nil {
		return true
	}
	if db, ok := dbOrTx.(*DB); ok && db == nil {
		return true
	}
	return false
}

func ErrorWithPosition(err error, sql string) error {
	if positionDesc := GetPosition(err, sql); positionDesc != "" {
		return errors.New(err.Error() + "\n" + positionDesc)
	}
	return err
}

func GetPosition(err error, sql string) string {
	pqError, ok := err.(*pq.Error)
	if !ok || pqError == nil {
		return ""
	}
	// Position: the field value is a decimal ASCII integer,
	// indicating an error cursor position as an index into the original query string.
	// The first character has index 1, and positions are measured in characters not bytes.
	pos := pqError.Position
	if pos == "" && pqError.InternalPosition != "" {
		pos = pqError.InternalPosition
		sql = pqError.InternalQuery
	}
	var positionDesc string
	if pos != "" {
		if offset, err := strconv.Atoi(pos); err == nil && offset >= 1 {
			positionDesc = position.Get([]rune(sql), int(offset-1))
		}
	}
	if positionDesc != "" {
		return positionDesc + "\n" + PrettyPrint(*pqError)
	}
	return PrettyPrint(*pqError)
}

func PrettyPrint(v interface{}) string {
	val := reflect.ValueOf(v)
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return fmt.Sprint(v)
	}
	var b strings.Builder
	b.WriteString("{\n")
	for i := 0; i < typ.NumField(); i++ {
		v := fmt.Sprint(val.FieldByIndex([]int{i}).Interface())
		b.WriteString("  " + typ.Field(i).Name + ": " + v + "\n")
	}
	b.WriteString("}")
	return b.String()
}
