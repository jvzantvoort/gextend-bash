package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GetString(cmd cobra.Command, name string) string {
	retv, _ := cmd.Flags().GetString(name)
	if len(retv) != 0 {
		log.Infof("Found %s as %s", name, retv)
	}
	return retv
}
