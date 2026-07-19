package utils

import (
	"strings"
	"testing"

	"github.com/fatih/color"
)

// requireTTY skips the test when stdin isn't a terminal: PrintStatus (via
// stripString) shells out to an ioctl on stdin to get the console width,
// which panics when there is no controlling terminal (e.g. under `go test`
// in CI).
func requireTTY(t *testing.T) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Skip("skipping: no controlling terminal available for getWidth()")
		}
	}()
	getWidth()
}

func TestPrintStatus(t *testing.T) {
	requireTTY(t)
	color.NoColor = true

	out := captureStdout(t, func() {
		PrintStatus(color.FgGreen, "CUSTOM", "doing %s", "work")
	})

	if !strings.Contains(out, "doing work") {
		t.Errorf("expected output to contain %q, got %q", "doing work", out)
	}
	if !strings.Contains(out, "[ CUSTOM ]") {
		t.Errorf("expected output to contain the status label, got %q", out)
	}
}

func TestPrintSuccess(t *testing.T) {
	requireTTY(t)
	color.NoColor = true

	out := captureStdout(t, func() {
		PrintSuccess("task %d", 1)
	})

	if !strings.Contains(out, "task 1") || !strings.Contains(out, "SUCCESS") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestPrintFailed(t *testing.T) {
	requireTTY(t)
	color.NoColor = true

	out := captureStdout(t, func() {
		PrintFailed("task %d", 2)
	})

	if !strings.Contains(out, "task 2") || !strings.Contains(out, "FAILED") {
		t.Errorf("unexpected output: %q", out)
	}
}
