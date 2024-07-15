package main

import (
	"fmt"

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
	filepath := "/home/jvzantvoort/Logs/2024/July/common.log"

	logmsgs := logging.NewLogMessages(filepath)
	for _, msg := range logmsgs.Messages {
		tag := ""
		if len(msg.Tag) != 0 {
			tag = fmt.Sprintf("<%s> ", msg.Tag)
		}
		fmt.Printf("%s %-8s %s%s\n",
			msg.Time.Format("2006-01-02 15:04:05"),
			msg.Priority,
			tag,
			msg.Message,
		)
	}
	/*
		for _, msg := range logmsgs.messages {

			fmt.Printf("%#v\n", msg)

		}
	*/
}
