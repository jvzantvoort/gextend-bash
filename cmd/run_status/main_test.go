package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestRun(t *testing.T) {
	color.NoColor = true

	tests := []struct {
		name       string
		command    []string
		message    string
		quiet      bool
		wantCode   int
		wantStdout []string
		wantStderr string
	}{
		{"success shows output and status", []string{"sh", "-c", "echo hello"}, "greet", false, 0,
			[]string{"hello\n", "greet.", "[ SUCCESS ]"}, ""},
		{"default message is the command line", []string{"true"}, "", false, 0,
			[]string{"true.", "[ SUCCESS ]"}, ""},
		{"failure keeps the exit code", []string{"sh", "-c", "exit 3"}, "fail", false, 3,
			[]string{"fail.", "[ FAILURE ]"}, ""},
		{"quiet hides output on success", []string{"sh", "-c", "echo noisy"}, "q", true, 0,
			[]string{"[ SUCCESS ]"}, ""},
		{"quiet shows output on failure", []string{"sh", "-c", "echo out; echo err >&2; exit 1"}, "q", true, 1,
			[]string{"[ FAILURE ]"}, "out\nerr\n"},
		{"missing command", []string{"gextend-no-such-command"}, "x", false, 127,
			[]string{"x (not found)", "[ FAILURE ]"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(execRunner, tt.command, tt.message, tt.quiet, &stdout, &stderr)
			if code != tt.wantCode {
				t.Errorf("run() = %d, want %d", code, tt.wantCode)
			}
			for _, want := range tt.wantStdout {
				if !strings.Contains(stdout.String(), want) {
					t.Errorf("stdout = %q, want it to contain %q", stdout.String(), want)
				}
			}
			if tt.quiet && strings.Contains(stdout.String(), "noisy") {
				t.Errorf("quiet run leaked output: %q", stdout.String())
			}
			if stderr.String() != tt.wantStderr {
				t.Errorf("stderr = %q, want %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}
