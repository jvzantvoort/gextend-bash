package utils

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestWrapLines(t *testing.T) {
	tests := []struct {
		name  string
		instr string
		width int
		want  []string
	}{
		{"short text", "hello world", 20, []string{"hello world"}},
		{"wraps on spaces", "hello world", 7, []string{"hello", "world"}},
		{"breaks words wider than a line", "abcdefghij", 4, []string{"abcd", "efgh", "ij"}},
		{"breaks by display width", "日本語です", 4, []string{"日本", "語で", "す"}},
		{"character wider than the line", "日本", 1, []string{"日", "本", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapLines(tt.instr, tt.width)
			if strings.Join(got, "|") != strings.Join(tt.want, "|") {
				t.Errorf("WrapLines(%q, %d) = %q, want %q", tt.instr, tt.width, got, tt.want)
			}
		})
	}
}

func TestArgsOrStdin(t *testing.T) {
	data, err := io.ReadAll(ArgsOrStdin([]string{"hello", "world"}))
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("content = %q, want %q", data, "hello world")
	}
	if ArgsOrStdin(nil) != os.Stdin {
		t.Error("expected os.Stdin when no args are given")
	}
}
