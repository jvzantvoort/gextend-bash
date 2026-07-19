package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestShortHostname(t *testing.T) {
	fqdn, err := os.Hostname()
	if err != nil {
		t.Skipf("cannot determine hostname: %s", err)
	}
	want := strings.ToLower(strings.SplitN(fqdn, ".", 2)[0])
	got := ShortHostname()
	if got != want {
		t.Errorf("ShortHostname() = %q, want %q", got, want)
	}
}

func TestGetHomeDir(t *testing.T) {
	got := GetHomeDir()
	if got == "" {
		t.Error("GetHomeDir() returned an empty string")
	}
}

func TestMkdirP(t *testing.T) {
	base := t.TempDir()

	t.Run("creates nested directory", func(t *testing.T) {
		target := filepath.Join(base, "a", "b", "c")
		if err := MkdirP(target, 0755); err != nil {
			t.Fatalf("MkdirP() error = %v", err)
		}
		info, err := os.Stat(target)
		if err != nil {
			t.Fatalf("expected directory to exist: %v", err)
		}
		if !info.IsDir() {
			t.Fatal("expected target to be a directory")
		}
	})

	t.Run("existing directory is a no-op", func(t *testing.T) {
		target := filepath.Join(base, "already-exists")
		if err := os.Mkdir(target, 0755); err != nil {
			t.Fatalf("setup Mkdir() error = %v", err)
		}
		if err := MkdirP(target, 0755); err != nil {
			t.Fatalf("MkdirP() on existing dir error = %v", err)
		}
	})

	t.Run("errors when target is a file", func(t *testing.T) {
		target := filepath.Join(base, "afile")
		if err := os.WriteFile(target, []byte("x"), 0644); err != nil {
			t.Fatalf("setup WriteFile() error = %v", err)
		}
		if err := MkdirP(target, 0755); err == nil {
			t.Fatal("expected an error when target exists and is not a directory")
		}
	})
}

func TestFileExists(t *testing.T) {
	base := t.TempDir()

	t.Run("existing file", func(t *testing.T) {
		target := filepath.Join(base, "file.txt")
		if err := os.WriteFile(target, []byte("hello"), 0644); err != nil {
			t.Fatalf("setup WriteFile() error = %v", err)
		}
		ok, info := FileExists(target)
		if !ok {
			t.Fatal("expected FileExists() to return true")
		}
		if info == nil || info.Size() != 5 {
			t.Fatalf("unexpected file info: %+v", info)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		ok, _ := FileExists(filepath.Join(base, "missing.txt"))
		if ok {
			t.Fatal("expected FileExists() to return false for a missing file")
		}
	})

	t.Run("directory is not a file", func(t *testing.T) {
		dir := filepath.Join(base, "adir")
		if err := os.Mkdir(dir, 0755); err != nil {
			t.Fatalf("setup Mkdir() error = %v", err)
		}
		ok, _ := FileExists(dir)
		if ok {
			t.Fatal("expected FileExists() to return false for a directory")
		}
	})
}

func TestFileIsExecutable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits are not meaningful on windows")
	}

	base := t.TempDir()

	t.Run("missing file", func(t *testing.T) {
		if FileIsExecutable(filepath.Join(base, "missing")) {
			t.Fatal("expected false for a missing file")
		}
	})

	t.Run("non-executable file", func(t *testing.T) {
		target := filepath.Join(base, "plain.txt")
		if err := os.WriteFile(target, []byte("x"), 0644); err != nil {
			t.Fatalf("setup WriteFile() error = %v", err)
		}
		if FileIsExecutable(target) {
			t.Fatal("expected false for a non-executable file")
		}
	})

	t.Run("executable file", func(t *testing.T) {
		target := filepath.Join(base, "run.sh")
		if err := os.WriteFile(target, []byte("x"), 0755); err != nil {
			t.Fatalf("setup WriteFile() error = %v", err)
		}
		if !FileIsExecutable(target) {
			t.Fatal("expected true for an executable file")
		}
	})
}
