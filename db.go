package bsql

import (
	"context"
	"database/sql"
	"time"

	"github.com/lovego/bsql/scan"
	"github.com/lovego/errs"
	"github.com/lovego/tracer"
)

type DB struct {
	db      *sql.DB
	timeout time.Duration
	FullSql bool // put full sql into error.
}

func New(db *sql.DB, timeout time.Duration) *DB {
	if timeout <= 0 {
		timeout = time.Minute
	}
	return &DB{db, timeout, true}
}

func (db *DB) GetDB() *sql.DB {
	return db.db
}

func (db *DB) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		db.timeout = timeout
	}
}

func (db *DB) Query(data interface{}, sql string, args ...interface{}) error {
	return db.QueryT(db.timeout, data, sql, args...)
}

func (db *DB) QueryT(duration time.Duration, data interface{}, sql string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	return db.query(ctx, data, sql, args)
}

func (db *DB) QueryCtx(ctx context.Context, opName string,
	data interface{}, sql string, args ...interface{},
) error {
	defer tracer.Finish(tracer.StartChild(ctx, opName))
	if ctx.Done() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, db.timeout)
		defer cancel()
	}
	return db.query(ctx, data, sql, args)
}

func (db *DB) query(ctx context.Context, data interface{}, sql string, args []interface{}) error {
	if debug {
		debugSql(sql, args)
	}
	rows, err := db.db.QueryContext(ctx, sql, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return WrapError(err, sql, db.FullSql)
	}
	if err := scan.Scan(rows, data); err != nil {
		return errs.Trace(err)
	}
	return nil
}

func (db *DB) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return db.ExecT(db.timeout, sql, args...)
}

func (db *DB) ExecT(duration time.Duration, sql string, args ...interface{}) (sql.Result, error) {
	if debug {
		debugSql(sql, args)
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	result, err := db.db.ExecContext(ctx, sql, args...)
	if err != nil {
		err = WrapError(err, sql, db.FullSql)
	}
	return result, err
}

func (db *DB) ExecCtx(ctx context.Context, opName string, sql string, args ...interface{}) (sql.Result, error) {
	defer tracer.Finish(tracer.StartChild(ctx, opName))
	if ctx.Done() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, db.timeout)
		defer cancel()
	}
	if debug {
		debugSql(sql, args)
	}
	result, err := db.db.ExecContext(ctx, sql, args...)
	if err != nil {
		err = WrapError(err, sql, db.FullSql)
	}
	return result, err
}

func (db *DB) RunInTransaction(fn func(*Tx) error) error {
	return db.RunInTransactionT(db.timeout, fn)
}

func (db *DB) RunInTransactionT(duration time.Duration, fn func(*Tx) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.Trace(err)
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()
	if err := fn(&Tx{tx, db.timeout, db.FullSql}); err != nil {
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
		ctx, cancel = context.WithTimeout(ctx, db.timeout)
		defer cancel()
	}

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.Trace(err)
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()
	if err := fn(&Tx{tx, db.timeout, db.FullSql}, ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return errs.Trace(err)
	}
	return nil
}
