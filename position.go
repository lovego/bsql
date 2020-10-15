package bsql

import (
	"fmt"
	"strconv"

	"github.com/lib/pq"
	"github.com/lovego/bsql/position"
)

func ErrorWithPosition(err error, sqlContent string, fullSql bool) error {
	if err == nil {
		return nil
	}
	if pqError, ok := err.(*pq.Error); ok {
		// Position: the field value is a decimal ASCII integer,
		// indicating an error cursor position as an index into the original query string.
		// The first character has index 1, and positions are measured in characters not bytes.
		if offset, err := strconv.Atoi(pqError.Position); err == nil && offset >= 1 {
			position := position.Get([]rune(sqlContent), int(offset-1))
			if position != "" {
				pqError.Message += "\n" + position
			} else {
				pqError.Message += fmt.Sprintf(" (Position: %s)", pqError.Position)
			}
		}
		if fullSql {
			pqError.Message += "\nSql: " + sqlContent
		}
	}
	return err
}
