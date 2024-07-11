/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"github.com/jvzantvoort/gextend-bash/messages"
	"github.com/spf13/cobra"
)

// ErrCmd represents the err command
var ErrCmd = &cobra.Command{
	Use:   "err",
	Short: "Log a err command",
	Long:  messages.GetLong("logging"),
	Run:   handleLogCmd,
}

func init() {
	rootCmd.AddCommand(ErrCmd)
	ErrCmd.Flags().StringP("tag", "t", "", "mark every line with this tag")
	ErrCmd.Flags().StringP("file", "f", "output.log", "log the contents of this file")
	ErrCmd.Flags().StringP("priority", "p", "", "mark given message with this priority")
	ErrCmd.Flags().BoolP("skip-empty", "e", false, "do not log empty lines when processing files")
	ErrCmd.Flags().BoolP("stderr", "s", false, "output message to standard error as well")
}
