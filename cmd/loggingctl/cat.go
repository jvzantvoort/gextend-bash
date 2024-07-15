/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"github.com/jvzantvoort/gextend-bash/messages"
	"github.com/spf13/cobra"
)

// CatCmd represents the cat command
var CatCmd = &cobra.Command{
	Use:   "cat",
	Short: "Log a cat command",
	Long:  messages.GetLong("logging"),
	Run:   handleLogCmd,
}

func init() {
	rootCmd.AddCommand(CatCmd)
	CatCmd.Flags().StringP("tag", "t", "", "mark every line with this tag")
	CatCmd.Flags().StringP("file", "f", "", "log the contents of this file")
	CatCmd.Flags().StringP("priority", "p", "", "mark given message with this priority")
	CatCmd.Flags().BoolP("skip-empty", "e", false, "do not log empty lines when processing files")
	CatCmd.Flags().BoolP("stderr", "s", false, "output message to standard error as well")
}
