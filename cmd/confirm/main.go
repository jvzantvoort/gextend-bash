// confirm asks a yes/no question; the answer is the exit code:
//
//	if confirm -d no "Delete the backups?"; then ... fi
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jvzantvoort/gextend-bash/version"
)

// Default is the answer used for an empty reply, end of input or a timeout.
type Default int

const (
	NoDefault Default = iota
	DefaultYes
	DefaultNo
)

func parseDefault(s string) (Default, error) {
	switch strings.ToLower(s) {
	case "":
		return NoDefault, nil
	case "y", "yes":
		return DefaultYes, nil
	case "n", "no":
		return DefaultNo, nil
	}
	return NoDefault, fmt.Errorf("default %q is not yes or no", s)
}

func (d Default) hint() string {
	switch d {
	case DefaultYes:
		return "[Y/n]"
	case DefaultNo:
		return "[y/N]"
	}
	return "[y/n]"
}

// answer is what an unanswered question resolves to: the default, or no.
func (d Default) answer() bool {
	return d == DefaultYes
}

// ask writes question to out and reads replies from in until it gets a yes
// or no. An empty reply takes the default; end of input or the timeout (0
// waits forever) takes the default, or no when there is none.
func ask(in io.Reader, out io.Writer, question string, def Default, timeout time.Duration) bool {
	lines := make(chan string)
	go func() {
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	var expired <-chan time.Time
	if timeout > 0 {
		expired = time.After(timeout)
	}

	for {
		_, _ = fmt.Fprintf(out, "%s %s ", question, def.hint())
		select {
		case line, ok := <-lines:
			if !ok {
				_, _ = fmt.Fprintln(out)
				return def.answer()
			}
			switch strings.ToLower(strings.TrimSpace(line)) {
			case "y", "yes":
				return true
			case "n", "no":
				return false
			case "":
				if def != NoDefault {
					return def.answer()
				}
			}
			_, _ = fmt.Fprintln(out, "Please answer yes or no.")
		case <-expired:
			_, _ = fmt.Fprintln(out)
			return def.answer()
		}
	}
}

func main() {
	var defaultStr string
	var timeout time.Duration
	var showVersion bool
	flag.StringVar(&defaultStr, "d", "", "default answer (yes or no) for an empty reply, end of input or timeout")
	flag.DurationVar(&timeout, "t", 0, "give up after this long, e.g. 30s (0 waits forever)")
	flag.BoolVar(&showVersion, "version", false, "print the version and exit")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-d yes|no] [-t timeout] question...\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Exits 0 for yes and 1 for no. Without a default, end of input or a timeout counts as no.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if showVersion {
		version.Print(os.Stdout, "confirm")
		return
	}

	def, err := parseDefault(defaultStr)
	if err != nil || flag.NArg() == 0 {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		flag.Usage()
		os.Exit(2)
	}

	if ask(os.Stdin, os.Stderr, strings.Join(flag.Args(), " "), def, timeout) {
		os.Exit(0)
	}
	os.Exit(1)
}
