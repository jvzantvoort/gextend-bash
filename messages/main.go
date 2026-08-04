package messages

import (
	"embed"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Content holds the embedded long-description text for each print_status subcommand.
//
//go:embed long/*
var Content embed.FS

func GetLong(name string) string {
	filename := fmt.Sprintf("long/%s", name)
	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		log.Error(err)
		msgstr = []byte("undefined")
	}
	return string(msgstr)
}
