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
	INDENT_NO     = 8
	CONSOLE_WIDTH = 80
	B_MARK        = "{#-"
	E_MARK        = "-#}"
)

// PrintCenterText frames instr (left-aligned) in jinja comment markers.
// instr is expected to fit; see utils.WrapLines.
func PrintCenterText(instr string) string {
	if instr == "" {
		return B_MARK + strings.Repeat(" ", CONSOLE_WIDTH-INDENT_NO+2) + E_MARK
	}

	pad := max(CONSOLE_WIDTH-INDENT_NO-runewidth.StringWidth(instr), 0)
	return B_MARK + " " + instr + strings.Repeat(" ", pad) + " " + E_MARK
}

func PrintCenterTextLines(instr string) {
	instr = strings.TrimSpace(instr)
	for _, line := range utils.WrapLines(instr, CONSOLE_WIDTH-INDENT_NO) {
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
