package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

func TestAsk(t *testing.T) {
	tests := []struct {
		name  string
		input string
		def   Default
		want  bool
		hint  string
	}{
		{"yes", "y\n", NoDefault, true, "[y/n]"},
		{"YES in capitals", "YES\n", NoDefault, true, "[y/n]"},
		{"no", "no\n", DefaultYes, false, "[Y/n]"},
		{"empty takes default yes", "\n", DefaultYes, true, "[Y/n]"},
		{"empty takes default no", "\n", DefaultNo, false, "[y/N]"},
		{"empty without default asks again", "\ny\n", NoDefault, true, "[y/n]"},
		{"invalid asks again", "maybe\nn\n", NoDefault, false, "[y/n]"},
		{"end of input without default is no", "", NoDefault, false, "[y/n]"},
		{"end of input takes the default", "", DefaultYes, true, "[Y/n]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			got := ask(strings.NewReader(tt.input), &out, "Continue?", tt.def, 0)
			if got != tt.want {
				t.Errorf("ask() = %v, want %v (output %q)", got, tt.want, out.String())
			}
			if !strings.HasPrefix(out.String(), "Continue? "+tt.hint) {
				t.Errorf("prompt = %q, want hint %s", out.String(), tt.hint)
			}
		})
	}
}

func TestAskTimeout(t *testing.T) {
	r, w := io.Pipe()
	defer func() { _ = w.Close() }()

	for _, tt := range []struct {
		def  Default
		want bool
	}{{NoDefault, false}, {DefaultYes, true}} {
		var out bytes.Buffer
		start := time.Now()
		if got := ask(r, &out, "Continue?", tt.def, 50*time.Millisecond); got != tt.want {
			t.Errorf("ask() with default %v after timeout = %v, want %v", tt.def, got, tt.want)
		}
		if time.Since(start) > 2*time.Second {
			t.Error("ask() did not honour the timeout")
		}
	}
}

func TestParseDefault(t *testing.T) {
	for in, want := range map[string]Default{"": NoDefault, "y": DefaultYes, "Yes": DefaultYes, "n": DefaultNo, "NO": DefaultNo} {
		got, err := parseDefault(in)
		if err != nil || got != want {
			t.Errorf("parseDefault(%q) = %v, %v; want %v", in, got, err, want)
		}
	}
	if _, err := parseDefault("maybe"); err == nil {
		t.Error("parseDefault(maybe) expected an error")
	}
}
