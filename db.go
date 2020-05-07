package bsql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/bsql/scan"
	"github.com/lovego/tracer"
)

type DB struct {
	db      *sql.DB
	timeout time.Duration
}

type DbOrTx interface {
	Query(data interface{}, sql string, args ...interface{}) error
	QueryT(duration time.Duration, data interface{}, sql string, args ...interface{}) error
	Exec(sql string, args ...interface{}) (sql.Result, error)
	ExecT(duration time.Duration, sql string, args ...interface{}) (sql.Result, error)
}

func New(db *sql.DB, timeout time.Duration) *DB {
	if timeout <= 0 {
		timeout = time.Minute
	}
	return &DB{db, timeout}
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
func (db *DB) query(ctx context.Context,
	data interface{}, sql string, args []interface{},
) error {
	if debug {
		debugSql(sql, args)
	}
	rows, err := db.db.QueryContext(ctx, sql, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return err
	}
	return scan.Scan(rows, data)
}

func (db *DB) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return db.ExecT(db.timeout, sql, args...)
}

func (db *DB) ExecT(duration time.Duration, sql string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	if debug {
		debugSql(sql, args)
	}
	return db.db.ExecContext(ctx, sql, args...)
}

func (db *DB) ExecCtx(ctx context.Context, opName string,
	sql string, args ...interface{}) (sql.Result, error) {
	defer tracer.Finish(tracer.StartChild(ctx, opName))
	if ctx.Done() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, db.timeout)
		defer cancel()
	}
	if debug {
		debugSql(sql, args)
	}
	return db.db.ExecContext(ctx, sql, args...)
}

func (db *DB) RunInTransaction(fn func(*Tx) error) error {
	return db.RunInTransactionT(db.timeout, fn)
}

func (db *DB) RunInTransactionT(duration time.Duration, fn func(*Tx) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()
	if err := fn(&Tx{tx, db.timeout}); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *DB) RunInTransactionCtx(ctx context.Context, opName string, fn func(*Tx, context.Context) error) error {
	ctx = tracer.StartChild(ctx, opName)
	defer tracer.Finish(ctx)

	if ctx.Done() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, db.timeout)
		defer cancel()
	}

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()
	if err := fn(&Tx{tx, db.timeout}, ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *DB) GetDB() *sql.DB {
	return db.db
}

func (db *DB) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		db.timeout = timeout
	}
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
