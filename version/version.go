// Package version reports the version of the gextend-bash tools.
//
// Release builds set Version through the linker:
//
//	go build -ldflags "-X github.com/jvzantvoort/gextend-bash/version.Version=1.2.3"
//
// Without it, the version is taken from the module and VCS information Go
// embeds in every binary (go install ...@version, or a build from a git
// checkout).
package version

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
)

// Version is set at build time; see the package documentation.
var Version = ""

// readBuildInfo is replaceable for tests.
var readBuildInfo = debug.ReadBuildInfo

// String returns a human readable version, e.g. "1.2.3" or
// "v0.0.0-20261015-abcdef (rev abcdef0, modified)".
func String() string {
	if Version != "" {
		return Version
	}

	info, ok := readBuildInfo()
	if !ok {
		return "unknown"
	}

	retv := info.Main.Version
	if retv == "" {
		retv = "(devel)"
	}

	var revision string
	var modified bool
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	if revision != "" {
		if len(revision) > 7 {
			revision = revision[:7]
		}
		retv = fmt.Sprintf("%s (rev %s", retv, revision)
		if modified {
			retv += ", modified"
		}
		retv += ")"
	}
	return retv
}

// IsFlag reports whether arg asks for the version.
func IsFlag(arg string) bool {
	return arg == "-version" || arg == "--version"
}

// Print writes "<name> <version>" to w.
func Print(w io.Writer, name string) {
	_, _ = fmt.Fprintf(w, "%s %s\n", name, String())
}

// HandleArgs prints the version and exits when the first argument asks for
// it. It is meant for tools that take their arguments without a flag parser.
func HandleArgs(args []string) {
	if len(args) > 1 && IsFlag(args[1]) {
		Print(os.Stdout, filepath.Base(args[0]))
		os.Exit(0)
	}
}
