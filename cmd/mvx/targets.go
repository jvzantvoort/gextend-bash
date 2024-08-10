package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// Get the size of a file and handle the error(ish)
func GetFileSize(target string) int64 {
	if fi, ok := os.Stat(target); ok == nil {
		return fi.Size()
	}
	return int64(0)

}

func deepCompare(srcfile, dstfile string) bool {
	StartFunc(CurFunc())
	defer EndFunc(CurFunc())

	// test on size first
	srcfile_size := GetFileSize(srcfile)
	dstfile_size := GetFileSize(dstfile)
	if srcfile_size != dstfile_size {
		LogDebugf("file sizes differ")
		return false
	}

	sf, err := os.Open(srcfile)
	if err != nil {
		log.Fatal(err)
	}

	df, err := os.Open(dstfile)
	if err != nil {
		log.Fatal(err)
	}

	sscan := bufio.NewScanner(sf)
	dscan := bufio.NewScanner(df)

	for sscan.Scan() {
		dscan.Scan()
		if !bytes.Equal(sscan.Bytes(), dscan.Bytes()) {
			return false
		}
	}

	return true
}

func DirectoryExists(target string) bool {
	info, err := os.Stat(target)
	if os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		return true
	}
	return false
}

// We need the inverse ...
func TargetExistsAndIsNotADirectory(target string) bool {

	info, err := os.Stat(target)
	if os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

// Make sure the directory exists or exit with an error
func ensureDir(dirname string) error {
	StartFunc(CurFunc())
	defer EndFunc(CurFunc())

	if DirectoryExists(dirname) {
		LogDebugf("directory exists")
		return nil
	}

	_, err := os.Stat(dirname)
	if !os.IsNotExist(err) {
		LogDebugf("target exists but is not a directory")
		return ErrNotADir
	}

	err = os.MkdirAll(dirname, 0755)

	if err == nil {
		return nil
	}

	if os.IsExist(err) {
		return nil
	}

	LogErrorf("%s", err)

	return err
}

func targetExists(target string) bool {
	StartFunc(CurFunc())
	defer EndFunc(CurFunc())
	msg := fmt.Sprintf("target found %s", lastN(target, 50))

	_, err := os.Stat(target)
	if err == nil {
		LogSuccess(msg)
		return true
	}

	if os.IsNotExist(err) {
		LogFailure(msg)
		return false
	}

	LogErrorf("%s: %s", target, err)
	return false

}
