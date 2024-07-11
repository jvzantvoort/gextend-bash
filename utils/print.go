package utils

import (
	log "github.com/sirupsen/logrus"
)

// PrintError if err is not nil print fmtstr as error.
func PrintError(fmtstr string, err error) error {
	if err == nil {
		return err
	}
	log.Errorf(fmtstr, err)
	return err
}

// PrintFatal missing godoc.
func PrintFatal(fmtstr string, err error) error {
	if err == nil {
		return err
	}
	log.Fatalf(fmtstr, err)
	return err
}

// PanicOnError missing godoc.
func PanicOnError(fmtstr string, err error) {
	PrintError(fmtstr, err)
	if err != nil {
		panic(err)
	}
}
