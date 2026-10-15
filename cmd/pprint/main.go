package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	tls "github.com/jvzantvoort/gextend-bash/printing"
	"github.com/jvzantvoort/gextend-bash/version"
)

// usage returns the usage line for the given names.
func usage(name string, names []string) string {
	return fmt.Sprintf("%s [%s] message", name, strings.Join(names, "|"))
}

func help(a ...string) {
	fmt.Println(usage(os.Args[0], a))
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
	version.HandleArgs(os.Args)

	cr := tls.NewCprint()
	args := os.Args[1:]

	if len(args) < 2 {
		help(colorNames(cr)...)
	}

	cr.Print(args...)

}
