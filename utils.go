package bsql

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func UpsertSql(table string, toInsert, conflictKeys, notToUpdate []string) string {
	toInsert = Fields2Columns(toInsert)
	conflictKeys = Fields2Columns(conflictKeys)
	notToUpdate = Fields2Columns(notToUpdate)

	// conflict keys should be inserted and not be updated.
	for _, key := range conflictKeys {
		if notIn(key, toInsert) {
			toInsert = append(toInsert, key)
		}
		if notIn(key, notToUpdate) {
			notToUpdate = append(notToUpdate, key)
		}
	}

	toUpdate := make([]string, 0, len(toInsert))
	excluded := make([]string, 0, len(toInsert))
	for _, column := range toInsert {
		if notIn(column, notToUpdate) {
			excluded = append(excluded, "excluded."+column)
			if column == "created_at" || column == "created_by" {
				column = strings.Replace(column, "created", "updated", 1)
			}
			toUpdate = append(toUpdate, "         "+column)
		}
	}

	return fmt.Sprintf(`INSERT INTO %s (%s)
VALUES %%s
ON CONFLICT (%s) DO UPDATE SET
(%s) =
(%s)`,
		table, strings.Join(toInsert, ", "), strings.Join(conflictKeys, ", "),
		strings.Join(toUpdate, ", "), strings.Join(excluded, ", "),
	)
}

func ConflictedUniqueIndex(err error) string {
	if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" { // unique_violation
		return pqError.Constraint
	}
	return ""
}
