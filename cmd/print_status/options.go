package main

import (
	"fmt"

	"github.com/jvzantvoort/gextend-bash/messages"
	"github.com/spf13/cobra"
)

var translations = map[string]string{
	"success":   "success",
	"ok":        "success",
	"oke":       "success",
	"failure":   "failure",
	"fail":      "failure",
	"nok":       "failure",
	"notice":    "notice",
	"info":      "notice",
	"warning":   "warning",
	"warn":      "warning",
	"unknown":   "unknown",
	"undefined": "unknown",
}

// Options lists the registered subcommand names (the keys of translations).
var Options []string

func init() {

	for keyn, keyv := range translations {
		short := fmt.Sprintf("print a %s message", keyv)

		rootCmd.AddCommand(
			&cobra.Command{
				Use:   keyn,
				Short: short,
				Long:  messages.GetLong(keyv),
				Run:   handlePrintCmd,
			})

		Options = append(Options, keyn)
	}
}
