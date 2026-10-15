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
	t.Run("empty input is framed with the jinja markers", func(t *testing.T) {
		got := PrintCenterText("")
		if !strings.HasPrefix(got, B_MARK) || !strings.HasSuffix(got, E_MARK) {
			t.Errorf("PrintCenterText(%q) = %q, want it framed with %q/%q", "", got, B_MARK, E_MARK)
		}
	})

	t.Run("non-empty input is wrapped with the markers and kept", func(t *testing.T) {
		got := PrintCenterText("hi")
		if !strings.HasPrefix(got, B_MARK) || !strings.HasSuffix(got, E_MARK) {
			t.Errorf("PrintCenterText(%q) = %q, want it framed with %q/%q", "hi", got, B_MARK, E_MARK)
		}
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
	if !strings.Contains(out, B_MARK) {
		t.Errorf("output = %q, want it to contain the marker %q", out, B_MARK)
	}
}

func TestPrintCenterTextWidth(t *testing.T) {
	want := runewidth.StringWidth(PrintCenterText(""))
	for _, in := range []string{"hi", "日本語", strings.Repeat("x", CONSOLE_WIDTH-INDENT_NO)} {
		if got := runewidth.StringWidth(PrintCenterText(in)); got != want {
			t.Errorf("PrintCenterText(%q) display width = %d, want %d", in, got, want)
		}
	}
}
