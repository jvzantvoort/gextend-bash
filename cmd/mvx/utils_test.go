package main

import "testing"

func TestCenterText(t *testing.T) {
	tests := []struct {
		name  string
		instr string
		width int
	}{
		{"short word", "hi", 10},
		{"exact width", "abcdef", 6},
		{"empty string", "", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CenterText(tt.instr, tt.width)
			if len(got) != tt.width {
				t.Errorf("CenterText(%q, %d) length = %d, want %d (got %q)", tt.instr, tt.width, len(got), tt.width, got)
			}
		})
	}
}

func TestLastN(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int
		want   string
	}{
		{"shorter than length returned as-is", "short", 10, "short"},
		{"equal to length returned as-is", "exactly10!", 10, "exactly10!"},
		{"longer than length is truncated with ellipsis", "this is a very long string", 10, "... string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lastN(tt.input, tt.length)
			if got != tt.want {
				t.Errorf("lastN(%q, %d) = %q, want %q", tt.input, tt.length, got, tt.want)
			}
		})
	}
}
