package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jvzantvoort/gextend-bash/utils"
	"github.com/jvzantvoort/gextend-bash/version"
	"github.com/mattn/go-runewidth"
)

const (
	CONSOLE_WIDTH = 80
)

// PrintCenterText centers instr in a '#' framed line of CONSOLE_WIDTH
// columns. instr is expected to fit; see utils.WrapLines.
func PrintCenterText(instr string) string {
	if instr == "" {
		return "#" + strings.Repeat(" ", CONSOLE_WIDTH-2) + "#"
	}
	// pad the string with half of the difference
	restlen := max(CONSOLE_WIDTH-4-runewidth.StringWidth(instr), 0)
	leftpad, rightpad := restlen/2, restlen-restlen/2

	return "# " + strings.Repeat(" ", leftpad) + instr + strings.Repeat(" ", rightpad) + " #"
}

func PrintCenterTextLines(instr string) {
	instr = strings.TrimSpace(instr)
	for _, line := range utils.WrapLines(instr, CONSOLE_WIDTH-4) {
		fmt.Println(PrintCenterText(line))
	}
}

func main() {
	version.HandleArgs(os.Args)

	b, err := io.ReadAll(utils.ArgsOrStdin(os.Args[1:]))
	if err != nil {
		log.Fatal(err)
	}
	PrintCenterTextLines(string(b))
}

// vim: noexpandtab filetype=go
