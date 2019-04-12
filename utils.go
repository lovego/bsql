package bsql

import (
	"fmt"
	"strconv"

	"github.com/lib/pq"
)

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
	content = content[:offset]
	var line, column, lastLineWidth int64 = 1, 0, 0
	for _, b := range content {
		column++
		if b == '\n' {
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
