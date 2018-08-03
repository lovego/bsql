package bsql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/tracer"
)

// DB indicates the database to connect to and the timeout.
type DB struct {
	db      *sql.DB
	timeout time.Duration
}

// DbOrTx declares the query and statement execution functions for DB and TX.
type DbOrTx interface {
	Query(data interface{}, sql string, args ...interface{}) error
	QueryT(duration time.Duration, data interface{}, sql string, args ...interface{}) error
	Exec(sql string, args ...interface{}) (sql.Result, error)
	ExecT(duration time.Duration, sql string, args ...interface{}) (sql.Result, error)
}

// New generates objects that connect to the database.
func New(db *sql.DB, timeout time.Duration) *DB {
	if timeout <= 0 {
		timeout = time.Minute
	}
	return &DB{db, timeout}
}

// Query executes the sql, scans the results into the data, and returns an error.
func (db *DB) Query(data interface{}, sql string, args ...interface{}) error {
	// var people struct {
	// 	Name string
	// 	Age  int
	// }
	// var db *DB
	// if err := db.Query(&people, `select name, age from peoples where id = $1`, 1); err != nil {
	// 	errs.Trace(err)
	// }
	return db.QueryT(db.timeout, data, sql, args...)
}

// Query executes the sql, sets the timeout, scans the result into data, and returns an error.
func (db *DB) QueryT(duration time.Duration, data interface{}, sql string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	return db.query(ctx, data, sql, args)
}

func (db *DB) QueryCtx(ctx context.Context, opName string,
	data interface{}, sql string, args ...interface{},
) error {
	defer tracer.StartSpan(ctx, opName).Finish()
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
	return scan(rows, data)
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
	defer tracer.StartSpan(ctx, opName).Finish()
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
	span := tracer.StartSpan(ctx, opName)
	defer span.Finish()

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
	if err := fn(&Tx{tx, db.timeout}, tracer.Context(context.Background(), span)); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
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
