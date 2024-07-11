/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// critCmd represents the crit command
var critCmd = &cobra.Command{
	Use:   "crit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("crit called")
	},
}

func init() {
	rootCmd.AddCommand(critCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// critCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// critCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
