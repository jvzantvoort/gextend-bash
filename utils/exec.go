package utils

import (
	"errors"
	"io/fs"
	"os/exec"
	"syscall"
)

// ExitCode maps the error of running a command to the exit status a shell
// would report: the command's own status, 128+signal when it was killed,
// 127 when it was not found and 126 when it could not be executed.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
			return 128 + int(status.Signal())
		}
		return exitErr.ExitCode()
	}

	if errors.Is(err, exec.ErrNotFound) || errors.Is(err, fs.ErrNotExist) {
		return 127
	}
	return 126
}
