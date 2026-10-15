package main

import (
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantCode int
		wantMsg  string
		wantErr  bool
	}{
		{"success", []string{"0", "all", "good"}, 0, "all good", false},
		{"failure code", []string{"3", "broken"}, 3, "broken", false},
		{"too few arguments", []string{"0"}, 0, "", true},
		{"no arguments", nil, 0, "", true},
		{"non-numeric code", []string{"x", "msg"}, 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, msg, err := parseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseArgs(%q) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
			if code != tt.wantCode || msg != tt.wantMsg {
				t.Errorf("parseArgs(%q) = (%d, %q), want (%d, %q)", tt.args, code, msg, tt.wantCode, tt.wantMsg)
			}
		})
	}
}

func TestRender(t *testing.T) {
	color.NoColor = true

	tests := []struct {
		code  int
		label string
	}{
		{0, "SUCCESS"},
		{1, "FAILURE"},
		{127, "FAILURE"},
	}

	for _, tt := range tests {
		got := render(tt.code, "deploy")
		if !strings.HasPrefix(got, "deploy.") {
			t.Errorf("render(%d) = %q, want it to start with the dotted message", tt.code, got)
		}
		if !strings.HasSuffix(got, "[ "+tt.label+" ]") {
			t.Errorf("render(%d) = %q, want label %s", tt.code, got, tt.label)
		}
	}
}
