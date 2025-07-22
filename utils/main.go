// Package utils provides general utility functions for file, directory, and system operations.
package utils

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
)

// ShortHostname returns the short hostname (first part of the FQDN) in lowercase.
func ShortHostname() string {
	fqdn, _ := os.Hostname()
	parts := strings.Split(fqdn, ".")
	return strings.ToLower(parts[0])
}

// GetHomeDir returns the current user's home directory.
func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

// MkdirP creates a directory and all necessary parents with the specified mode.
// Returns an error if the directory cannot be created or if a non-directory file exists at the path.
//
//	err := utils.MkdirP("/lala", int(0755))
//	if err != nil {
//	  panic(err)
//	}
func MkdirP(dirname string, mode int) error {

	target_stat, err := os.Stat(dirname)
	if err == nil {
		if target_stat.IsDir() {
			return nil
		} else {
			return fmt.Errorf("target exists %s but is not a directory", dirname)
		}
	}

	if err := os.MkdirAll(dirname, os.FileMode(mode)); err != nil {
		return fmt.Errorf("directory cannot be created: %s", dirname)
	}
	return nil
}

// FileExists checks if the target exists and is a file.
// Returns true and the file info if it exists and is a file, otherwise false.
//
//	check, info := utils.FileExists("/etc/passwd")
//	if check {
//	   fmt.Printf("size: %d\n", info.Size())
//	}
func FileExists(fpath string) (bool, os.FileInfo) {
	info, err := os.Stat(fpath)
	if err != nil {
		return false, info
	}

	// is a directory
	if info.IsDir() {
		return false, info
	}

	return true, info
}

// FileIsExecutable checks if the file exists and is executable by owner, group, or others.
// On Windows, returns true if the file exists.
func FileIsExecutable(fpath string) bool {
	exists, info := FileExists(fpath)
	if !exists {
		return false
	}

	goos := runtime.GOOS

	// windows doesn't do that
	if goos == "windows" {
		return true
	}

	mode := info.Mode()

	// Exec owner
	if mode&0100 != 0 {
		return true
	}

	// Exec group
	if mode&0010 != 0 {
		return true
	}

	// Exec other
	if mode&0001 != 0 {
		return true
	}
	return false
}
