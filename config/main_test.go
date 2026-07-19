package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetHomeDir(t *testing.T) {
	got, err := GetHomeDir()
	if err != nil {
		t.Fatalf("GetHomeDir() error = %v", err)
	}
	if got == "" {
		t.Error("GetHomeDir() returned an empty string")
	}
}

func TestExpandHome(t *testing.T) {
	home, err := GetHomeDir()
	if err != nil {
		t.Fatalf("GetHomeDir() error = %v", err)
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty path", "", ""},
		{"absolute path unchanged", "/etc/passwd", "/etc/passwd"},
		{"relative path unchanged", "relative/path", "relative/path"},
		{"tilde alone", "~", home},
		{"tilde with subpath", "~/foo/bar", filepath.Join(home, "foo", "bar")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandHome(tt.input)
			if err != nil {
				t.Fatalf("ExpandHome(%q) error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ExpandHome(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestConfigSetDefaultHomeDir(t *testing.T) {
	t.Run("populates an empty HomeDir", func(t *testing.T) {
		c := &Config{}
		if err := c.SetDefaultHomeDir(); err != nil {
			t.Fatalf("SetDefaultHomeDir() error = %v", err)
		}
		if c.HomeDir == "" {
			t.Error("expected HomeDir to be populated")
		}
	})

	t.Run("leaves an existing HomeDir untouched", func(t *testing.T) {
		c := &Config{HomeDir: "/already/set"}
		if err := c.SetDefaultHomeDir(); err != nil {
			t.Fatalf("SetDefaultHomeDir() error = %v", err)
		}
		if c.HomeDir != "/already/set" {
			t.Errorf("HomeDir = %q, want %q", c.HomeDir, "/already/set")
		}
	})
}

func TestConfigSetDefaultConfigDir(t *testing.T) {
	t.Run("uses the env var when set", func(t *testing.T) {
		dir := t.TempDir()
		envDir := filepath.Join(dir, "from-env")
		t.Setenv(ConfigDirEnv, envDir)

		c := &Config{}
		c.SetDefaultConfigDir()
		if c.ConfigDir != envDir {
			t.Errorf("ConfigDir = %q, want %q", c.ConfigDir, envDir)
		}
	})

	t.Run("falls back to ~/.config and creates it", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv(ConfigDirEnv, "")
		os.Unsetenv(ConfigDirEnv)

		c := &Config{HomeDir: home}
		c.SetDefaultConfigDir()

		want := filepath.Join(home, ".config", ConfigDirName)
		if c.ConfigDir != want {
			t.Errorf("ConfigDir = %q, want %q", c.ConfigDir, want)
		}
		info, err := os.Stat(c.ConfigDir)
		if err != nil {
			t.Fatalf("expected ConfigDir to be created: %v", err)
		}
		if !info.IsDir() {
			t.Error("expected ConfigDir to be a directory")
		}
	})

	t.Run("is a no-op when ConfigDir already set", func(t *testing.T) {
		c := &Config{ConfigDir: "/already/set"}
		c.SetDefaultConfigDir()
		if c.ConfigDir != "/already/set" {
			t.Errorf("ConfigDir = %q, want %q", c.ConfigDir, "/already/set")
		}
	})
}

func TestConfigInitialize(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(ConfigDirEnv, dir)

	c := &Config{}
	c.Initialize()

	if c.AppName != ApplicationName {
		t.Errorf("AppName = %q, want %q", c.AppName, ApplicationName)
	}
	if c.HomeDir == "" {
		t.Error("expected HomeDir to be populated")
	}
	if c.ConfigDir != dir {
		t.Errorf("ConfigDir = %q, want %q", c.ConfigDir, dir)
	}
}

func TestNewConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(ConfigDirEnv, dir)

	c := NewConfig()
	if c == nil {
		t.Fatal("NewConfig() returned nil")
	}
	if c.ConfigDir != dir {
		t.Errorf("ConfigDir = %q, want %q", c.ConfigDir, dir)
	}
}
