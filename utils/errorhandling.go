package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func WarningOnError(err error) {
	if err != nil {
		log.Warningf("error %v\n", err)
	}
}

// ExitOnError check error and exit if not nil
func ExitOnError(err error) {
	if err != nil {
		log.Errorf("error %v\n", err)
		os.Exit(1)
	}
}
