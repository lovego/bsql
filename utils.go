package bsql

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

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

func ErrorWithPosition(err error, sqlContent string) error {
	if err == nil {
		return nil
	}
	if pqError, ok := err.(*pq.Error); ok {
		position, err := strconv.ParseInt(pqError.Position, 10, 64)
		var line, column int64
		if err == nil {
			line, column = OffsetToLineAndColumn(sqlContent, position)
		}
		if line > 0 && column > 0 {
			pqError.Message += fmt.Sprintf(" (Line: %d, Column: %d)", line, column)
		} else {
			pqError.Message += fmt.Sprintf(" (Position: %s)", pqError.Position)
		}
	}
	return err
}

// offset should begin at 1.
func OffsetToLineAndColumn(content string, offset int64) (int64, int64) {
	if offset <= 0 || offset > int64(len(content)) {
		return 0, 0
	}
	var line, column, lastLineWidth int64 = 1, 0, 0
	for i := int64(0); i < offset; i++ {
		if len(content) == 0 {
			return 0, 0
		}
		char, size := utf8.DecodeRuneInString(content)
		if char == utf8.RuneError || size == 0 {
			return 0, 0
		}
		content = content[size:]

		column++
		if char == int32('\n') {
			line++
			lastLineWidth = column
			column = 0
		}
	}
	if column == 0 {
		line--
		column = lastLineWidth
	}
	return line, column
}
