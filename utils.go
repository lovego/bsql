package bsql

import "github.com/lib/pq"

func ConflictedUniqueIndex(err error) string {
	if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" { // unique_violation
		return pqError.Constraint
	}
	return ""
}
