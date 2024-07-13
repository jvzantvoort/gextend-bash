/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"github.com/jvzantvoort/gextend-bash/messages"
	"github.com/spf13/cobra"
)

// WarningCmd represents the warning command
var WarningCmd = &cobra.Command{
	Use:   "warning",
	Short: "Log a warning command",
	Long:  messages.GetLong("logging"),
	Run:   handleLogCmd,
}

func init() {
	rootCmd.AddCommand(WarningCmd)
	WarningCmd.Flags().StringP("tag", "t", "", "mark every line with this tag")
	WarningCmd.Flags().StringP("file", "f", "", "log the contents of this file")
	WarningCmd.Flags().StringP("priority", "p", "", "mark given message with this priority")
	WarningCmd.Flags().BoolP("skip-empty", "e", false, "do not log empty lines when processing files")
	WarningCmd.Flags().BoolP("stderr", "s", false, "output message to standard error as well")
}
