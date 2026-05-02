// bash autocomplete helper
package main

import (
	"embed"
	"fmt"
	"os"
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

	if os.Args[1] == "ssh" {
		err := SecureShellHelper()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		AbortMe()
	}
}
