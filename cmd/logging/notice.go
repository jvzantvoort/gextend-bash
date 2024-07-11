/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"github.com/jvzantvoort/gextend-bash/messages"
	"github.com/spf13/cobra"
)

// NoticeCmd represents the notice command
var NoticeCmd = &cobra.Command{
	Use:   "notice",
	Short: "Log a notice command",
	Long:  messages.GetLong("logging"),
	Run:   handleLogCmd,
}

func init() {
	rootCmd.AddCommand(NoticeCmd)
	NoticeCmd.Flags().StringP("tag", "t", "", "mark every line with this tag")
	NoticeCmd.Flags().StringP("file", "f", "output.log", "log the contents of this file")
	NoticeCmd.Flags().StringP("priority", "p", "", "mark given message with this priority")
	NoticeCmd.Flags().BoolP("skip-empty", "e", false, "do not log empty lines when processing files")
	NoticeCmd.Flags().BoolP("stderr", "s", false, "output message to standard error as well")
}
