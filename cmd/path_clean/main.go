package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func AppendIfMissing(slice []string, i string) []string {
	for _, element := range slice {
		if element == i {
			return slice
		}
	}
	return append(slice, i)
}

func FilterExists(slice []string) []string {
	var retv []string
	for _, dirn := range slice {
		if len(dirn) == 0 {
			continue
		}
		dirn, _ = filepath.Abs(dirn)
		info, err := os.Stat(dirn)
		if os.IsNotExist(err) {
			continue
		}
		if info.IsDir() {
			retv = AppendIfMissing(retv, dirn)
		}
	}
	return retv
}

func main() {
	path := strings.Join(os.Args[1:], " ")

	// Default to the environment variable
	if len(path) == 0 {
		path = os.Getenv("PATH")
	}

	pathlist := FilterExists(strings.Split(path, ":"))
	fmt.Printf("%s\n", strings.Join(pathlist, ":"))

}
