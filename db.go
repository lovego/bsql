package psql

import (
	"database/sql"
)

type DB struct {
	*sql.DB
}

func (db *DB) RunInTransaction(fn func(*sql.Tx) error) error {
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

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *DB) Query(data interface{}, sql string, args ...interface{}) error {
	rows, err := db.DB.Query(sql, args)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return err
	}
	return Scan(rows, data)
}
