package bsql

import (
	"context"
	"database/sql"
	"time"
)

type DB struct {
	*sql.DB
	Timeout time.Duration
}

func (db *DB) Query(data interface{}, sql string, args ...interface{}) error {
	if db.Timeout > 0 {
		return db.QueryT(db.Timeout, data, sql, args...)
	} else {
		return db.QueryT(time.Minute, data, sql, args...)
	}
}

func (db *DB) QueryT(duration time.Duration, data interface{}, sql string, args ...interface{}) error {
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
