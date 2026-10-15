package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/jvzantvoort/gextend-bash/version"
)

func sec2Human(insec int) string {
	var retv string
	hours := 0
	minutes := 0
	seconds := insec

	if seconds < 1 {
		return "None"
	}

	if seconds >= 3600 {
		hours = int(seconds / 3600)
		seconds = seconds % 3600
	}

	if seconds >= 60 {
		minutes = int(seconds / 60)
		seconds = seconds % 60
	}

	if hours > 0 {
		retv = fmt.Sprintf("%02d:", hours)
	}

	if minutes > 0 {
		retv = fmt.Sprintf("%s%02d:", retv, minutes)
	} else if (minutes == 0) && (hours > 0) {
		retv = fmt.Sprintf("%s00:", retv)
	}

	if len(retv) > 0 {
		retv = fmt.Sprintf("%s%02d", retv, seconds)
	} else {
		retv = strconv.Itoa(seconds) + " seconds"
	}

	return retv
}

// calculate the variouse changes
func getTimes(epoch int, duration int) (int, string) {
	cur_ts := time.Now()
	cur_sec := int(cur_ts.Unix())
	delta_sec := cur_sec - epoch
	remaining := duration - delta_sec
	percent := int(remaining * 100 / duration)
	return percent, sec2Human(remaining)
}

// barColor picks the gauge color from the percentage of time remaining:
// green while there is plenty left, then yellow, then red near the end.
func barColor(percent int) ui.Color {
	switch {
	case percent > 65:
		return ui.ColorGreen
	case percent > 33:
		return ui.ColorYellow
	default:
		return ui.ColorRed
	}
}

// parseDuration accepts a number of seconds ("90") or a Go duration
// ("90s", "5m", "1h30m") and returns whole seconds, at least 1.
func parseDuration(arg string) (int, error) {
	if seconds, err := strconv.Atoi(arg); err == nil {
		if seconds < 1 {
			return 0, fmt.Errorf("duration %q must be at least one second", arg)
		}
		return seconds, nil
	}
	d, err := time.ParseDuration(arg)
	if err != nil {
		return 0, fmt.Errorf("duration %q is not a number of seconds or a duration like 5m or 1h30m", arg)
	}
	if d < time.Second {
		return 0, fmt.Errorf("duration %q must be at least one second", arg)
	}
	return int(d.Round(time.Second) / time.Second), nil
}

// run shows the countdown and reports whether it ran to completion (false
// when the user quit early).
func run(seconds int) (bool, error) {
	if err := ui.Init(); err != nil {
		return false, fmt.Errorf("failed to initialize termui: %w", err)
	}
	defer ui.Close()

	now := time.Now()
	start_secs := int(now.Unix())

	screenwidth, screenheight := ui.TerminalDimensions()
	padding := 10
	screenheight = (screenheight - (screenheight % 10)) / 2

	g0 := widgets.NewGauge()
	g0.Title = "Countdown"
	g0.Percent = 100
	g0.SetRect(padding, screenheight, screenwidth-padding, screenheight+7)

	g0.BarColor = ui.ColorGreen
	g0.BorderStyle.Fg = ui.ColorWhite
	g0.TitleStyle.Fg = ui.ColorCyan

	draw := func() {
		percent, tsstr := getTimes(start_secs, seconds)
		g0.BarColor = barColor(percent)
		g0.Percent = percent
		g0.Title = fmt.Sprintf("Countdown %s", tsstr)
		ui.Render(g0)
	}

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C

	for {
		percent, _ := getTimes(start_secs, seconds)
		if percent < 1 {
			return true, nil
		}
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return false, nil
			}
		case <-ticker:
			draw()
		}
	}
}

func main() {
	var noBell, showVersion bool
	flag.BoolVar(&noBell, "no-bell", false, "do not ring the terminal bell when the countdown ends")
	flag.BoolVar(&showVersion, "version", false, "print the version and exit")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--no-bell] <seconds|duration>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Durations use Go syntax, e.g. 90, 90s, 5m, 1h30m. Press q to quit.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if showVersion {
		version.Print(os.Stdout, "countdown")
		return
	}
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	seconds, err := parseDuration(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	finished, err := run(seconds)
	if err != nil {
		log.Fatal(err)
	}
	if finished && !noBell {
		fmt.Print("\a")
	}
}
