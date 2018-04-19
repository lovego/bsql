package bsql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

type DB struct {
	*sql.DB
	Timeout time.Duration
}
type DbOrTx interface {
	Query(data interface{}, sql string, args ...interface{}) error
	QueryT(duration time.Duration, data interface{}, sql string, args ...interface{}) error
	Exec(sql string, args ...interface{}) (sql.Result, error)
	ExecT(duration time.Duration, sql string, args ...interface{}) (sql.Result, error)
}

func (db *DB) Query(data interface{}, sql string, args ...interface{}) error {
	if db.Timeout > 0 {
		return db.QueryT(db.Timeout, data, sql, args...)
	} else {
		return db.QueryT(time.Minute, data, sql, args...)
	}
}

func (db *DB) QueryT(duration time.Duration, data interface{}, sql string, args ...interface{}) error {
	debugBsql(sql, args...)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	rows, err := db.DB.QueryContext(ctx, sql, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return err
	}
	return scan(rows, data)
}

func (db *DB) Exec(sql string, args ...interface{}) (sql.Result, error) {
	if db.Timeout > 0 {
		return db.ExecT(db.Timeout, sql, args...)
	} else {
		return db.ExecT(time.Minute, sql, args...)
	}
}

func (db *DB) ExecT(duration time.Duration, sql string, args ...interface{}) (sql.Result, error) {
	debugBsql(sql, args...)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	return db.DB.ExecContext(ctx, sql, args...)
}

func (db *DB) RunInTransaction(fn func(*Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()

	if err := fn(&Tx{tx, db.Timeout}); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

var DebugBsql = os.Getenv(`DebugBsql`) != ``

func debugBsql(sql string, args ...interface{}) {
	if DebugBsql {
		color.Green(sql)
		argsString := ``
		for _, arg := range args {
			argsString += fmt.Sprintf("%#v", arg)
		}
		color.Blue(argsString)
	}
}
