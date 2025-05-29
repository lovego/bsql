package bsql

import (
	"context"
	"database/sql"
	"io"
	"time"

	"github.com/lovego/bsql/scan"
	"github.com/lovego/errs"
	"github.com/lovego/tracer"
)

type Tx struct {
	Tx            *sql.Tx
	Context       context.Context
	Timeout       time.Duration // default timeout for Query or Exec.
	Debug         bool
	DebugOutput   io.Writer
	PutSqlInError bool // put sql into returned error if error happend .
}

func NewTx(tx *sql.Tx, timeout time.Duration) *Tx {
	if timeout <= 0 {
		timeout = time.Minute
	}
	return &Tx{Tx: tx, Timeout: timeout, PutSqlInError: true}
}

func (tx *Tx) Query(data interface{}, sql string, args ...interface{}) error {
	return tx.QueryT(tx.Timeout, data, sql, args...)
}

// query and return to data when do exec and returning
func (tx *Tx) QueryR(data interface{}, sql string, args ...interface{}) error {
	ctx, cancel := tx.context(tx.Timeout)
	defer cancel()
	return tx.query(ctx, data, sql, args, true)
}

func (tx *Tx) QueryT(duration time.Duration, data interface{}, sql string, args ...interface{}) error {
	ctx, cancel := tx.context(duration)
	defer cancel()
	return tx.query(ctx, data, sql, args)
}

func (tx *Tx) QueryCtx(ctx context.Context, opName string,
	data interface{}, sql string, args ...interface{},
) error {
	defer tracer.Finish(tracer.StartChild(ctx, opName))
	if ctx.Done() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, tx.Timeout)
		defer cancel()
	}
	return tx.query(ctx, data, sql, args)
}

func (tx *Tx) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return tx.ExecT(tx.Timeout, sql, args...)
}

func (tx *Tx) ExecT(duration time.Duration, sql string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := tx.context(duration)
	defer cancel()
	return tx.exec(ctx, sql, args)
}

func (tx *Tx) ExecCtx(ctx context.Context, opName string, sql string, args ...interface{}) (sql.Result, error) {
	defer tracer.Finish(tracer.StartChild(ctx, opName))
	if ctx.Done() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, tx.Timeout)
		defer cancel()
	}
	return tx.exec(ctx, sql, args)
}

func (tx *Tx) query(ctx context.Context, data interface{}, sql string, args []interface{}, append ...bool) error {
	return run(tx.Debug, tx.DebugOutput, tx.PutSqlInError, sql, args,
		func() (scanAt time.Time, err error) {
			rows, err := tx.Tx.QueryContext(ctx, sql, args...)
			if rows != nil {
				defer rows.Close()
			}
			if err != nil {
				return scanAt, errs.Trace(err)
			}
			if tx.Debug {
				scanAt = time.Now()
			}
			if err := scan.Scan(rows, data, append...); err != nil {
				return scanAt, errs.Trace(err)
			}
			return
		})
}

func (tx *Tx) exec(
	ctx context.Context, sql string, args []interface{},
) (result sql.Result, err error) {
	err = run(tx.Debug, tx.DebugOutput, tx.PutSqlInError, sql, args,
		func() (time.Time, error) {
			result, err = tx.Tx.ExecContext(ctx, sql, args...)
			return time.Time{}, errs.Trace(err)
		})
	return
}

func (tx *Tx) context(timeout time.Duration) (context.Context, context.CancelFunc) {
	var ctx = tx.Context
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, timeout)
}
