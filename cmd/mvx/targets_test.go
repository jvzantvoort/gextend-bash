package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestDeepCompare(t *testing.T) {
	big := bytes.Repeat([]byte{0}, 3*compareChunkSize+17)
	bigDiffLast := bytes.Clone(big)
	bigDiffLast[len(bigDiffLast)-1] = 'X'

	tests := []struct {
		name string
		src  []byte
		dst  []byte
		want bool
	}{
		{"identical text", []byte("a\nb\n"), []byte("a\nb\n"), true},
		{"different text of equal size", []byte("a\nb\n"), []byte("a\nc\n"), false},
		{"different sizes", []byte("abc"), []byte("abcd"), false},
		{"empty files", nil, nil, true},
		{"identical large binary", big, bytes.Clone(big), true},
		// regression: a line-based scanner gave up on lines over 64KB and
		// reported these as identical, after which mvx deleted the source
		{"large binary differing in the last byte", bigDiffLast, big, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			src := filepath.Join(dir, "src")
			dst := filepath.Join(dir, "dst")
			if err := os.WriteFile(src, tt.src, 0644); err != nil {
				t.Fatalf("os.WriteFile() error = %v", err)
			}
			if err := os.WriteFile(dst, tt.dst, 0644); err != nil {
				t.Fatalf("os.WriteFile() error = %v", err)
			}
			if got := deepCompare(src, dst); got != tt.want {
				t.Errorf("deepCompare() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("unreadable file is never reported identical", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "src")
		if err := os.WriteFile(src, []byte("x"), 0644); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}
		if deepCompare(src, filepath.Join(dir, "missing")) {
			t.Error("deepCompare() = true for a missing destination")
		}
	})
}

func TestMoveFileKeepsDifferentSameSizeFile(t *testing.T) {
	dir := t.TempDir()
	content := bytes.Repeat([]byte{0}, 100000)
	src := filepath.Join(dir, "f.bin")
	srcContent := bytes.Clone(content)
	srcContent[len(srcContent)-1] = 'X'
	if err := os.WriteFile(src, srcContent, 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	dst := filepath.Join(dir, "out")
	if err := os.MkdirAll(dst, 0755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dst, "f.bin"), content, 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	if err := MoveFile(src, dst); err != nil {
		t.Fatalf("MoveFile() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dst, "f.1.bin"))
	if err != nil {
		t.Fatalf("expected the source to be moved to f.1.bin: %v", err)
	}
	if !bytes.Equal(got, srcContent) {
		t.Error("moved file content differs from the original source")
	}
}

func TestDirectoryExistsUnreadableParent(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("permissions are not enforced for root")
	}
	dir := t.TempDir()
	locked := filepath.Join(dir, "locked")
	if err := os.MkdirAll(filepath.Join(locked, "sub"), 0755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}
	if err := os.Chmod(locked, 0); err != nil {
		t.Fatalf("os.Chmod() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0755) })

	if DirectoryExists(filepath.Join(locked, "sub")) {
		t.Error("expected DirectoryExists() to be false when the path cannot be stat'ed")
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.sh")
	if err := os.WriteFile(src, []byte("#!/bin/sh\n"), 0750); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	mtime := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	if err := os.Chtimes(src, mtime, mtime); err != nil {
		t.Fatalf("os.Chtimes() error = %v", err)
	}
	dst := filepath.Join(dir, "dst.sh")

	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile() error = %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("os.Stat() error = %v", err)
	}
	if info.Mode().Perm() != 0750 {
		t.Errorf("mode = %v, want %v", info.Mode().Perm(), os.FileMode(0750))
	}
	if !info.ModTime().Equal(mtime) {
		t.Errorf("mtime = %v, want %v", info.ModTime(), mtime)
	}

	if err := copyFile(src, dst); err == nil {
		t.Error("expected copyFile() to refuse to overwrite an existing file")
	}
}
