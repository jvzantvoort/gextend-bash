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
	INDENT_NO = 8
	CONSOLE_WIDTH = 80
	B_MARK = "{#-"
	E_MARK = "-#}"

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
		return B_MARK + strings.Repeat(" ", CONSOLE_WIDTH-INDENT_NO+2) + E_MARK
	}

	stringfmt := fmt.Sprintf("%s %%-%ds %s", B_MARK, CONSOLE_WIDTH-INDENT_NO, E_MARK)

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
	wrapped := wordwrap.WrapString(instr, CONSOLE_WIDTH-INDENT_NO)
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
