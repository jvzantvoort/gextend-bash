// Package utils provides utility functions for error printing and handling.
package utils

import (
	log "github.com/sirupsen/logrus"
)

// PrintError prints an error message using the provided format string if err is not nil.
// Returns the error for further handling.
func PrintError(fmtstr string, err error) error {
	if err == nil {
		return err
	}
	log.Errorf(fmtstr, err)
	return err
}

// PrintFatal prints a fatal error message and exits if err is not nil.
// Returns the error for further handling.
func PrintFatal(fmtstr string, err error) error {
	if err == nil {
		return err
	}
	log.Fatalf(fmtstr, err)
	return err
}

// PanicOnError prints an error and panics if err is not nil.
func PanicOnError(fmtstr string, err error) {
	err = PrintError(fmtstr, err)
	if err != nil {
		panic(err)
	}
}
