package main

import (
	"fmt"
	"strings"
)

func CenterText(instr string, width int) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	// pad the string with half of the difference
	restlen := width - len(instr)
	instr = strings.Repeat(" ", (restlen/2)) + instr

	stringfmt := fmt.Sprintf("%%-%ds", width)

	return fmt.Sprintf(stringfmt, instr)
}

func lastN(input string, length int) string {
	if len(input) <= length {
		return input
	}
	arglen := len(input) - length + 4

	return "... " + input[arglen:]
}
