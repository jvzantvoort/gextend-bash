package utils

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCenterLine(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		width int
		want  string
	}{
		{"short word centered", "hi", 10, "    hi    "},
		{"trims surrounding whitespace", "  hi  ", 10, "    hi    "},
		{"empty line", "", 6, "      "},
		{"line as wide as width", "abcdef", 6, "abcdef"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CenterLine(tt.line, tt.width)
			if got != tt.want {
				t.Errorf("CenterLine(%q, %d) = %q, want %q", tt.line, tt.width, got, tt.want)
			}
			if len(got) != tt.width {
				t.Errorf("CenterLine(%q, %d) length = %d, want %d", tt.line, tt.width, len(got), tt.width)
			}
		})
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = orig }()

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

func TestTextBox(t *testing.T) {
	out := captureStdout(t, func() {
		TextBox("Title", "hello %s", "world")
	})

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least a header and footer line, got %d lines: %q", len(lines), out)
	}
	if !strings.HasPrefix(lines[0], "+-Title") {
		t.Errorf("expected header to start with %q, got %q", "+-Title", lines[0])
	}
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected body to contain %q, got %q", "hello world", out)
	}
	last := lines[len(lines)-1]
	if !strings.HasPrefix(last, "+-") || !strings.HasSuffix(last, "-+") {
		t.Errorf("expected footer to look like a boxed line, got %q", last)
	}
}

func TestErrorBox(t *testing.T) {
	out := captureStdout(t, func() {
		ErrorBox("something %s", "broke")
	})

	if !strings.Contains(out, "Error") {
		t.Errorf("expected output to contain the title %q, got %q", "Error", out)
	}
	if !strings.Contains(out, "something broke") {
		t.Errorf("expected output to contain %q, got %q", "something broke", out)
	}
}
