package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
)

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

func TestGetWidthWithoutTerminal(t *testing.T) {
	t.Setenv("COLUMNS", "")
	if got := getWidth(); got <= 0 {
		t.Errorf("getWidth() = %d, want a positive width", got)
	}
}

func TestStripString(t *testing.T) {
	width := max(min(getWidth(), MAXWITH)-WIDTHSUBS, minMsgWidth)

	tests := []struct {
		name   string
		msg    string
		prefix string
	}{
		{"short message is padded", "hello", "hello."},
		{"word longer than the line is truncated", strings.Repeat("x", 500), strings.Repeat("x", width)},
		{"wide characters are measured by display width", "日本語", "日本語."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripString("%s", tt.msg)
			if gotWidth := runewidth.StringWidth(got); gotWidth != width {
				t.Errorf("stripString() display width = %d, want %d (%q)", gotWidth, width, got)
			}
			if !strings.HasPrefix(got, tt.prefix) {
				t.Errorf("stripString() = %q, want prefix %q", got, tt.prefix)
			}
		})
	}
}
