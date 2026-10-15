package utils

import (
	"io"
	"os"
	"strings"
)

// ArgsOrStdin returns a reader over args joined by spaces, or os.Stdin when
// there are no args.
func ArgsOrStdin(args []string) io.Reader {
	if len(args) > 0 {
		return strings.NewReader(strings.Join(args, " "))
	}
	return os.Stdin
}
