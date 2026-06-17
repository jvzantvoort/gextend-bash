/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jvzantvoort/gextend-bash/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// monitorState holds the display buffer for the monitor view.
// rawLines stores raw JSON strings that have already passed the level filter,
// newest entry at index 0.
type monitorState struct {
	mu       sync.Mutex
	rawLines []string
	maxLines int
	width    int
	minLevel int
}

// prepend adds a new raw JSON line at the top of the buffer, evicting the
// oldest entry when the buffer is at capacity.
func (m *monitorState) prepend(raw string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rawLines = append([]string{raw}, m.rawLines...)
	if len(m.rawLines) > m.maxLines {
		m.rawLines = m.rawLines[:m.maxLines]
	}
}

// resize updates terminal dimensions and trims the buffer to fit.
func (m *monitorState) resize() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.width = consoleWidth()
	m.maxLines = consoleHeight()
	if len(m.rawLines) > m.maxLines {
		m.rawLines = m.rawLines[:m.maxLines]
	}
}

// redraw repaints the entire alternate screen in place without clearing first,
// which avoids the flicker of a full erase. Each line is overwritten and then
// cleared to end-of-line; any leftover screen rows below the buffer are erased.
func (m *monitorState) redraw() {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Print("\033[H") // move cursor to top-left
	for _, raw := range m.rawLines {
		line := formatJSONLine(raw, m.minLevel, m.width)
		if line == "" {
			line = raw
		}
		fmt.Printf("%s\033[K\n", line)
	}
	fmt.Print("\033[J") // erase from cursor to end of screen
}

// enterAlternateScreen switches to the xterm alternate screen buffer and hides
// the cursor so the display looks clean.
func enterAlternateScreen() {
	fmt.Print("\033[?1049h\033[?25l")
}

// exitAlternateScreen restores the normal screen buffer and cursor.
func exitAlternateScreen() {
	fmt.Print("\033[?25h\033[?1049l")
}

// runMonitor is the main loop: reads initial history, then tails the file
// printing newest messages at the top of the alternate screen.
func runMonitor(filename string, initialN int, minLevel int) {
	width := consoleWidth()
	height := consoleHeight()

	state := &monitorState{
		maxLines: height,
		width:    width,
		minLevel: minLevel,
	}

	// Load initial history oldest-first so prepending leaves newest at top.
	f, err := os.Open(filename)
	if err != nil {
		log.Errorf("Cannot open %s: %s", filename, err)
		os.Exit(1)
	}

	var allRaw []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		raw := strings.TrimSpace(scanner.Text())
		if formatJSONLine(raw, minLevel, width) != "" {
			allRaw = append(allRaw, raw)
		}
	}
	if err := scanner.Err(); err != nil {
		f.Close()
		log.Errorf("Error reading %s: %s", filename, err)
		os.Exit(1)
	}

	if len(allRaw) > initialN {
		allRaw = allRaw[len(allRaw)-initialN:]
	}
	for _, raw := range allRaw {
		state.prepend(raw)
	}

	// Seek to end before entering follow mode.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		log.Errorf("Seek failed on %s: %s", filename, err)
		os.Exit(1)
	}

	enterAlternateScreen()

	// Restore terminal on SIGINT / SIGTERM.
	sigExit := make(chan os.Signal, 1)
	signal.Notify(sigExit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigExit
		exitAlternateScreen()
		os.Exit(0)
	}()

	// Repaint on terminal resize (SIGWINCH).
	sigWinch := make(chan os.Signal, 1)
	signal.Notify(sigWinch, syscall.SIGWINCH)
	go func() {
		for range sigWinch {
			state.resize()
			state.redraw()
		}
	}()

	state.redraw()

	// Follow the file, prepending each new line that passes the filter.
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(200 * time.Millisecond)
				continue
			}
			exitAlternateScreen()
			log.Errorf("Error reading %s: %s", filename, err)
			os.Exit(1)
		}
		raw := strings.TrimRight(line, "\n")
		if formatJSONLine(raw, minLevel, width) != "" {
			state.prepend(raw)
			state.redraw()
		}
	}
}

var monitorCmd = &cobra.Command{
	Use:   "monitor [file]",
	Short: "Live log monitor — newest messages appear at the top",
	Long: `Watch a JSON-structured log file in real time. New entries are inserted at the
top of the screen; existing lines shift down. The display fills the terminal
height and is refreshed in place without flickering.

Press Ctrl+C to exit and return to the normal terminal.

Severity levels (highest to lowest): EMERG, ALERT, CRIT, ERR, WARNING, NOTICE, INFO, DEBUG`,
	Args: cobra.MaximumNArgs(1),
	Run:  runMonitorCmd,
}

var (
	monitorLines int
	monitorLevel string
	monitorFile  string
)

func runMonitorCmd(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	filename := monitorFile
	if len(args) > 0 {
		filename = args[0]
	}
	if filename == "" {
		cfg := config.NewConfigLogging()
		var err error
		filename, err = cfg.LogfilePath()
		if err != nil {
			log.Errorf("Cannot determine log file path: %s", err)
			os.Exit(1)
		}
	}

	log.Debugf("log file: %s", filename)
	runMonitor(filename, monitorLines, levelNum(monitorLevel))
}

func init() {
	rootCmd.AddCommand(monitorCmd)
	monitorCmd.Flags().IntVarP(&monitorLines, "lines", "n", 50, "number of historical lines to load on startup")
	monitorCmd.Flags().StringVarP(&monitorLevel, "level", "l", "DEBUG", "minimum severity level to display (EMERG|ALERT|CRIT|ERR|WARNING|NOTICE|INFO|DEBUG)")
	monitorCmd.Flags().StringVar(&monitorFile, "file", "", "log file to read (default: from config)")
}
