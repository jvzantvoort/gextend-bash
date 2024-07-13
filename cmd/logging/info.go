/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"github.com/jvzantvoort/gextend-bash/messages"
	"github.com/spf13/cobra"
)

// InfoCmd represents the info command
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Log a info command",
	Long:  messages.GetLong("logging"),
	Run:   handleLogCmd,
}

func init() {
	rootCmd.AddCommand(InfoCmd)
	InfoCmd.Flags().StringP("tag", "t", "", "mark every line with this tag")
	InfoCmd.Flags().StringP("file", "f", "", "log the contents of this file")
	InfoCmd.Flags().StringP("priority", "p", "", "mark given message with this priority")
	InfoCmd.Flags().BoolP("skip-empty", "e", false, "do not log empty lines when processing files")
	InfoCmd.Flags().BoolP("stderr", "s", false, "output message to standard error as well")
}
