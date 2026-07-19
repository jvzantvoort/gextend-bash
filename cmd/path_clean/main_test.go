package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppendIfMissing(t *testing.T) {
	t.Run("appends a new element", func(t *testing.T) {
		got := AppendIfMissing([]string{"a", "b"}, "c")
		want := []string{"a", "b", "c"}
		if len(got) != len(want) {
			t.Fatalf("AppendIfMissing() = %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("does not duplicate an existing element", func(t *testing.T) {
		got := AppendIfMissing([]string{"a", "b"}, "b")
		want := []string{"a", "b"}
		if len(got) != len(want) {
			t.Fatalf("AppendIfMissing() = %v, want %v", got, want)
		}
	})
}

func TestFilterExists(t *testing.T) {
	dir := t.TempDir()
	existingDir := filepath.Join(dir, "existing")
	if err := os.Mkdir(existingDir, 0755); err != nil {
		t.Fatalf("os.Mkdir() error = %v", err)
	}
	existingFile := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(existingFile, []byte("x"), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	missingDir := filepath.Join(dir, "missing")

	got := FilterExists([]string{"", existingDir, existingFile, missingDir, existingDir})

	want := []string{existingDir}
	if len(got) != len(want) {
		t.Fatalf("FilterExists() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
