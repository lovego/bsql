package bsql

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/errs"
)

func run(
	debug bool, debugOutput io.Writer, putSqlInError bool,
	sql string, args []interface{}, work func() (time.Time, error),
) error {
	if !debug {
		_, err := work()
		return WrapError(err, sql, putSqlInError)
	}
	var startAt = time.Now()
	var scanAt, err = work()

	var s = make([]string, 0, 5)
	s = append(s, fmt.Sprintf("bsql: total(%s)", time.Since(startAt)))
	if !scanAt.IsZero() {
		s = append(s, fmt.Sprintf("scan(%s)", time.Since(scanAt)))
	}
	s = append(s, color.GreenString(sql))

	if len(args) > 0 {
		s = append(s, "\n", color.BlueString(argsToString(args)))
	}
	if debugOutput == nil {
		debugOutput = os.Stderr
	}
	fmt.Fprintln(debugOutput, strings.Join(s, " "))

	return WrapError(err, sql, putSqlInError)
}

func WrapError(err error, sql string, fullSql bool) error {
	if err == nil {
		return nil
	}
	erro := errs.Trace(err).(*errs.Error)
	if fullSql {
		erro.SetData(GetPosition(err, sql) + "\n" + sql)
	} else {
		erro.SetData(GetPosition(err, sql))
	}
	return erro
}

func argsToString(args []interface{}) string {
	var s = make([]string, len(args))
	for _, arg := range args {
		s = append(s, fmt.Sprintf("%#v", arg))
	}
	return strings.Join(s, " ")
}
