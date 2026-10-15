package main

import "testing"

func TestMunge(t *testing.T) {
	tests := []struct {
		name string
		path string
		dirs []string
		mode Mode
		want string
	}{
		{"prepend new", "/usr/bin:/bin", []string{"/opt/bin"}, Prepend, "/opt/bin:/usr/bin:/bin"},
		{"prepend keeps given order", "/bin", []string{"/a", "/b"}, Prepend, "/a:/b:/bin"},
		{"prepend existing is a no-op", "/usr/bin:/opt/bin", []string{"/opt/bin"}, Prepend, "/usr/bin:/opt/bin"},
		{"prepend compares cleaned paths", "/usr/bin:/opt/bin/", []string{"/opt//bin"}, Prepend, "/usr/bin:/opt/bin/"},
		{"prepend deduplicates arguments", "/bin", []string{"/a", "/a/"}, Prepend, "/a:/bin"},
		{"prepend to empty", "", []string{"/a"}, Prepend, "/a"},
		{"append new", "/usr/bin", []string{"/opt/bin"}, Append, "/usr/bin:/opt/bin"},
		{"append existing is a no-op", "/opt/bin:/usr/bin", []string{"/opt/bin"}, Append, "/opt/bin:/usr/bin"},
		{"append deduplicates arguments", "/bin", []string{"/a", "/a"}, Append, "/bin:/a"},
		{"remove all occurrences", "/a:/usr/bin:/a/:/bin", []string{"/a"}, Remove, "/usr/bin:/bin"},
		{"remove missing is a no-op", "/usr/bin", []string{"/a"}, Remove, "/usr/bin"},
		{"remove keeps empty entries", "/a::/bin", []string{"/a"}, Remove, ":/bin"},
		{"empty dirs are ignored", "/bin", []string{""}, Prepend, "/bin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := munge(tt.path, tt.dirs, tt.mode); got != tt.want {
				t.Errorf("munge(%q, %q) = %q, want %q", tt.path, tt.dirs, got, tt.want)
			}
		})
	}
}
