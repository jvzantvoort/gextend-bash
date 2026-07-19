package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestOpenStdinOrFile(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	t.Run("uses args when provided", func(t *testing.T) {
		os.Args = []string{"cmd", "hello", "world"}
		r := openStdinOrFile()
		data, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("io.ReadAll() error = %v", err)
		}
		if string(data) != "hello world" {
			t.Errorf("content = %q, want %q", data, "hello world")
		}
	})

	t.Run("falls back to stdin when no args", func(t *testing.T) {
		os.Args = []string{"cmd"}
		r := openStdinOrFile()
		if r != os.Stdin {
			t.Error("expected os.Stdin when no extra args are given")
		}
	})
}

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
	w.Close()
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
