package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/jvzantvoort/gextend-bash/colors"
)

func help(a ...string) {
	fmt.Printf("%s [%s] message\n", os.Args[0], strings.Join(a, "|"))
	os.Exit(0)
}

func main() {
	cr := colors.NewCprint()
	args := os.Args[1:]

	if len(args) < 2 {
		keys := []string{}
		for keyn := range cr.Colors {
			keys = append(keys, keyn)
		}
		help(keys...)
	}

	cr.Print(args...)

}
