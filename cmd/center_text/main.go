package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

func PrintCenterText(instr string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	if instr == "" {
		return "#" + strings.Repeat(" ", CONSOLE_WIDTH-2) + "#"
	}
	// pad the string with half of the difference
	restlen := CONSOLE_WIDTH - 4 - len(instr)
	instr = strings.Repeat(" ", (restlen/2)) + instr

	stringfmt := fmt.Sprintf("# %%-%ds #", CONSOLE_WIDTH-4)

	return fmt.Sprintf(stringfmt, instr)
}

func PrintCenterTextLines(instr string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in PrintCenterTextLines", r)
			fmt.Println(instr)
		}
	}()

	instr = strings.TrimSpace(instr)
	wrapped := wordwrap.WrapString(instr, CONSOLE_WIDTH-4)
	textslice := strings.Split(wrapped, "\n")
	for _, instr := range textslice {
		fmt.Println(PrintCenterText(instr))
	}
}

func main() {

	r := openStdinOrFile()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	worstr := string(b)
	PrintCenterTextLines(worstr)
}

// vim: noexpandtab filetype=go
