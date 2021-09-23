package bsql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/lib/pq"
	"github.com/lovego/bsql/position"
	"github.com/lovego/errs"
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

var debug = os.Getenv(`DebugBsql`) != ``

func debugSql(sql string, args []interface{}) {
	color.Green(sql)
	argsString := ``
	for _, arg := range args {
		argsString += fmt.Sprintf("%#v ", arg)
	}
	color.Blue(argsString)
}

func WrapError(err error, sql string, fullSql bool) error {
	if err == nil {
		return nil
	}
	var positionDesc string
	if pqError, ok := err.(*pq.Error); ok {
		// Position: the field value is a decimal ASCII integer,
		// indicating an error cursor position as an index into the original query string.
		// The first character has index 1, and positions are measured in characters not bytes.
		if offset, err := strconv.Atoi(pqError.Position); err == nil && offset >= 1 {
			positionDesc = position.Get([]rune(sql), int(offset-1))
			if positionDesc == "" {
				positionDesc = fmt.Sprintf("(Position: %s)", pqError.Position)
			}
		}
	}

	erro := errs.Trace(err).(*errs.Error)
	if fullSql {
		erro.SetData(positionDesc + "\n" + sql)
	} else {
		erro.SetData(positionDesc)
	}
	return erro
}
