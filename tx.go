package psql

import (
	"database/sql"
)

type Tx struct {
	*sql.Tx
}
