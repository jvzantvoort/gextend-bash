package main

import (
	"os"

	"github.com/jvzantvoort/gextend-bash/logging"
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

func handleLogCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) == 0 {
		log.Error("No message provided")
		cmd.Help()
		os.Exit(1)
	}
	logmsg := logging.NewLogMessage(cmd.Use)
	logmsg.ImportArgs(cmd, args)
	err := logmsg.Print()
	if err != nil {
		log.Error(err)

	}
}
