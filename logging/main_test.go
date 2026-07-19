package logging

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jvzantvoort/gextend-bash/config"
	"github.com/spf13/cobra"
)

func TestLogMessageSetLevel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"emerg", "emerg", "EMERG"},
		{"alert", "ALERT", "ALERT"},
		{"crit", "crit", "CRIT"},
		{"err", "err", "ERR"},
		{"warning", "warning", "WARNING"},
		{"notice", "notice", "NOTICE"},
		{"info", "info", "INFO"},
		{"debug", "debug", "DEBUG"},
		{"panic aliases to emerg", "panic", "EMERG"},
		{"error aliases to err", "error", "ERR"},
		{"warn aliases to warning", "warn", "WARNING"},
		{"invalid defaults to notice", "bogus", "NOTICE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LogMessage{}
			l.SetLevel(tt.input)
			if l.Priority != tt.want {
				t.Errorf("SetLevel(%q): Priority = %q, want %q", tt.input, l.Priority, tt.want)
			}
		})
	}

	t.Run("empty level keeps the existing priority", func(t *testing.T) {
		l := &LogMessage{Priority: "INFO"}
		l.SetLevel("")
		if l.Priority != "INFO" {
			t.Errorf("Priority = %q, want %q", l.Priority, "INFO")
		}
	})
}

func TestLogMessageMakeString(t *testing.T) {
	ts := time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC)

	t.Run("without a tag", func(t *testing.T) {
		l := LogMessage{Priority: "INFO", Message: "hello", Time: ts}
		got := string(l.MakeString())
		want := ts.Format(time.RFC3339) + " INFO hello\n"
		if got != want {
			t.Errorf("MakeString() = %q, want %q", got, want)
		}
	})

	t.Run("with a tag", func(t *testing.T) {
		l := LogMessage{Priority: "INFO", Tag: "mytag", Message: "hello", Time: ts}
		got := string(l.MakeString())
		want := ts.Format(time.RFC3339) + " INFO [mytag] hello\n"
		if got != want {
			t.Errorf("MakeString() = %q, want %q", got, want)
		}
	})
}

func TestLogMessageMakeJSONString(t *testing.T) {
	ts := time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC)
	l := LogMessage{Priority: "INFO", Tag: "mytag", Message: "hello", Time: ts}

	data, err := l.MakeJSONString()
	if err != nil {
		t.Fatalf("MakeJSONString() error = %v", err)
	}

	var decoded LogMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if decoded.Priority != "INFO" || decoded.Tag != "mytag" || decoded.Message != "hello" {
		t.Errorf("decoded message = %+v, want Priority=INFO Tag=mytag Message=hello", decoded)
	}
	if !decoded.Time.Equal(ts) {
		t.Errorf("decoded Time = %v, want %v", decoded.Time, ts)
	}
}

func TestGetString(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("tag", "", "")

	t.Run("returns empty when unset", func(t *testing.T) {
		if got := GetString(*cmd, "tag"); got != "" {
			t.Errorf("GetString() = %q, want empty string", got)
		}
	})

	t.Run("returns the flag value when set", func(t *testing.T) {
		if err := cmd.Flags().Set("tag", "mytag"); err != nil {
			t.Fatalf("Flags().Set() error = %v", err)
		}
		if got := GetString(*cmd, "tag"); got != "mytag" {
			t.Errorf("GetString() = %q, want %q", got, "mytag")
		}
	})
}

func TestLogMessagePrint(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(config.ConfigDirEnv, filepath.Join(dir, "cfg"))

	target := filepath.Join(dir, "nested", "out.log")
	ts := time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC)
	l := LogMessage{
		Priority: "INFO",
		Tag:      "mytag",
		Message:  "hello world",
		Time:     ts,
		File:     target,
	}

	if err := l.Print(); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	var decoded LogMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, data = %q", err, data)
	}
	if decoded.Message != "hello world" || decoded.Tag != "mytag" {
		t.Errorf("decoded message = %+v", decoded)
	}
}

func TestNewLogMessage(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(config.ConfigDirEnv, dir)

	l := NewLogMessage("info")
	if l == nil {
		t.Fatal("NewLogMessage() returned nil")
	}
	if l.Priority != "INFO" {
		t.Errorf("Priority = %q, want %q", l.Priority, "INFO")
	}
	if l.Time.IsZero() {
		t.Error("expected Time to be populated")
	}
}

func TestNewLogMessages(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "messages.jsonl")

	older := LogMessage{Priority: "INFO", Message: "first", Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	newer := LogMessage{Priority: "INFO", Message: "second", Time: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)}

	newerJSON, err := newer.MakeJSONString()
	if err != nil {
		t.Fatalf("MakeJSONString() error = %v", err)
	}
	olderJSON, err := older.MakeJSONString()
	if err != nil {
		t.Fatalf("MakeJSONString() error = %v", err)
	}

	// Written out of order on purpose to verify NewLogMessages sorts by time.
	content := string(newerJSON) + "\n" + string(olderJSON) + "\n"
	if err := os.WriteFile(input, []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	msgs := NewLogMessages(input)
	if msgs == nil {
		t.Fatal("NewLogMessages() returned nil")
	}
	if len(msgs.messages) != 2 {
		t.Fatalf("len(messages) = %d, want 2", len(msgs.messages))
	}
	if msgs.messages[0].Message != "first" {
		t.Errorf("messages[0].Message = %q, want %q (oldest first)", msgs.messages[0].Message, "first")
	}
	if msgs.messages[1].Message != "second" {
		t.Errorf("messages[1].Message = %q, want %q", msgs.messages[1].Message, "second")
	}
}
