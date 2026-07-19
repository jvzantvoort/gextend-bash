package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormatDestFile(t *testing.T) {
	tests := []struct {
		name       string
		destdir    string
		sourcefile string
		num        int
		want       string
	}{
		{"no suffix uses base name", "/dest", "/some/path/file.txt", 0, "/dest/file.txt"},
		{"suffix with extension", "/dest", "/some/path/file.txt", 2, "/dest/file.2.txt"},
		{"suffix without extension", "/dest", "/some/path/file", 3, "/dest/file.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDestFile(tt.destdir, tt.sourcefile, tt.num)
			if got != tt.want {
				t.Errorf("formatDestFile(%q, %q, %d) = %q, want %q", tt.destdir, tt.sourcefile, tt.num, got, tt.want)
			}
		})
	}
}

func TestGetFileSize(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		dir := t.TempDir()
		target := filepath.Join(dir, "f.txt")
		if err := os.WriteFile(target, []byte("hello"), 0644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}
		if got := GetFileSize(target); got != 5 {
			t.Errorf("GetFileSize() = %d, want 5", got)
		}
	})

	t.Run("missing file returns zero", func(t *testing.T) {
		if got := GetFileSize(filepath.Join(t.TempDir(), "missing")); got != 0 {
			t.Errorf("GetFileSize() = %d, want 0", got)
		}
	})
}

func TestDirectoryExists(t *testing.T) {
	dir := t.TempDir()

	if !DirectoryExists(dir) {
		t.Error("expected DirectoryExists() to be true for an existing directory")
	}

	file := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	if DirectoryExists(file) {
		t.Error("expected DirectoryExists() to be false for a regular file")
	}

	if DirectoryExists(filepath.Join(dir, "missing")) {
		t.Error("expected DirectoryExists() to be false for a missing path")
	}
}

func TestTargetExistsAndIsNotADirectory(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	if !TargetExistsAndIsNotADirectory(file) {
		t.Error("expected true for an existing regular file")
	}
	if TargetExistsAndIsNotADirectory(dir) {
		t.Error("expected false for a directory")
	}
	if TargetExistsAndIsNotADirectory(filepath.Join(dir, "missing")) {
		t.Error("expected false for a missing path")
	}
}

func TestGetNextTarget(t *testing.T) {
	t.Run("returns the base name when free", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "src.txt")
		if err := os.WriteFile(src, []byte("source content"), 0644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}
		destdir := filepath.Join(dir, "dest")
		if err := os.MkdirAll(destdir, 0755); err != nil {
			t.Fatalf("os.MkdirAll() error = %v", err)
		}

		got, err := GetNextTarget(destdir, src)
		if err != nil {
			t.Fatalf("GetNextTarget() error = %v", err)
		}
		want := filepath.Join(destdir, "src.txt")
		if got != want {
			t.Errorf("GetNextTarget() = %q, want %q", got, want)
		}
	})

	t.Run("increments the suffix when the target exists but differs", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "src.txt")
		if err := os.WriteFile(src, []byte("source content"), 0644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}
		destdir := filepath.Join(dir, "dest")
		if err := os.MkdirAll(destdir, 0755); err != nil {
			t.Fatalf("os.MkdirAll() error = %v", err)
		}
		existing := filepath.Join(destdir, "src.txt")
		if err := os.WriteFile(existing, []byte("different content"), 0644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}

		got, err := GetNextTarget(destdir, src)
		if err != nil {
			t.Fatalf("GetNextTarget() error = %v", err)
		}
		want := filepath.Join(destdir, "src.1.txt")
		if got != want {
			t.Errorf("GetNextTarget() = %q, want %q", got, want)
		}
	})

	t.Run("returns ErrSameFile when the base target is identical", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "src.txt")
		if err := os.WriteFile(src, []byte("source content"), 0644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}
		destdir := filepath.Join(dir, "dest")
		if err := os.MkdirAll(destdir, 0755); err != nil {
			t.Fatalf("os.MkdirAll() error = %v", err)
		}
		identical := filepath.Join(destdir, "src.txt")
		if err := os.WriteFile(identical, []byte("source content"), 0644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}

		_, err := GetNextTarget(destdir, src)
		if err != ErrSameFile {
			t.Errorf("GetNextTarget() error = %v, want %v", err, ErrSameFile)
		}
	})
}

func TestMoveFile(t *testing.T) {
	t.Run("moves the file into the destination directory", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "a.txt")
		if err := os.WriteFile(src, []byte("payload"), 0644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}
		dst := filepath.Join(dir, "out")

		if err := MoveFile(src, dst); err != nil {
			t.Fatalf("MoveFile() error = %v", err)
		}

		if _, err := os.Stat(src); !os.IsNotExist(err) {
			t.Errorf("expected source to be removed, stat err = %v", err)
		}

		moved := filepath.Join(dst, "a.txt")
		data, err := os.ReadFile(moved)
		if err != nil {
			t.Fatalf("os.ReadFile() error = %v", err)
		}
		if string(data) != "payload" {
			t.Errorf("moved file content = %q, want %q", data, "payload")
		}
	})

	t.Run("returns ErrSrcNoExist for a missing source", func(t *testing.T) {
		dir := t.TempDir()
		err := MoveFile(filepath.Join(dir, "missing"), filepath.Join(dir, "out"))
		if err != ErrSrcNoExist {
			t.Errorf("MoveFile() error = %v, want %v", err, ErrSrcNoExist)
		}
	})

	t.Run("removes source when identical file already exists at destination", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "a.txt")
		if err := os.WriteFile(src, []byte("same"), 0644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}
		dst := filepath.Join(dir, "out")
		if err := os.MkdirAll(dst, 0755); err != nil {
			t.Fatalf("os.MkdirAll() error = %v", err)
		}
		if err := os.WriteFile(filepath.Join(dst, "a.txt"), []byte("same"), 0644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}

		if err := MoveFile(src, dst); err != nil {
			t.Fatalf("MoveFile() error = %v", err)
		}
		if _, err := os.Stat(src); !os.IsNotExist(err) {
			t.Errorf("expected source to be removed, stat err = %v", err)
		}
	})
}
