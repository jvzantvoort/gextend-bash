package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jvzantvoort/gextend-bash/utils"
	"github.com/jvzantvoort/gextend-bash/version"
)

func Die(instr ...string) {
	if len(instr) != 0 {
		fmt.Fprintln(os.Stderr, strings.Join(instr, " "))
	}
	fmt.Fprintf(os.Stderr, "USAGE:\n\n")
	fmt.Fprintf(os.Stderr, "\t%s <exitcode> <message>\n\n", os.Args[0])
	os.Exit(1)
}

// parseArgs splits the command line into the exit code and the message.
func parseArgs(args []string) (int, string, error) {
	if len(args) < 2 {
		return 0, "", errors.New("not enough arguments")
	}
	exitcode, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, "", fmt.Errorf("error: %s is not a number", args[0])
	}
	return exitcode, strings.Join(args[1:], " "), nil
}

// render returns the status line for exitcode: SUCCESS for 0, FAILURE
// otherwise.
func render(exitcode int, message string) string {
	if exitcode == 0 {
		return utils.MakeStatus("SUCCESS", "%s", message)
	}
	return utils.MakeStatus("FAILURE", "%s", message)
}

func main() {
	version.HandleArgs(os.Args)

	exitcode, message, err := parseArgs(os.Args[1:])
	if err != nil {
		Die(err.Error())
	}

	fmt.Println(render(exitcode, message))
	os.Exit(exitcode)
}
