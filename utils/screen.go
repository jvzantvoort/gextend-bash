// Package utils provides terminal screen utilities for formatting and displaying text boxes and centered lines.
package utils

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/mitchellh/go-wordwrap"
)

// winsize represents the size of the terminal window.
type winsize struct {
	Row    uint16 // Number of rows
	Col    uint16 // Number of columns
	Xpixel uint16 // Width in pixels
	Ypixel uint16 // Height in pixels
}

// getWidth returns the width of the terminal in columns.
func getWidth() int {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return int(ws.Col)
}

// CenterLine centers a line of text within the given width.
func CenterLine(line string, width int) string {
	line = strings.TrimSpace(line)

	padleft := (width - len(line)) / 2          // number of spaces to add before the line
	retv := strings.Repeat(" ", padleft) + line // add spaces before the line
	strfmt := fmt.Sprintf("%%-%ds", width)
	retv = fmt.Sprintf(strfmt, retv)
	return retv
}

// TextBox prints a formatted text box with a title and message, wrapping the message as needed.
func TextBox(title, format string, args ...interface{}) {
	boxwidth := int(60)
	msg := format

	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	// build the list of lines
	lines := strings.Split(wordwrap.WrapString(msg, uint(boxwidth)), "\n")

	// pad line before and after
	lines = append([]string{""}, lines...)
	lines = append(lines, "")

	header := "+-" + title + strings.Repeat("-", boxwidth-len(title)-3) + "+"
	footer := "+" + strings.Repeat("-", boxwidth-2) + "+"

	fmt.Printf("%s\n", header)
	for _, line := range lines {
		fmt.Printf("| %s |\n", CenterLine(line, boxwidth-4))
	}
	fmt.Printf("%s\n", footer)

}

// ErrorBox prints a text box with the title "Error" and the provided message.
func ErrorBox(format string, args ...interface{}) {
	TextBox("Error", format, args...)
}
