package main

import (
	"fmt"
	"strings"
)

func CenterText(instr string, width int) string {
	// pad the string with half of the difference
	restlen := max(width-len(instr), 0)
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
