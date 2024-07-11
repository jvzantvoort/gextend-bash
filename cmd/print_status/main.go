package main

import (
	"fmt"
	"os"
	"strings"
	tls "github.com/jvzantvoort/gextend-bash/printing"
)

func help(a ...string) {
	fmt.Printf("%s [%s] message\n", os.Args[0], strings.Join(a, "|"))
	os.Exit(0)
}

func main() {
	cr := tls.NewCprint()
	args := os.Args[1:]

	if len(args) < 2 {
		keys := []string{}
		for keyn, _ := range cr.Colors {
			keys = append(keys, keyn)
		}
		help(keys...)
	}

	cr.Print(args...)

}
