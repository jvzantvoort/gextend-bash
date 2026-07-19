package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
)

func TestLevelNum(t *testing.T) {
	tests := []struct {
		level string
		want  int
	}{
		{"EMERG", 0},
		{"emerg", 0},
		{"err", 3},
		{"error", 3},
		{"warn", 4},
		{"WARNING", 4},
		{"DEBUG", 7},
		{"not-a-level", 7},
		{"", 7},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			if got := levelNum(tt.level); got != tt.want {
				t.Errorf("levelNum(%q) = %d, want %d", tt.level, got, tt.want)
			}
		})
	}
}

func TestFormatLogMessage(t *testing.T) {
	color.NoColor = true
	ts := time.Date(2026, 7, 19, 10, 30, 0, 0, time.UTC)

	t.Run("includes timestamp, priority, tag and message", func(t *testing.T) {
		msg := logMessage{Tag: "mytag", Priority: "info", Message: "hello", Time: ts}
		got := formatLogMessage(msg, 0)
		want := "2026-07-19 10:30:00 INFO    [mytag] hello"
		if got != want {
			t.Errorf("formatLogMessage() = %q, want %q", got, want)
		}
	})

	t.Run("omits the tag block when empty", func(t *testing.T) {
		msg := logMessage{Priority: "info", Message: "hello", Time: ts}
		got := formatLogMessage(msg, 0)
		if strings.Contains(got, "[]") {
			t.Errorf("formatLogMessage() = %q, want no empty tag brackets", got)
		}
	})

	t.Run("truncates to the given width", func(t *testing.T) {
		msg := logMessage{Priority: "info", Message: strings.Repeat("x", 100), Time: ts}
		got := formatLogMessage(msg, 20)
		if len(got) != 20 {
			t.Errorf("formatLogMessage() length = %d, want %d", len(got), 20)
		}
	})
}

func TestFormatJSONLine(t *testing.T) {
	color.NoColor = true
	ts := time.Date(2026, 7, 19, 10, 30, 0, 0, time.UTC)

	t.Run("empty line returns empty", func(t *testing.T) {
		if got := formatJSONLine("   ", 7, 0); got != "" {
			t.Errorf("formatJSONLine() = %q, want empty", got)
		}
	})

	t.Run("filtered out by level", func(t *testing.T) {
		msg := logMessage{Priority: "DEBUG", Message: "hello", Time: ts}
		data, _ := json.Marshal(msg)
		if got := formatJSONLine(string(data), levelNum("ERR"), 0); got != "" {
			t.Errorf("formatJSONLine() = %q, want empty (filtered by level)", got)
		}
	})

	t.Run("passes the level filter", func(t *testing.T) {
		msg := logMessage{Priority: "ERR", Message: "boom", Time: ts}
		data, _ := json.Marshal(msg)
		got := formatJSONLine(string(data), levelNum("DEBUG"), 0)
		if !strings.Contains(got, "boom") {
			t.Errorf("formatJSONLine() = %q, want it to contain %q", got, "boom")
		}
	})

	t.Run("unparseable JSON is returned truncated but uncolored", func(t *testing.T) {
		raw := "not json " + strings.Repeat("x", 50)
		got := formatJSONLine(raw, 7, 10)
		if got != raw[:10] {
			t.Errorf("formatJSONLine() = %q, want %q", got, raw[:10])
		}
	})
}

func TestSeverityColorReturnsAFunctionForEveryLevel(t *testing.T) {
	for _, level := range []string{"EMERG", "ALERT", "CRIT", "ERR", "ERROR", "WARNING", "WARN", "NOTICE", "INFO", "DEBUG", "unknown"} {
		if fn := severityColor(level); fn == nil {
			t.Errorf("severityColor(%q) returned nil", level)
		}
	}
}

func TestPrintLastN(t *testing.T) {
	color.NoColor = true
	dir := t.TempDir()
	target := filepath.Join(dir, "log.jsonl")

	var lines []string
	ts := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)
	for range 5 {
		msg := logMessage{Priority: "INFO", Message: "line", Time: ts}
		data, _ := json.Marshal(msg)
		lines = append(lines, string(data))
	}
	if err := os.WriteFile(target, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	printErr := printLastN(target, 2, levelNum("DEBUG"), 0)
	w.Close()
	os.Stdout = orig

	if printErr != nil {
		t.Fatalf("printLastN() error = %v", printErr)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}
	got := strings.TrimRight(buf.String(), "\n")
	gotLines := strings.Split(got, "\n")
	if len(gotLines) != 2 {
		t.Fatalf("printLastN() printed %d lines, want 2: %q", len(gotLines), got)
	}
}

func TestPrintLastNMissingFile(t *testing.T) {
	err := printLastN(filepath.Join(t.TempDir(), "missing.jsonl"), 2, levelNum("DEBUG"), 0)
	if err == nil {
		t.Fatal("expected an error for a missing file")
	}
}
