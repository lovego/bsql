package bsql

import (
	"context"
	"database/sql"
	"time"

	"github.com/lovego/bsql/scan"
	"github.com/lovego/errs"
	"github.com/lovego/tracer"
)

type Tx struct {
	tx      *sql.Tx
	timeout time.Duration
	FullSql bool // put full sql in error.
}

func (tx *Tx) Query(data interface{}, sql string, args ...interface{}) error {
	return tx.QueryT(tx.timeout, data, sql, args...)
}

func (tx *Tx) QueryT(duration time.Duration,
	data interface{}, sql string, args ...interface{},
) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	return tx.query(ctx, data, sql, args)
}

func (tx *Tx) QueryCtx(ctx context.Context, opName string,
	data interface{}, sql string, args ...interface{},
) error {
	defer tracer.Finish(tracer.StartChild(ctx, opName))
	if ctx.Done() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, tx.timeout)
		defer cancel()
	}
	return tx.query(ctx, data, sql, args)
}

func (tx *Tx) query(ctx context.Context, data interface{}, sql string, args []interface{}) error {
	if debug {
		debugSql(sql, args)
	}
	rows, err := tx.tx.QueryContext(ctx, sql, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return errs.Trace(ErrorWithPosition(err, sql, tx.FullSql))
	}
	if err := scan.Scan(rows, data); err != nil {
		return errs.Trace(err)
	}
	return nil
}

func (tx *Tx) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return tx.ExecT(tx.timeout, sql, args...)
}

func (tx *Tx) ExecT(
	duration time.Duration, sql string, args ...interface{},
) (sql.Result, error) {
	if debug {
		debugSql(sql, args)
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	result, err := tx.tx.ExecContext(ctx, sql, args...)
	if err != nil {
		err = errs.Trace(ErrorWithPosition(err, sql, tx.FullSql))
	}
	return result, err
}

func (tx *Tx) ExecCtx(
	ctx context.Context, opName string, sql string, args ...interface{},
) (sql.Result, error) {
	defer tracer.Finish(tracer.StartChild(ctx, opName))
	if ctx.Done() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, tx.timeout)
		defer cancel()
	}
	if debug {
		debugSql(sql, args)
	}
	result, err := tx.tx.ExecContext(ctx, sql, args...)
	if err != nil {
		err = errs.Trace(ErrorWithPosition(err, sql, tx.FullSql))
	}
	return result, err
}

func (tx *Tx) GetTx() *sql.Tx {
	return tx.tx
}
