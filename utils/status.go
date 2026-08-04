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
	width := getWidth()

	if width > MAXWITH {
		width = MAXWITH
	}
	width -= WIDTHSUBS

	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	textslice := strings.Split(wordwrap.WrapString(msg, uint(width)), "\n")

	msg = textslice[0]

	return strings.Join([]string{msg, strings.Repeat(".", width-len(msg))}, "")
}

// NormalizeStatus maps common status aliases (e.g. "OK", "WARN") to their
// canonical form (e.g. "SUCCESS", "WARNING"). Unrecognized values are
// returned upper-cased and otherwise unchanged.
func NormalizeStatus(status string) string {

	ustatus := strings.ToUpper(status)
	retv := ustatus

	switch ustatus {
	case "OK", "OKE":
		retv = "SUCCESS"
	case "NOK", "FAIL", "FAILED":
		retv = "FAILURE"
	case "INFO":
		retv = "NOTICE"
	case "WARN":
		retv = "WARNING"
	case "UNDEFINED":
		retv = "UNKNOWN"
	}
	return retv
}

// MakeStatus formats a status message with a status label colored according
// to the normalized status (see NormalizeStatus). The label itself is
// printed as given, not normalized.
func MakeStatus(status, format string, args ...interface{}) string {

	msg := stripString(format, args...)
	state_color := color.New(NoticeColor)
	ustatus := NormalizeStatus(status)

	switch ustatus {
	case "SUCCESS":
		state_color = color.New(SuccessColor)
	case "FAILURE":
		state_color = color.New(FailureColor)
	case "NOTICE":
		state_color = color.New(NoticeColor)
	case "WARNING":
		state_color = color.New(WarningColor)
	case "UNKNOWN":
		state_color = color.New(UnknownColor)
	}

	return fmt.Sprintf("%s [ %-7s ]", msg, state_color.Sprint(status))
}
