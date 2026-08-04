package main

import (
	"os"
	"strings"

	"github.com/jvzantvoort/gextend-bash/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type PrettyPrint struct {
	Status  string `json:"priority"`
	Message string `json:"message"`
}

func (l PrettyPrint) Print() error {
	// make parent dirs
	msg := []byte(utils.MakeStatus(l.Status, "%s", l.Message) + "\n")

	if _, err := os.Stderr.Write(msg); err != nil {
		return err
	}
	return nil

}

func NewPrettyPrint(level string) *PrettyPrint {
	retv := &PrettyPrint{}
	retv.Status = strings.ToUpper(level)

	return retv
}

func handlePrintCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("%s: start", cmd.Use)
	defer log.Debugf("%s: end", cmd.Use)

	if len(args) == 0 {
		log.Error("No message provided")
		if err := cmd.Help(); err != nil {
			log.Error(err)
		}
		os.Exit(1)
	}
	logmsg := NewPrettyPrint(cmd.Use)
	logmsg.Message = strings.Join(args, " ")
	err := logmsg.Print()
	if err != nil {
		log.Error(err)

	}
}
