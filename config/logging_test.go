package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTmplLookupEnv(t *testing.T) {
	t.Run("returns the env var when set", func(t *testing.T) {
		t.Setenv("GEXTEND_BASH_TEST_VAR", "value")
		got := TmplLookupEnv("GEXTEND_BASH_TEST_VAR")
		if got != "value" {
			t.Errorf("TmplLookupEnv() = %q, want %q", got, "value")
		}
	})

	t.Run("returns the default when unset", func(t *testing.T) {
		os.Unsetenv("GEXTEND_BASH_TEST_VAR_UNSET")
		got := TmplLookupEnv("GEXTEND_BASH_TEST_VAR_UNSET", "fallback")
		if got != "fallback" {
			t.Errorf("TmplLookupEnv() = %q, want %q", got, "fallback")
		}
	})

	t.Run("returns empty string when unset and no default given", func(t *testing.T) {
		os.Unsetenv("GEXTEND_BASH_TEST_VAR_UNSET")
		got := TmplLookupEnv("GEXTEND_BASH_TEST_VAR_UNSET")
		if got != "" {
			t.Errorf("TmplLookupEnv() = %q, want empty string", got)
		}
	})
}

func newTestConfigLogging(t *testing.T) *ConfigLogging {
	t.Helper()
	dir := t.TempDir()
	t.Setenv(ConfigDirEnv, dir)

	cl := &ConfigLogging{}
	cl.Config = *NewConfig()
	cl.AppName = cl.Config.AppName
	cl.ConfigDir = cl.Config.ConfigDir
	return cl
}

func TestConfigLoggingParse(t *testing.T) {
	cl := newTestConfigLogging(t)

	t.Run("static template returns the same string", func(t *testing.T) {
		got, err := cl.Parse("plain-text")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}
		if got != "plain-text" {
			t.Errorf("Parse() = %q, want %q", got, "plain-text")
		}
	})

	t.Run("field substitution", func(t *testing.T) {
		got, err := cl.Parse("{{ .AppName }}")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}
		if got != cl.AppName {
			t.Errorf("Parse() = %q, want %q", got, cl.AppName)
		}
	})

	t.Run("env template function", func(t *testing.T) {
		t.Setenv("GEXTEND_BASH_TEST_PARSE_VAR", "envvalue")
		got, err := cl.Parse(`{{ env "GEXTEND_BASH_TEST_PARSE_VAR" }}`)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}
		if got != "envvalue" {
			t.Errorf("Parse() = %q, want %q", got, "envvalue")
		}
	})

	t.Run("invalid template returns an error", func(t *testing.T) {
		_, err := cl.Parse("{{ .Missing")
		if err == nil {
			t.Fatal("expected an error for a malformed template")
		}
	})
}

func TestSectLoggingInitialize(t *testing.T) {
	sl := &SectLogging{}
	sl.Initialize()

	if sl.OutputDir != "~/Logs" {
		t.Errorf("OutputDir = %q, want %q", sl.OutputDir, "~/Logs")
	}
	if sl.OutputFile != "common.log" {
		t.Errorf("OutputFile = %q, want %q", sl.OutputFile, "common.log")
	}
	if sl.FileMode != ConstFileMode {
		t.Errorf("FileMode = %v, want %v", sl.FileMode, ConstFileMode)
	}
	if sl.MaxLines != 10000 {
		t.Errorf("MaxLines = %d, want %d", sl.MaxLines, 10000)
	}
}

func TestNewSectLogging(t *testing.T) {
	sl := NewSectLogging()
	if sl == nil {
		t.Fatal("NewSectLogging() returned nil")
	}
	if sl.OutputDir != "~/Logs" {
		t.Errorf("OutputDir = %q, want %q", sl.OutputDir, "~/Logs")
	}
}

func TestConfigLoggingOutputPaths(t *testing.T) {
	cl := newTestConfigLogging(t)
	cl.SectLogging = *NewSectLogging()

	home, err := GetHomeDir()
	if err != nil {
		t.Fatalf("GetHomeDir() error = %v", err)
	}

	gotDir, err := cl.GetOutputDir()
	if err != nil {
		t.Fatalf("GetOutputDir() error = %v", err)
	}
	wantDir := filepath.Join(home, "Logs")
	if gotDir != wantDir {
		t.Errorf("GetOutputDir() = %q, want %q", gotDir, wantDir)
	}

	gotFile, err := cl.GetOutputFile()
	if err != nil {
		t.Fatalf("GetOutputFile() error = %v", err)
	}
	if gotFile != "common.log" {
		t.Errorf("GetOutputFile() = %q, want %q", gotFile, "common.log")
	}

	gotPath, err := cl.LogfilePath()
	if err != nil {
		t.Fatalf("LogfilePath() error = %v", err)
	}
	wantPath := filepath.Join(wantDir, "common.log")
	if gotPath != wantPath {
		t.Errorf("LogfilePath() = %q, want %q", gotPath, wantPath)
	}
}

func TestConfigLoggingWriteReadRoundTrip(t *testing.T) {
	cl := newTestConfigLogging(t)
	cl.ConfigFile = filepath.Join(cl.ConfigDir, "logging.ini")
	cl.SectLogging = *NewSectLogging()
	cl.SectLogging.OutputDir = "~/CustomLogs"
	cl.SectLogging.MaxLines = 42

	if cl.FileExists() {
		t.Fatal("expected the config file not to exist yet")
	}

	if err := cl.Write(); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if !cl.FileExists() {
		t.Fatal("expected the config file to exist after Write()")
	}

	info, err := os.Stat(cl.ConfigFile)
	if err != nil {
		t.Fatalf("os.Stat() error = %v", err)
	}
	if info.Mode().Perm() != os.FileMode(ConstFileMode) {
		t.Errorf("file mode = %v, want %v", info.Mode().Perm(), os.FileMode(ConstFileMode))
	}

	readBack := &ConfigLogging{}
	readBack.ConfigFile = cl.ConfigFile
	if err := readBack.Read(); err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if readBack.OutputDir != "~/CustomLogs" {
		t.Errorf("OutputDir = %q, want %q", readBack.OutputDir, "~/CustomLogs")
	}
	if readBack.MaxLines != 42 {
		t.Errorf("MaxLines = %d, want %d", readBack.MaxLines, 42)
	}
}

func TestConfigLoggingCreateIfEmpty(t *testing.T) {
	cl := newTestConfigLogging(t)
	cl.ConfigFile = filepath.Join(cl.ConfigDir, "logging.ini")

	t.Run("creates the file when missing", func(t *testing.T) {
		if err := cl.CreateIfEmpty(); err != nil {
			t.Fatalf("CreateIfEmpty() error = %v", err)
		}
		if !cl.FileExists() {
			t.Fatal("expected the config file to have been created")
		}
	})

	t.Run("is a no-op when the file already exists", func(t *testing.T) {
		before, err := os.ReadFile(cl.ConfigFile)
		if err != nil {
			t.Fatalf("os.ReadFile() error = %v", err)
		}
		if err := cl.CreateIfEmpty(); err != nil {
			t.Fatalf("CreateIfEmpty() error = %v", err)
		}
		after, err := os.ReadFile(cl.ConfigFile)
		if err != nil {
			t.Fatalf("os.ReadFile() error = %v", err)
		}
		if string(before) != string(after) {
			t.Error("expected the existing config file to be left untouched")
		}
	})
}

func TestConfigLoggingInitialize(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(ConfigDirEnv, dir)

	cl := &ConfigLogging{}
	cl.Initialize()

	if cl.AppName != ApplicationName {
		t.Errorf("AppName = %q, want %q", cl.AppName, ApplicationName)
	}
	if cl.ConfigDir != dir {
		t.Errorf("ConfigDir = %q, want %q", cl.ConfigDir, dir)
	}
	wantConfigFile := filepath.Join(dir, "logging.ini")
	if cl.ConfigFile != wantConfigFile {
		t.Errorf("ConfigFile = %q, want %q", cl.ConfigFile, wantConfigFile)
	}
	if !cl.FileExists() {
		t.Error("expected Initialize() to have written the config file")
	}
	if cl.Hostname == "" {
		t.Error("expected Hostname to be populated")
	}
}

func TestNewConfigLogging(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(ConfigDirEnv, dir)

	cl := NewConfigLogging()
	if cl == nil {
		t.Fatal("NewConfigLogging() returned nil")
	}
	if cl.Now.IsZero() {
		t.Error("expected Now to be populated")
	}
}
