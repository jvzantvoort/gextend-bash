package version

import (
	"bytes"
	"runtime/debug"
	"testing"
)

func withBuildInfo(t *testing.T, info *debug.BuildInfo, ok bool) {
	t.Helper()
	orig := readBuildInfo
	readBuildInfo = func() (*debug.BuildInfo, bool) { return info, ok }
	t.Cleanup(func() { readBuildInfo = orig })
}

func TestString(t *testing.T) {
	t.Run("linker version wins", func(t *testing.T) {
		Version = "1.2.3"
		t.Cleanup(func() { Version = "" })
		if got := String(); got != "1.2.3" {
			t.Errorf("String() = %q, want %q", got, "1.2.3")
		}
	})

	tests := []struct {
		name string
		info *debug.BuildInfo
		ok   bool
		want string
	}{
		{"no build info", nil, false, "unknown"},
		{"module version only", &debug.BuildInfo{Main: debug.Module{Version: "v1.0.0"}}, true, "v1.0.0"},
		{"devel with revision", &debug.BuildInfo{
			Main: debug.Module{Version: "(devel)"},
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "0123456789abcdef"},
				{Key: "vcs.modified", Value: "true"},
			},
		}, true, "(devel) (rev 0123456, modified)"},
		{"empty version", &debug.BuildInfo{
			Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "abc"}},
		}, true, "(devel) (rev abc)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withBuildInfo(t, tt.info, tt.ok)
			if got := String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsFlag(t *testing.T) {
	for arg, want := range map[string]bool{"-version": true, "--version": true, "version": false, "-v": false} {
		if got := IsFlag(arg); got != want {
			t.Errorf("IsFlag(%q) = %v, want %v", arg, got, want)
		}
	}
}

func TestPrint(t *testing.T) {
	Version = "9.9.9"
	t.Cleanup(func() { Version = "" })
	var buf bytes.Buffer
	Print(&buf, "tool")
	if buf.String() != "tool 9.9.9\n" {
		t.Errorf("Print() wrote %q", buf.String())
	}
}
