package utils

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
)

// GetHomeDir get the user's homedir
func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

// MkdirP
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
			return fmt.Errorf("Target exists %s but is not a directory", dirname)
		}
	}

	if err := os.MkdirAll(dirname, os.FileMode(mode)); err != nil {
		return fmt.Errorf("Directory cannot be created: %s", dirname)
	}
	return nil
}

// FileExists check if a target exists and is a file.
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

// FileIsExecutable file exists and is executable
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
