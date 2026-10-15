// bash autocomplete helper
package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/jvzantvoort/gextend-bash/version"
)

//go:embed message/*
var Content embed.FS

func AbortMe() {
	msgstr, err := Content.ReadFile("message/abort")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		msgstr = []byte("undefined")
	}
	fmt.Fprintf(os.Stderr, "%s", string(msgstr))
	os.Exit(1)
}

func main() {
	if len(os.Args) == 1 {
		AbortMe()
	}

	version.HandleArgs(os.Args)

	if os.Args[1] == "ssh" {
		knownHosts := false
		for _, arg := range os.Args[2:] {
			switch arg {
			case "--known-hosts", "-known-hosts":
				knownHosts = true
			default:
				AbortMe()
			}
		}
		err := SecureShellHelper(knownHosts)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		AbortMe()
	}
}
