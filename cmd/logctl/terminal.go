package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWinsize() (cols, rows int) {
	ws := &winsize{}
	rc, _, _ := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)
	if rc != 0 || ws.Col == 0 {
		return 0, 0
	}
	return int(ws.Col), int(ws.Row)
}

func consoleWidth() int {
	cols, _ := getWinsize()
	if cols > 0 {
		return cols
	}
	if v := os.Getenv("COLUMNS"); v != "" {
		var w int
		if _, err := fmt.Sscanf(v, "%d", &w); err == nil && w > 0 {
			return w
		}
	}
	return 80
}

func consoleHeight() int {
	_, rows := getWinsize()
	if rows > 0 {
		return rows
	}
	if v := os.Getenv("LINES"); v != "" {
		var h int
		if _, err := fmt.Sscanf(v, "%d", &h); err == nil && h > 0 {
			return h
		}
	}
	return 24
}
