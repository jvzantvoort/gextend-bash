// Package utils provides terminal screen utilities for formatting and displaying text boxes and centered lines.
package utils

import (
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

// defaultWidth is used when the terminal width cannot be determined, e.g.
// when running from cron, in CI, or with all standard streams redirected.
const defaultWidth = 80

// getWidth returns the width of the terminal in columns. It asks stdout,
// stderr and stdin in turn, then falls back to $COLUMNS and finally to
// defaultWidth.
func getWidth() int {
	for _, fd := range []int{int(os.Stdout.Fd()), int(os.Stderr.Fd()), int(os.Stdin.Fd())} {
		if ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ); err == nil && ws.Col > 0 {
			return int(ws.Col)
		}
	}
	if cols, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && cols > 0 {
		return cols
	}
	return defaultWidth
}
