package bsql

import (
	"context"
	"database/sql"
	"time"
)

type Tx struct {
	*sql.Tx
	Timeout time.Duration
}

func (tx *Tx) Query(data interface{}, sql string, args ...interface{}) error {
	if tx.Timeout > 0 {
		return tx.QueryT(tx.Timeout, data, sql, args...)
	} else {
		return tx.QueryT(time.Minute, data, sql, args...)
	}
}

func (tx *Tx) QueryT(duration time.Duration, data interface{}, sql string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	rows, err := tx.Tx.QueryContext(ctx, sql, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return err
	}
	return Scan(rows, data)
}

func (tx *Tx) Exec(sql string, args ...interface{}) (sql.Result, error) {
	if tx.Timeout > 0 {
		return tx.ExecT(tx.Timeout, sql, args)
	} else {
		return tx.ExecT(time.Minute, sql, args)
	}
}

func (tx *Tx) ExecT(duration time.Duration, sql string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	return tx.Tx.ExecContext(ctx, sql, args)
}
