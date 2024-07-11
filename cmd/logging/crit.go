/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"github.com/jvzantvoort/gextend-bash/messages"
	"github.com/spf13/cobra"
)

// CritCmd represents the crit command
var CritCmd = &cobra.Command{
	Use:   "crit",
	Short: "Log a crit command",
	Long:  messages.GetLong("logging"),
	Run:   handleLogCmd,
}

func init() {
	rootCmd.AddCommand(CritCmd)
	CritCmd.Flags().StringP("tag", "t", "", "mark every line with this tag")
	CritCmd.Flags().StringP("file", "f", "output.log", "log the contents of this file")
	CritCmd.Flags().StringP("priority", "p", "", "mark given message with this priority")
	CritCmd.Flags().BoolP("skip-empty", "e", false, "do not log empty lines when processing files")
	CritCmd.Flags().BoolP("stderr", "s", false, "output message to standard error as well")
}
