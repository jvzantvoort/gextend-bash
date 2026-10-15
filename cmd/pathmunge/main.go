// pathmunge adds or removes entries in a colon separated path variable
// without creating duplicates, and prints the result:
//
//	PATH="$(pathmunge /opt/tool/bin)"
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jvzantvoort/gextend-bash/version"
)

// Mode selects what munge does with the given directories.
type Mode int

const (
	Prepend Mode = iota
	Append
	Remove
)

// munge returns path with dirs prepended, appended or removed. Directories
// already present are left where they are (like the classic Red Hat
// pathmunge); entries are compared after filepath.Clean. Empty entries of
// path are kept as they are.
func munge(path string, dirs []string, mode Mode) string {
	var entries []string
	if path != "" {
		entries = strings.Split(path, ":")
	}

	present := func(list []string, dir string) bool {
		for _, entry := range list {
			if entry != "" && filepath.Clean(entry) == dir {
				return true
			}
		}
		return false
	}

	var cleaned []string
	for _, dir := range dirs {
		if dir != "" {
			cleaned = append(cleaned, filepath.Clean(dir))
		}
	}

	switch mode {
	case Remove:
		var kept []string
		for _, entry := range entries {
			if entry == "" || !present(cleaned, filepath.Clean(entry)) {
				kept = append(kept, entry)
			}
		}
		entries = kept
	case Append:
		for _, dir := range cleaned {
			if !present(entries, dir) {
				entries = append(entries, dir)
			}
		}
	default:
		// keep the order the directories were given in
		var head []string
		for _, dir := range cleaned {
			if !present(entries, dir) && !present(head, dir) {
				head = append(head, dir)
			}
		}
		entries = append(head, entries...)
	}

	return strings.Join(entries, ":")
}

func main() {
	var after, remove, showVersion bool
	var variable string
	flag.BoolVar(&after, "a", false, "append instead of prepend")
	flag.BoolVar(&after, "after", false, "same as -a")
	flag.BoolVar(&remove, "r", false, "remove the directories")
	flag.BoolVar(&remove, "remove", false, "same as -r")
	flag.StringVar(&variable, "v", "PATH", "environment variable to read")
	flag.BoolVar(&showVersion, "version", false, "print the version and exit")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-a|-r] [-v VAR] dir [dir...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Prints $VAR (default PATH) with the directories prepended (default),\n")
		fmt.Fprintf(os.Stderr, "appended (-a) or removed (-r), without creating duplicates:\n\n")
		fmt.Fprintf(os.Stderr, "  PATH=\"$(%s /opt/tool/bin)\"\n\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if showVersion {
		version.Print(os.Stdout, "pathmunge")
		return
	}
	if flag.NArg() == 0 || (after && remove) {
		flag.Usage()
		os.Exit(2)
	}

	mode := Prepend
	switch {
	case after:
		mode = Append
	case remove:
		mode = Remove
	}

	fmt.Println(munge(os.Getenv(variable), flag.Args(), mode))
}
