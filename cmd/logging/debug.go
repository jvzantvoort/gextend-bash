/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"github.com/jvzantvoort/gextend-bash/messages"
	"github.com/spf13/cobra"
)

// DebugCmd represents the debug command
var DebugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Log a debug command",
	Long:  messages.GetLong("logging"),
	Run:   handleLogCmd,
}

func init() {
	rootCmd.AddCommand(DebugCmd)
	DebugCmd.Flags().StringP("tag", "t", "", "mark every line with this tag")
	DebugCmd.Flags().StringP("file", "f", "", "log the contents of this file")
	DebugCmd.Flags().StringP("priority", "p", "", "mark given message with this priority")
	DebugCmd.Flags().BoolP("skip-empty", "e", false, "do not log empty lines when processing files")
	DebugCmd.Flags().BoolP("stderr", "s", false, "output message to standard error as well")
}
