package bsql

import (
	"fmt"
	"strconv"

	"github.com/lib/pq"
	runewidth "github.com/mattn/go-runewidth"
)

func ErrorWithPosition(err error, sqlContent string) error {
	if err == nil {
		return nil
	}
	if pqError, ok := err.(*pq.Error); ok {
		// Position: the field value is a decimal ASCII integer,
		// indicating an error cursor position as an index into the original query string.
		// The first character has index 1, and positions are measured in characters not bytes.
		if offset, err := strconv.Atoi(pqError.Position); err == nil && offset >= 1 {
			position := GetPosition([]rune(sqlContent), int(offset-1))
			if position != "" {
				pqError.Message += "\n" + position
			} else {
				pqError.Message += fmt.Sprintf(" (Position: %s)", pqError.Position)
			}
		}
	}
	return err
}

func GetPosition(content []rune, offset int) string {
	line, column := OffsetToLineAndColumn(content, offset)
	if line <= 0 || column <= 0 {
		return ""
	}
	position := fmt.Sprintf("Line %d: ", line)

	lineStart := offset - (column - 1)
	if lineStart < 0 {
		lineStart = 0
	}
	lineEnd := GetLineEnd(content, offset)
	if lineStart > lineEnd {
		return position
	}

	lineContent := content[lineStart:lineEnd]
	padding := makePadding(append([]rune(position), lineContent[:column-1]...))

	return position + string(lineContent) + "\n" + padding + "^"
}

// line and column begins at 1
func OffsetToLineAndColumn(content []rune, offset int) (int, int) {
	if offset < 0 || offset >= len(content) {
		return 0, 0
	}
	var line, column, lastLineWidth int = 1, 0, 0
	for i := 0; i <= offset; i++ {
		column++
		if content[i] == '\n' {
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

func GetLineEnd(content []rune, offset int) int {
	if offset < 0 || offset >= len(content) {
		return -1
	}

	for i := offset; i < len(content); i++ {
		if content[i] == '\n' {
			return i
		}
	}
	return len(content)
}

func makePadding(content []rune) (result string) {
	for _, char := range content {
		if char == '\t' {
			result += "\t"
		} else {
			for i := 0; i < runewidth.RuneWidth(char); i++ {
				result += " "
			}
		}
	}
	return
}
