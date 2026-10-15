package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

const (
	CONSOLE_WIDTH = 180
)

var (
	indentchar    string = "  "
	indent        int    = 0
	debug         bool   = false
	ErrNotADir           = errors.New("target exists but is not a directory")
	ErrMaxCount          = errors.New("max count exceeded")
	ErrSameFile          = errors.New("sourcefile and target are the same")
	ErrSrcNoExist        = errors.New("sourcefile does not exist")
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          false,
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		//		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func formatDestFile(destdir, sourcefile string, num int) string {
	StartFunc(CurFunc())
	defer EndFunc(CurFunc())

	sourcefile = filepath.Base(sourcefile) // strip away leading gunc
	LogDebugf("sourcefile: %s", sourcefile)

	if num == 0 {
		target := filepath.Join(destdir, filepath.Base(sourcefile))
		LogDebugf("return: %s", target)
		return target
	}

	ext := filepath.Ext(sourcefile)

	LogDebugf("ext: %s", ext)
	if len(ext) == 0 {
		return filepath.Join(destdir, fmt.Sprintf("%s.%d", sourcefile, num))
	}

	ext_indx := len(sourcefile) - len(ext)
	LogDebugf("%s", sourcefile[:ext_indx])
	basename := sourcefile[:ext_indx]
	return filepath.Join(destdir, fmt.Sprintf("%s.%d.%s", basename, num, ext[1:]))

}

func GetNextTarget(destdir, sourcefile string) (string, error) {
	StartFunc(CurFunc())
	defer EndFunc(CurFunc())

	target := formatDestFile(destdir, sourcefile, 0)

	if !targetExists(target) {
		LogDebugf("target does not exist, returning %s", target)
		return target, nil
	}

	// targets are the same
	if deepCompare(sourcefile, target) {
		return "", ErrSameFile
	}

	for i := 1; i < 100; i++ {
		target = formatDestFile(destdir, sourcefile, i)
		if !targetExists(target) {
			LogDebugf("return %s", target)
			return target, nil
		}
	}
	return "", ErrMaxCount
}

func MoveFile(src, dst string) error {
	StartFunc(CurFunc())
	defer EndFunc(CurFunc())

	// translate src to absolute
	src, err := filepath.Abs(src)
	if err != nil {
		LogErrorf("cannot get abspath %q", err)
		return err
	}
	LogDebugf("sourcefile %s", src)
	LogDebugf("destdir %s", dst)

	// fail if the source file does not exist
	if !targetExists(src) {
		return ErrSrcNoExist
	} else {
		LogDebugf("sourcefile %s exists", src)
	}

	// src_shortname := filepath.Base(src)
	err = ensureDir(dst)
	if err != nil {
		LogErrorf("%q", err)
		return err
	}

	dst_path, err := GetNextTarget(dst, src)
	if errors.Is(err, ErrSameFile) {
		return os.Remove(src)
	}
	if err != nil {
		return err
	}

	LogDebugf("mv %s %s", src, dst_path)
	err = os.Rename(src, dst_path)
	if errors.Is(err, syscall.EXDEV) {
		// different filesystem: rename cannot work, copy and remove instead
		LogDebugf("cross-device move, copying %s", src)
		if err = copyFile(src, dst_path); err != nil {
			return err
		}
		return os.Remove(src)
	}

	return err
}

// copyFile copies src to dst, preserving the permission bits and the
// modification time. A partially written dst is removed on failure.
func copyFile(src, dst string) (err error) {
	StartFunc(CurFunc())
	defer EndFunc(CurFunc())

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, info.Mode().Perm())
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			_ = os.Remove(dst)
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	if err = out.Sync(); err != nil {
		return err
	}
	return os.Chtimes(dst, info.ModTime(), info.ModTime())
}

func main() {
	appname := path.Base(os.Args[0])
	StartFunc(appname)
	defer EndFunc(appname)

	flag.BoolVar(&debug, "debug", false, "debug messages")
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}
	arguments := flag.Args()
	arg_indx := len(arguments)

	if arg_indx < 2 {
		fmt.Printf("\n\n\tNot enough arguments: %d\n\n", arg_indx)
		os.Exit(1)
	}

	DestTarget := arguments[arg_indx-1]
	SourceTargets := arguments[:arg_indx-1]

	LogDebugf("Sources: %s", strings.Join(SourceTargets, ", "))
	LogDebugf("Destination: %s", DestTarget)

	DestTarget, err := filepath.Abs(DestTarget)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", appname, err)
		os.Exit(1)
	}

	failed := 0
	for _, indx := range SourceTargets {
		msg := fmt.Sprintf("move %s", lastN(indx, 60))
		err := MoveFile(indx, DestTarget)
		if err == nil {
			LogSuccess("%s", msg)
		} else {
			LogFailure("%s", msg)
			fmt.Fprintf(os.Stderr, "%s: %s: %s\n", appname, indx, err)
			failed++
		}
	}

	if failed > 0 {
		os.Exit(1)
	}
}
