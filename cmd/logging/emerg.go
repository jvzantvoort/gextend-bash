/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"github.com/jvzantvoort/gextend-bash/messages"
	"github.com/spf13/cobra"
)

// EmergCmd represents the emerg command
var EmergCmd = &cobra.Command{
	Use:   "emerg",
	Short: "Log a emerg command",
	Long:  messages.GetLong("logging"),
	Run:   handleLogCmd,
}

func init() {
	rootCmd.AddCommand(EmergCmd)
	EmergCmd.Flags().StringP("tag", "t", "", "mark every line with this tag")
	EmergCmd.Flags().StringP("file", "f", "output.log", "log the contents of this file")
	EmergCmd.Flags().StringP("priority", "p", "", "mark given message with this priority")
	EmergCmd.Flags().BoolP("skip-empty", "e", false, "do not log empty lines when processing files")
	EmergCmd.Flags().BoolP("stderr", "s", false, "output message to standard error as well")
}
