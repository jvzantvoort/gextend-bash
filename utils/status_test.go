package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fatih/color"
)

// requireTTY skips the test when stdin isn't a terminal: SprintStatus and
// MakeStatus (via stripString) shell out to an ioctl on stdin to get the
// console width, which panics when there is no controlling terminal (e.g.
// under `go test` in CI).
func requireTTY(t *testing.T) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Skip("skipping: no controlling terminal available for getWidth()")
		}
	}()
	getWidth()
}

func TestNormalizeStatus(t *testing.T) {
	tests := []struct {
		status string
		want   string
	}{
		{"OK", "SUCCESS"},
		{"oke", "SUCCESS"},
		{"nok", "FAILURE"},
		{"FAIL", "FAILURE"},
		{"failed", "FAILURE"},
		{"info", "NOTICE"},
		{"WARN", "WARNING"},
		{"undefined", "UNKNOWN"},
		{"success", "SUCCESS"},
		{"something-else", "SOMETHING-ELSE"},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			if got := NormalizeStatus(tt.status); got != tt.want {
				t.Errorf("NormalizeStatus(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestMakeStatus(t *testing.T) {
	requireTTY(t)
	color.NoColor = true

	t.Run("known status uses the given label as-is", func(t *testing.T) {
		out := MakeStatus("ok", "task %d", 1)
		if !strings.Contains(out, "task 1") {
			t.Errorf("expected output to contain %q, got %q", "task 1", out)
		}
		wantLabel := fmt.Sprintf("[ %-7s ]", "ok")
		if !strings.Contains(out, wantLabel) {
			t.Errorf("expected output to contain %q, got %q", wantLabel, out)
		}
		if strings.HasSuffix(out, "\n") {
			t.Errorf("expected MakeStatus output not to end with a newline, got %q", out)
		}
	})

	t.Run("unknown status still renders the message and label", func(t *testing.T) {
		out := MakeStatus("mystery", "task %d", 2)
		if !strings.Contains(out, "task 2") {
			t.Errorf("expected output to contain %q, got %q", "task 2", out)
		}
		wantLabel := fmt.Sprintf("[ %-7s ]", "mystery")
		if !strings.Contains(out, wantLabel) {
			t.Errorf("expected output to contain %q, got %q", wantLabel, out)
		}
	})
}
