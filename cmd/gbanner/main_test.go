package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestOpenStdinOrFile(t *testing.T) {
	t.Run("uses args when provided", func(t *testing.T) {
		r := openStdinOrFile([]string{"hello", "world"})
		data, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("io.ReadAll() error = %v", err)
		}
		if string(data) != "hello world" {
			t.Errorf("openStdinOrFile() content = %q, want %q", data, "hello world")
		}
	})

	t.Run("falls back to stdin when no args", func(t *testing.T) {
		r := openStdinOrFile(nil)
		if r != os.Stdin {
			t.Error("expected openStdinOrFile() to return os.Stdin when no args given")
		}
	})
}

func TestColorize(t *testing.T) {
	color.NoColor = false

	t.Run("empty color name leaves line unmodified", func(t *testing.T) {
		if got := colorize("line", ""); got != "line" {
			t.Errorf("colorize() = %q, want %q", got, "line")
		}
	})

	t.Run("none leaves line unmodified", func(t *testing.T) {
		if got := colorize("line", "none"); got != "line" {
			t.Errorf("colorize() = %q, want %q", got, "line")
		}
	})

	t.Run("unknown color leaves line unmodified", func(t *testing.T) {
		if got := colorize("line", "not-a-color"); got != "line" {
			t.Errorf("colorize() = %q, want %q", got, "line")
		}
	})

	t.Run("known color wraps the line", func(t *testing.T) {
		got := colorize("line", "red")
		if got == "line" {
			t.Error("expected the line to be colorized")
		}
		if !strings.Contains(got, "line") {
			t.Errorf("colorize() = %q, want it to still contain the original text", got)
		}
	})

	t.Run("rainbow colorizes every character", func(t *testing.T) {
		got := colorize("ab", Rainbow)
		if got == "ab" {
			t.Error("expected rainbow output to differ from the plain input")
		}
	})
}

func TestBuildLine(t *testing.T) {
	t.Run("centers short text", func(t *testing.T) {
		got := buildLine("hi", 10, false)
		want := Vertical + " " + "  hi  " + " " + Vertical
		if got != want {
			t.Errorf("buildLine() = %q, want %q", got, want)
		}
	})

	t.Run("left-aligns when requested", func(t *testing.T) {
		got := buildLine("hi", 10, true)
		want := Vertical + " " + "hi    " + " " + Vertical
		if got != want {
			t.Errorf("buildLine() = %q, want %q", got, want)
		}
	})

	t.Run("truncates text wider than the box", func(t *testing.T) {
		got := buildLine("abcdefgh", 8, false)
		// width = boxWidth-4 = 4, text truncated to "abcd"
		want := Vertical + " " + "abcd" + " " + Vertical
		if got != want {
			t.Errorf("buildLine() = %q, want %q", got, want)
		}
	})
}

func TestPrintBanner(t *testing.T) {
	color.NoColor = true

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	PrintBanner("hello", "", 0, 0, false)
	w.Close()
	os.Stdout = orig

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}
	out := buf.String()

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected at least top, content and bottom lines, got %d: %q", len(lines), out)
	}
	if !strings.HasPrefix(lines[0], TopLeft) {
		t.Errorf("expected first line to start with %q, got %q", TopLeft, lines[0])
	}
	last := lines[len(lines)-1]
	if !strings.HasPrefix(last, BottomLeft) {
		t.Errorf("expected last line to start with %q, got %q", BottomLeft, last)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected output to contain %q, got %q", "hello", out)
	}
}
