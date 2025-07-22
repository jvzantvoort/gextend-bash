// Package utils provides utility functions for formatting and printing status messages.
package utils

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/mitchellh/go-wordwrap"
)

// stripString formats a string with optional arguments, wraps it to the terminal width minus WIDTHSUBS,
// and pads it with dots to fill the line.
func stripString(format string, args ...interface{}) string {

	msg := format
	width := getWidth() - WIDTHSUBS

	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	textslice := strings.Split(wordwrap.WrapString(msg, uint(width)), "\n")

	msg = textslice[0]

	return strings.Join([]string{msg, strings.Repeat(".", width-len(msg))}, "")
}

// PrintStatus prints a formatted status message with a colored status label.
func PrintStatus(colorattr color.Attribute, status, format string, args ...interface{}) {

	msg := stripString(format, args...)
	state_color := color.New(colorattr)

	fmt.Printf("%s [ %s ]\n", msg, state_color.Sprint(status))
}

// PrintSuccess prints a success status message in green.
func PrintSuccess(format string, args ...interface{}) {
	PrintStatus(SuccessColor, "SUCCESS", format, args...)
}

// PrintFailed prints a failed status message in red.
func PrintFailed(format string, args ...interface{}) {
	PrintStatus(FailureColor, "FAILED", format, args...)
}
