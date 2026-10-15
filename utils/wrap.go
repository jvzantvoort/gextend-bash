package utils

import (
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
)

// WrapLines word-wraps text to width display columns, breaking words that
// are wider than a line on their own.
func WrapLines(text string, width int) []string {
	var lines []string
	for _, line := range strings.Split(wordwrap.WrapString(text, uint(width)), "\n") {
		for runewidth.StringWidth(line) > width {
			head := runewidth.Truncate(line, width, "")
			if head == "" {
				// a single character wider than the line; emit it as-is
				_, size := utf8.DecodeRuneInString(line)
				head = line[:size]
			}
			lines = append(lines, head)
			line = line[len(head):]
		}
		lines = append(lines, line)
	}
	return lines
}
