// retry runs a command until it succeeds or the attempts run out.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/jvzantvoort/gextend-bash/utils"
	"github.com/jvzantvoort/gextend-bash/version"
)

// Options controls the retry loop.
type Options struct {
	Attempts int           // total number of runs, at least 1
	Delay    time.Duration // pause before the second run
	Backoff  bool          // double the pause after every failure
	MaxDelay time.Duration // upper bound for the pause when backing off
}

// nextDelay returns the pause after a failure that followed a pause of d.
func (o Options) nextDelay(d time.Duration) time.Duration {
	if !o.Backoff {
		return d
	}
	d *= 2
	if o.MaxDelay > 0 && d > o.MaxDelay {
		d = o.MaxDelay
	}
	return d
}

// retry calls attempt until it returns 0, the attempts run out, or the
// command cannot be run at all (126/127, retrying would not help). It
// returns the last exit code.
func retry(opts Options, attempt func() int, sleep func(time.Duration), log io.Writer) int {
	delay := opts.Delay
	code := 0
	for i := 1; i <= opts.Attempts; i++ {
		code = attempt()
		if code == 0 {
			return 0
		}
		if code == 126 || code == 127 {
			_, _ = fmt.Fprintf(log, "retry: command could not be run (exit %d), giving up\n", code)
			return code
		}
		if i == opts.Attempts {
			break
		}
		_, _ = fmt.Fprintf(log, "retry: attempt %d/%d failed (exit %d), retrying in %s\n", i, opts.Attempts, code, delay)
		sleep(delay)
		delay = opts.nextDelay(delay)
	}
	_, _ = fmt.Fprintf(log, "retry: giving up after %d attempts (exit %d)\n", opts.Attempts, code)
	return code
}

func main() {
	opts := Options{}
	var showVersion bool
	flag.IntVar(&opts.Attempts, "n", 3, "number of attempts")
	flag.DurationVar(&opts.Delay, "d", time.Second, "delay between attempts")
	flag.BoolVar(&opts.Backoff, "b", false, "double the delay after every failed attempt")
	flag.DurationVar(&opts.MaxDelay, "max-delay", time.Minute, "maximum delay when backing off")
	flag.BoolVar(&showVersion, "version", false, "print the version and exit")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-n attempts] [-d delay] [-b] [-max-delay d] [--] command [args...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Runs command until it succeeds; exits with its last exit code.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if showVersion {
		version.Print(os.Stdout, "retry")
		return
	}
	if flag.NArg() == 0 || opts.Attempts < 1 || opts.Delay < 0 {
		flag.Usage()
		os.Exit(2)
	}

	args := flag.Args()
	attempt := func() int {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		return utils.ExitCode(cmd.Run())
	}

	os.Exit(retry(opts, attempt, time.Sleep, os.Stderr))
}
