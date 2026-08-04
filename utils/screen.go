// Package utils provides terminal screen utilities for formatting and displaying text boxes and centered lines.
package utils

import (
	"syscall"
	"unsafe"
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
