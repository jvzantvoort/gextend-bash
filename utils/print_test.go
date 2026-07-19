package utils

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestPrintError(t *testing.T) {
	var buf bytes.Buffer
	orig := log.StandardLogger().Out
	log.SetOutput(&buf)
	defer log.SetOutput(orig)

	t.Run("nil error is a no-op", func(t *testing.T) {
		buf.Reset()
		if err := PrintError("failed: %s", nil); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if buf.Len() != 0 {
			t.Errorf("expected no log output, got %q", buf.String())
		}
	})

	t.Run("non-nil error is logged and returned", func(t *testing.T) {
		buf.Reset()
		wantErr := errors.New("boom")
		gotErr := PrintError("failed: %s", wantErr)
		if gotErr != wantErr {
			t.Errorf("expected returned error to be the same instance, got %v", gotErr)
		}
		if !strings.Contains(buf.String(), "boom") {
			t.Errorf("expected log output to contain %q, got %q", "boom", buf.String())
		}
	})
}

func TestPrintFatal(t *testing.T) {
	var buf bytes.Buffer
	orig := log.StandardLogger().Out
	log.SetOutput(&buf)
	defer log.SetOutput(orig)

	t.Run("nil error is a no-op and does not exit", func(t *testing.T) {
		buf.Reset()
		if err := PrintFatal("failed: %s", nil); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if buf.Len() != 0 {
			t.Errorf("expected no log output, got %q", buf.String())
		}
	})

	t.Run("non-nil error is logged without exiting the process", func(t *testing.T) {
		exited := false
		origExit := log.StandardLogger().ExitFunc
		log.StandardLogger().ExitFunc = func(int) { exited = true }
		defer func() { log.StandardLogger().ExitFunc = origExit }()

		buf.Reset()
		wantErr := errors.New("boom")
		gotErr := PrintFatal("failed: %s", wantErr)
		if gotErr != wantErr {
			t.Errorf("expected returned error to be the same instance, got %v", gotErr)
		}
		if !strings.Contains(buf.String(), "boom") {
			t.Errorf("expected log output to contain %q, got %q", "boom", buf.String())
		}
		if !exited {
			t.Error("expected the logger's ExitFunc to be invoked")
		}
	})
}

func TestPanicOnError(t *testing.T) {
	var buf bytes.Buffer
	orig := log.StandardLogger().Out
	log.SetOutput(&buf)
	defer log.SetOutput(orig)

	t.Run("nil error does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("did not expect a panic, got %v", r)
			}
		}()
		PanicOnError("failed: %s", nil)
	})

	t.Run("non-nil error panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected a panic for a non-nil error")
			}
		}()
		PanicOnError("failed: %s", errors.New("boom"))
	})
}
