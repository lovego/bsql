package position

import (
	"fmt"

	runewidth "github.com/mattn/go-runewidth"
)

func Get(content []rune, offset int) string {
	line, column := OffsetToLineAndColumn(content, offset)
	if line <= 0 || column <= 0 {
		return ""
	}

	lineDesc := []rune(fmt.Sprintf("Line %d: ", line))
	columnDesc := []rune(fmt.Sprintf("Char %d: ", column))

	position := lineDesc
	position = append(position, GetLineContent(content, offset, column)...)
	position = append(position, '\n')
	position = append(position, columnDesc...)

	if columnEnd := len(lineDesc) + column; columnEnd > len(columnDesc) {
		position = append(position, makePadding(position[len(columnDesc):columnEnd-1])...)
		position = append(position, '^')
	}
	return string(position)
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

func GetLineContent(content []rune, offset, column int) []rune {
	lineStart := offset - (column - 1)
	if lineStart < 0 {
		lineStart = 0
	}
	lineEnd := GetLineEnd(content, offset)
	return content[lineStart:lineEnd]
}

func GetLineEnd(content []rune, offset int) int {
	if offset < 0 || offset >= len(content) {
		return len(content)
	}

	for i := offset; i < len(content); i++ {
		if content[i] == '\n' {
			return i
		}
	}
	return len(content)
}

func makePadding(content []rune) (result []rune) {
	for _, char := range content {
		if char == '\t' {
			result = append(result, '\t')
		} else {
			for i := 0; i < runewidth.RuneWidth(char); i++ {
				result = append(result, ' ')
			}
		}
	}
	return
}
