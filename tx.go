package bsql

import (
	"database/sql"
)

type Tx struct {
	*sql.Tx
}
