package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

func TestPrintCenterText(t *testing.T) {
	t.Run("empty input produces a blank padded line", func(t *testing.T) {
		got := PrintCenterText("")
		if len(got) != CONSOLE_WIDTH {
			t.Errorf("length = %d, want %d", len(got), CONSOLE_WIDTH)
		}
		if !strings.HasPrefix(got, "#") || !strings.HasSuffix(got, "#") {
			t.Errorf("PrintCenterText(%q) = %q, want it framed with '#'", "", got)
		}
	})

	t.Run("non-empty input is centered", func(t *testing.T) {
		got := PrintCenterText("hi")
		if !strings.Contains(got, "hi") {
			t.Errorf("PrintCenterText() = %q, want it to contain %q", got, "hi")
		}
	})
}

func TestPrintCenterTextLines(t *testing.T) {
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	PrintCenterTextLines("hello world")
	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() error = %v", err)
	}
	os.Stdout = orig

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "hello world") {
		t.Errorf("output = %q, want it to contain %q", out, "hello world")
	}
}

func TestPrintCenterTextLinesLongWord(t *testing.T) {
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	word := strings.Repeat("y", 100)
	PrintCenterTextLines(word)
	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() error = %v", err)
	}
	os.Stdout = orig

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}

	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
	var text strings.Builder
	for _, line := range lines {
		if len(line) != CONSOLE_WIDTH {
			t.Errorf("line %q has length %d, want %d", line, len(line), CONSOLE_WIDTH)
		}
		text.WriteString(strings.Trim(line, "# "))
	}
	if text.String() != word {
		t.Errorf("output text = %q, want the full word %q", text.String(), word)
	}
}

func TestPrintCenterTextWide(t *testing.T) {
	got := PrintCenterText("日本語")
	if w := runewidth.StringWidth(got); w != CONSOLE_WIDTH {
		t.Errorf("display width = %d, want %d (%q)", w, CONSOLE_WIDTH, got)
	}
}
