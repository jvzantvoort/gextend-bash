package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jvzantvoort/gextend-bash/version"
)

func AppendIfMissing(slice []string, i string) []string {
	for _, element := range slice {
		if element == i {
			return slice
		}
	}
	return append(slice, i)
}

// FilterExists returns the existing directories in slice, cleaned and
// without duplicates. Empty and relative entries (which resolve against
// whatever the current directory happens to be) are dropped.
func FilterExists(slice []string) []string {
	var retv []string
	for _, dirn := range slice {
		if !filepath.IsAbs(dirn) {
			continue
		}
		dirn = filepath.Clean(dirn)
		info, err := os.Stat(dirn)
		if err != nil {
			continue
		}
		if info.IsDir() {
			retv = AppendIfMissing(retv, dirn)
		}
	}
	return retv
}

func main() {
	version.HandleArgs(os.Args)

	path := strings.Join(os.Args[1:], " ")

	// Default to the environment variable
	if len(path) == 0 {
		path = os.Getenv("PATH")
	}

	pathlist := FilterExists(strings.Split(path, ":"))
	fmt.Printf("%s\n", strings.Join(pathlist, ":"))

}
