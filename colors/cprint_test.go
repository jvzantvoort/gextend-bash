package colors

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() error = %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}
	return buf.String()
}

func TestNewCprint(t *testing.T) {
	c := NewCprint()
	if c == nil {
		t.Fatal("NewCprint() returned nil")
	}
	if len(c.Colors) == 0 {
		t.Error("expected Colors map to be populated")
	}
	if c.Colors["ok"] != color.FgGreen {
		t.Errorf("Colors[ok] = %v, want %v", c.Colors["ok"], color.FgGreen)
	}
	if c.Colors["nok"] != color.FgRed {
		t.Errorf("Colors[nok] = %v, want %v", c.Colors["nok"], color.FgRed)
	}
}

func TestCprintSetColor(t *testing.T) {
	color.NoColor = false

	t.Run("known color adds its SGR attribute", func(t *testing.T) {
		c := NewCprint()
		c.SetColor("green")
		got := c.color.Sprint("x")
		want := strconv.Itoa(int(color.FgGreen))
		if !strings.Contains(got, want) {
			t.Errorf("expected output %q to contain the FgGreen SGR code %s", got, want)
		}
	})

	t.Run("unknown color leaves the color unset", func(t *testing.T) {
		c := NewCprint()
		c.SetColor("does-not-exist")
		got := c.color.Sprint("x")
		if got != "\x1b[mx\x1b[m" {
			t.Errorf("expected no color attributes to be set, got %q", got)
		}
	})
}

func TestCprintSetToken(t *testing.T) {
	tests := []struct {
		name string
		want rune
	}{
		{"ok", HEAVY_CHECK_MARK},
		{"oke", HEAVY_CHECK_MARK},
		{"nok", WARNING_SIGN},
		{"warn", EXCLAMATION_MARK},
		{"profile", BLACK_DIAMOND_SUIT},
		{"unknown-token", BLACK_DIAMOND_SUIT},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCprint()
			c.SetToken(tt.name)
			if c.token != tt.want {
				t.Errorf("SetToken(%q): token = %q, want %q", tt.name, c.token, tt.want)
			}
		})
	}
}

func TestCprintSetFormat(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"profile", "%s profile (%s) sourced\n"},
		{"platform", "%s platform %s sourced\n"},
		{"workspace", "%s workspace %s sourced\n"},
		{"unknown", "%s %s\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCprint()
			c.SetFormat(tt.name)
			if c.Format != tt.want {
				t.Errorf("SetFormat(%q): Format = %q, want %q", tt.name, c.Format, tt.want)
			}
		})
	}
}

func TestCprintPrint(t *testing.T) {
	color.NoColor = true
	c := NewCprint()

	out := captureStdout(t, func() {
		c.Print("ok", "everything", "is", "fine")
	})

	if !strings.Contains(out, "everything is fine") {
		t.Errorf("expected output to contain the joined message, got %q", out)
	}
}

func TestCprintPrintProfileFormat(t *testing.T) {
	color.NoColor = true
	c := NewCprint()

	out := captureStdout(t, func() {
		c.Print("profile", "myprofile")
	})

	if !strings.Contains(out, "profile (myprofile) sourced") {
		t.Errorf("unexpected output: %q", out)
	}
}
