package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	tls "github.com/jvzantvoort/gextend-bash/printing"
)

func help(a ...string) {
	fmt.Printf("%s [%s] message\n", os.Args[0], strings.Join(a, "|"))
	os.Exit(0)
}

// colorNames returns the known color/status names in sorted order.
func colorNames(cr *tls.Cprint) []string {
	keys := make([]string, 0, len(cr.Colors))
	for keyn := range cr.Colors {
		keys = append(keys, keyn)
	}
	slices.Sort(keys)
	return keys
}

func main() {
	cr := tls.NewCprint()
	args := os.Args[1:]

	if len(args) < 2 {
		help(colorNames(cr)...)
	}

	cr.Print(args...)

}
