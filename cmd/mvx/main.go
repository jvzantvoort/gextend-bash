package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

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
		FullTimestamp:   false,
		DisableTimestamp: true,
		DisableLevelTruncation: true,
		PadLevelText: true,
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
		return fmt.Sprintf("%s/%s.%d", destdir, sourcefile, num)
	}

	ext_indx := len(sourcefile) - len(ext)
	LogDebugf("%s", sourcefile[:ext_indx])
	basename := sourcefile[:ext_indx]
	return fmt.Sprintf("%s/%s.%d.%s", destdir, basename, num, ext[1:])

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
	if err != nil {
		if errors.Is(err, ErrMaxCount) {
			return err
		}
		if errors.Is(err, ErrSameFile) {
			err = os.Remove(src)
			return err
		}

		return err
	}

	LogDebugf("mv %s %s", src, dst_path)
	err = os.Rename(src, dst_path)
	if err != nil {
		log.Fatal(err)
	}

	return nil
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

	DestTarget := arguments[arg_indx - 1]
	SourceTargets := arguments[:arg_indx -1]

	LogDebugf("Sources: %s", strings.Join(SourceTargets, ", "))
	LogDebugf("Destination: %s", DestTarget)

	DestTarget, _ = filepath.Abs(DestTarget)

	for _, indx := range SourceTargets {
		msg := fmt.Sprintf("move %s", lastN(indx, 60))
		err := MoveFile(indx, DestTarget)
		if err == nil {
			LogSuccess(msg)
		} else {
			LogFailure(msg)
		}
	}
}
