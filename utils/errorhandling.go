// Package utils provides error handling utilities for the gextend-bash project.
package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// WarningOnError logs a warning if the provided error is not nil.
func WarningOnError(err error) {
	if err != nil {
		log.Warningf("error %v\n", err)
	}
}

// ExitOnError checks the error and exits the program if the error is not nil.
func ExitOnError(err error) {
	if err != nil {
		log.Errorf("error %v\n", err)
		os.Exit(1)
	}
}
