/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"
	"github.com/jvzantvoort/gextend-bash/logging"
	"github.com/jvzantvoort/gextend-bash/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NoticeCmd represents the notice command
var NoticeCmd = &cobra.Command{
	Use:   "notice",
	Short: "Log a notice command",
	Long:  messages.GetLong("logging"),
	Run:   handleLogNoticeCmd,
}

func handleLogNoticeCmd(cmd *cobra.Command, args []string) {
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
	logmsg.Print()
}

func init() {
	rootCmd.AddCommand(NoticeCmd)
	NoticeCmd.Flags().StringP("tag", "t", "", "mark every line with this tag")
	NoticeCmd.Flags().StringP("file", "f", "output.log", "log the contents of this file")
	NoticeCmd.Flags().StringP("priority", "p", "", "mark given message with this priority")
	NoticeCmd.Flags().BoolP("skip-empty", "e", false, "do not log empty lines when processing files")
	NoticeCmd.Flags().BoolP("stderr", "s", false, "output message to standard error as well")
}
