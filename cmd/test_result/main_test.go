package main

import "testing"

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
