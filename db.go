package bsql

import (
	"context"
	"database/sql"
	"io"
	"os"
	"time"

	"github.com/lovego/bsql/scan"
	"github.com/lovego/errs"
	"github.com/lovego/tracer"
)

type DB struct {
	DB            *sql.DB
	Context       context.Context
	Timeout       time.Duration // default timeout for Query, Exec and transactions.
	Debug         bool
	DebugOutput   io.Writer
	PutSqlInError bool // put sql into returned error if error happend .
}

func New(db *sql.DB, timeout time.Duration) *DB {
	if timeout <= 0 {
		timeout = time.Minute
	}
	return &DB{
		DB:            db,
		Timeout:       timeout,
		Debug:         os.Getenv(`DebugBsql`) != ``,
		PutSqlInError: true,
	}
}

func (db *DB) GetDB() *sql.DB {
	return db.DB
}

func (db *DB) Query(data interface{}, sql string, args ...interface{}) error {
	return db.QueryT(db.Timeout, data, sql, args...)
}

func (db *DB) QueryT(duration time.Duration, data interface{}, sql string, args ...interface{}) error {
	ctx, cancel := db.context(duration)
	defer cancel()
	return db.query(ctx, data, sql, args)
}

func (db *DB) QueryCtx(ctx context.Context, opName string,
	data interface{}, sql string, args ...interface{},
) error {
	defer tracer.Finish(tracer.StartChild(ctx, opName))
	if ctx.Done() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, db.Timeout)
		defer cancel()
	}
	return db.query(ctx, data, sql, args)
}

func (db *DB) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return db.ExecT(db.Timeout, sql, args...)
}

func (db *DB) ExecT(duration time.Duration, sql string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := db.context(duration)
	defer cancel()
	return db.exec(ctx, sql, args)
}

func (db *DB) ExecCtx(
	ctx context.Context, opName string, sql string, args ...interface{},
) (sql.Result, error) {
	defer tracer.Finish(tracer.StartChild(ctx, opName))
	if ctx.Done() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, db.Timeout)
		defer cancel()
	}
	return db.exec(ctx, sql, args)
}

func (db *DB) query(ctx context.Context, data interface{}, sql string, args []interface{}) error {
	return run(db.Debug, db.DebugOutput, db.PutSqlInError, sql, args,
		func() (scanAt time.Time, err error) {
			rows, err := db.DB.QueryContext(ctx, sql, args...)
			if rows != nil {
				defer rows.Close()
			}
			if err != nil {
				return scanAt, errs.Trace(err)
			}
			if db.Debug {
				scanAt = time.Now()
			}
			if err := scan.Scan(rows, data); err != nil {
				return scanAt, errs.Trace(err)
			}
			return scanAt, nil
		})
}

func (db *DB) exec(
	ctx context.Context, sql string, args []interface{},
) (result sql.Result, err error) {
	err = run(db.Debug, db.DebugOutput, db.PutSqlInError, sql, args,
		func() (time.Time, error) {
			result, err = db.DB.ExecContext(ctx, sql, args...)
			return time.Time{}, errs.Trace(err)
		})
	return
}

func (db *DB) RunInTransaction(fn func(*Tx) error) error {
	return db.RunInTransactionT(db.Timeout, fn)
}

func (db *DB) RunInTransactionT(duration time.Duration, fn func(*Tx) error) error {
	ctx, cancel := db.context(duration)
	defer cancel()

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return errs.Trace(err)
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()
	if err := fn(&Tx{
		Tx: tx, Context: db.Context, Timeout: db.Timeout, PutSqlInError: db.PutSqlInError,
	}); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return errs.Trace(err)
	}
	return nil
}

func (db *DB) RunInTransactionCtx(
	ctx context.Context, opName string, fn func(*Tx, context.Context) error,
) error {
	ctx = tracer.StartChild(ctx, opName)
	defer tracer.Finish(ctx)

	if ctx.Done() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, db.Timeout)
		defer cancel()
	}

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return errs.Trace(err)
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()
	if err := fn(&Tx{
		Tx: tx, Context: db.Context, Timeout: db.Timeout, PutSqlInError: db.PutSqlInError,
	}, ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return errs.Trace(err)
	}
	return nil
}

func (db *DB) context(timeout time.Duration) (context.Context, context.CancelFunc) {
	var ctx = db.Context
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, timeout)
}
