package main

import (
	"bytes"
	"io"
	"os"
)

// compareChunkSize is the block size deepCompare reads at a time.
const compareChunkSize = 64 * 1024

// Get the size of a file and handle the error(ish)
func GetFileSize(target string) int64 {
	if fi, ok := os.Stat(target); ok == nil {
		return fi.Size()
	}
	return int64(0)

}

// deepCompare reports whether srcfile and dstfile have identical content.
// Any error while reading either file is treated as "not identical", so the
// caller never discards a file it could not fully verify.
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
		LogErrorf("%s", err)
		return false
	}
	defer func() { _ = sf.Close() }()

	df, err := os.Open(dstfile)
	if err != nil {
		LogErrorf("%s", err)
		return false
	}
	defer func() { _ = df.Close() }()

	sbuf := make([]byte, compareChunkSize)
	dbuf := make([]byte, compareChunkSize)
	for {
		sn, serr := io.ReadFull(sf, sbuf)
		dn, derr := io.ReadFull(df, dbuf)
		if !bytes.Equal(sbuf[:sn], dbuf[:dn]) {
			return false
		}
		sdone := serr == io.EOF || serr == io.ErrUnexpectedEOF
		ddone := derr == io.EOF || derr == io.ErrUnexpectedEOF
		if (serr != nil && !sdone) || (derr != nil && !ddone) {
			LogErrorf("compare %s %s: %v %v", srcfile, dstfile, serr, derr)
			return false
		}
		if sdone || ddone {
			return sdone && ddone
		}
	}
}

func DirectoryExists(target string) bool {
	info, err := os.Stat(target)
	if err != nil {
		return false
	}
	return info.IsDir()
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

	_, err := os.Stat(target)
	if err == nil {
		LogDebugf("target found %s", target)
		return true
	}

	if os.IsNotExist(err) {
		LogDebugf("target not found %s", target)
		return false
	}

	LogErrorf("%s: %s", target, err)
	return false
}
