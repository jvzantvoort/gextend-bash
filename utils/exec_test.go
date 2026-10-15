package utils

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestExitCode(t *testing.T) {
	notExecutable := filepath.Join(t.TempDir(), "script")
	if err := os.WriteFile(notExecutable, []byte("#!/bin/sh\n"), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	tests := []struct {
		name string
		cmd  *exec.Cmd
		want int
	}{
		{"success", exec.Command("sh", "-c", "exit 0"), 0},
		{"failure status", exec.Command("sh", "-c", "exit 3"), 3},
		{"killed by signal", exec.Command("sh", "-c", "kill -TERM $$"), 128 + 15},
		{"not found", exec.Command("gextend-no-such-command"), 127},
		{"not executable", exec.Command(notExecutable), 126},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExitCode(tt.cmd.Run()); got != tt.want {
				t.Errorf("ExitCode() = %d, want %d", got, tt.want)
			}
		})
	}

	if got := ExitCode(errors.New("other")); got != 126 {
		t.Errorf("ExitCode(other) = %d, want 126", got)
	}
}
