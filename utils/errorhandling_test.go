package utils

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestWarningOnError(t *testing.T) {
	var buf bytes.Buffer
	orig := log.StandardLogger().Out
	log.SetOutput(&buf)
	defer log.SetOutput(orig)

	t.Run("nil error logs nothing", func(t *testing.T) {
		buf.Reset()
		WarningOnError(nil)
		if buf.Len() != 0 {
			t.Errorf("expected no output for nil error, got %q", buf.String())
		}
	})

	t.Run("error is logged as a warning", func(t *testing.T) {
		buf.Reset()
		WarningOnError(errors.New("boom"))
		if !strings.Contains(buf.String(), "boom") {
			t.Errorf("expected warning output to contain %q, got %q", "boom", buf.String())
		}
	})
}

func TestExitOnError_NilDoesNotExit(t *testing.T) {
	// If this panics or exits the process, the test runner will report the failure.
	ExitOnError(nil)
}

// TestExitOnError_Exits verifies the process-exit path via a re-exec subprocess,
// since os.Exit cannot be intercepted in-process.
func TestExitOnError_Exits(t *testing.T) {
	if os.Getenv("GEXTEND_BASH_EXIT_ON_ERROR_HELPER") == "1" {
		ExitOnError(errors.New("boom"))
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExitOnError_Exits")
	cmd.Env = append(os.Environ(), "GEXTEND_BASH_EXIT_ON_ERROR_HELPER=1")
	err := cmd.Run()

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected process to exit with a non-zero status, got err = %v", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
	}
}
