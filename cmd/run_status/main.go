// run_status runs a command and prints a status line for it, replacing the
// `command; test_result $? "message"` pattern.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/jvzantvoort/gextend-bash/utils"
	"github.com/jvzantvoort/gextend-bash/version"
)

// runner executes a command; replaceable for tests.
type runner func(name string, args []string, stdout, stderr io.Writer) error

func execRunner(name string, args []string, stdout, stderr io.Writer) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

// statusLine returns the status line for a command that exited with code.
func statusLine(code int, message string) string {
	switch code {
	case 0:
		return utils.MakeStatus("SUCCESS", "%s", message)
	case 127:
		return utils.MakeStatus("FAILURE", "%s", message+" (not found)")
	default:
		return utils.MakeStatus("FAILURE", "%s", message)
	}
}

// run executes command, writes its output (all of it, or with quiet only
// when it fails) and the status line to stdout/stderr, and returns the
// command's exit code.
func run(run runner, command []string, message string, quiet bool, stdout, stderr io.Writer) int {
	if message == "" {
		message = strings.Join(command, " ")
	}

	var buf bytes.Buffer
	cmdOut, cmdErr := stdout, stderr
	if quiet {
		cmdOut, cmdErr = &buf, &buf
	}

	err := run(command[0], command[1:], cmdOut, cmdErr)
	code := utils.ExitCode(err)

	if quiet && code != 0 {
		_, _ = buf.WriteTo(stderr)
	}
	_, _ = fmt.Fprintln(stdout, statusLine(code, message))
	return code
}

func main() {
	var message string
	var quiet, showVersion bool
	flag.StringVar(&message, "m", "", "message to print (default: the command line)")
	flag.BoolVar(&quiet, "q", false, "hide the command's output unless it fails")
	flag.BoolVar(&showVersion, "version", false, "print the version and exit")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-m message] [-q] [--] command [args...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Runs command, prints a SUCCESS/FAILURE status line and exits with its exit code.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if showVersion {
		version.Print(os.Stdout, "run_status")
		return
	}
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	os.Exit(run(execRunner, flag.Args(), message, quiet, os.Stdout, os.Stderr))
}
