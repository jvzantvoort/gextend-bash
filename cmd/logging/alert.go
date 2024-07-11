/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"github.com/jvzantvoort/gextend-bash/messages"
	"github.com/spf13/cobra"
)

// AlertCmd represents the alert command
var AlertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Log a alert command",
	Long:  messages.GetLong("logging"),
	Run:   handleLogCmd,
}

func init() {
	rootCmd.AddCommand(AlertCmd)
	AlertCmd.Flags().StringP("tag", "t", "", "mark every line with this tag")
	AlertCmd.Flags().StringP("file", "f", "output.log", "log the contents of this file")
	AlertCmd.Flags().StringP("priority", "p", "", "mark given message with this priority")
	AlertCmd.Flags().BoolP("skip-empty", "e", false, "do not log empty lines when processing files")
	AlertCmd.Flags().BoolP("stderr", "s", false, "output message to standard error as well")
}
