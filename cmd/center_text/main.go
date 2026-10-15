package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
)

const (
	CONSOLE_WIDTH = 80
)

func openStdinOrFile() io.Reader {
	if len(os.Args) > 1 {
		instr := strings.Join(os.Args[1:], " ")
		return strings.NewReader(instr)
	}
	return os.Stdin
}

// PrintCenterText centers instr in a '#' framed line of CONSOLE_WIDTH
// columns. instr is expected to fit; see wrapLines.
func PrintCenterText(instr string) string {
	if instr == "" {
		return "#" + strings.Repeat(" ", CONSOLE_WIDTH-2) + "#"
	}
	// pad the string with half of the difference
	restlen := max(CONSOLE_WIDTH-4-runewidth.StringWidth(instr), 0)
	leftpad, rightpad := restlen/2, restlen-restlen/2

	return "# " + strings.Repeat(" ", leftpad) + instr + strings.Repeat(" ", rightpad) + " #"
}

// wrapLines word-wraps instr to width columns, breaking words that are
// wider than a line on their own.
func wrapLines(instr string, width int) []string {
	var lines []string
	for _, line := range strings.Split(wordwrap.WrapString(instr, uint(width)), "\n") {
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

func PrintCenterTextLines(instr string) {
	instr = strings.TrimSpace(instr)
	for _, line := range wrapLines(instr, CONSOLE_WIDTH-4) {
		fmt.Println(PrintCenterText(line))
	}
}

func main() {

	r := openStdinOrFile()
	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	worstr := string(b)
	PrintCenterTextLines(worstr)
}

// vim: noexpandtab filetype=go
