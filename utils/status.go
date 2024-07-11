package utils

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/mitchellh/go-wordwrap"
)

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

func PrintStatus(colorattr color.Attribute, status, format string, args ...interface{}) {

	msg := stripString(format, args...)
	state_color := color.New(colorattr)

	fmt.Printf("%s [ %s ]\n", msg, state_color.Sprint(status))
}

func PrintSuccess(format string, args ...interface{}) {
	PrintStatus(SuccessColor, "SUCCESS", format, args...)
}

func PrintFailed(format string, args ...interface{}) {
	PrintStatus(FailureColor, "FAILED", format, args...)
}
