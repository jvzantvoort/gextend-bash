package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/jvzantvoort/gextend-bash/utils"
	"github.com/spf13/cobra"
)

// requireTTY skips the test when stdin isn't a terminal: PrettyPrint.Print
// (via utils.MakeStatus) shells out to an ioctl on stdin to get the console
// width, which panics when there is no controlling terminal (e.g. under
// `go test` in CI).
func requireTTY(t *testing.T) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Skip("skipping: no controlling terminal available for getWidth()")
		}
	}()
	utils.MakeStatus("probe", "probe")
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stderr = w
	defer func() { os.Stderr = orig }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() error = %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}
	return buf.String()
}

func TestNewPrettyPrint(t *testing.T) {
	got := NewPrettyPrint("success")
	if got.Status != "SUCCESS" {
		t.Errorf("Status = %q, want %q", got.Status, "SUCCESS")
	}
	if got.Message != "" {
		t.Errorf("Message = %q, want empty", got.Message)
	}
}

func TestPrettyPrintPrint(t *testing.T) {
	requireTTY(t)
	color.NoColor = true

	pp := NewPrettyPrint("warning")
	pp.Message = "disk almost full"

	out := captureStderr(t, func() {
		if err := pp.Print(); err != nil {
			t.Fatalf("Print() error = %v", err)
		}
	})

	if !strings.Contains(out, "disk almost full") {
		t.Errorf("expected output to contain %q, got %q", "disk almost full", out)
	}
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected output to contain the status label, got %q", out)
	}
}

func TestHandlePrintCmd(t *testing.T) {
	requireTTY(t)
	color.NoColor = true

	cmd := &cobra.Command{Use: "success"}
	out := captureStderr(t, func() {
		handlePrintCmd(cmd, []string{"hello", "world"})
	})

	if !strings.Contains(out, "hello world") {
		t.Errorf("expected output to contain %q, got %q", "hello world", out)
	}
	if !strings.Contains(out, "SUCCESS") {
		t.Errorf("expected output to contain the status label, got %q", out)
	}
}
