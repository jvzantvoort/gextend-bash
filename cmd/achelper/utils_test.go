package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFileAsList(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "data.txt")
	content := "first\n# a comment\n\nsecond   # inline comment\n   \nthird\n"
	if err := os.WriteFile(target, []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	got, err := readFileAsList(target)
	if err != nil {
		t.Fatalf("readFileAsList() error = %v", err)
	}

	want := []string{"first", "second", "third"}
	if len(got) != len(want) {
		t.Fatalf("readFileAsList() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("readFileAsList()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestReadFileAsListMissingFile(t *testing.T) {
	_, err := readFileAsList(filepath.Join(t.TempDir(), "missing.txt"))
	if err == nil {
		t.Fatal("expected an error for a missing file")
	}
}

func TestExpand(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("cannot determine home dir: %s", err)
	}

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"empty path", "", "", false},
		{"absolute path unchanged", "/etc/passwd", "/etc/passwd", false},
		{"relative path unchanged", "relative/path", "relative/path", false},
		{"tilde alone", "~", home, false},
		{"tilde with subpath", "~/foo/bar", filepath.Join(home, "foo", "bar"), false},
		{"tilde-prefixed user is unsupported", "~someuser/path", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Expand(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Expand(%q) expected an error, got none", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expand(%q) error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("Expand(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
