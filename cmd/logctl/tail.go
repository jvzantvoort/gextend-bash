/*
Copyright © 2024 John van Zantvoort <john@vanzantvoort.org>
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jvzantvoort/gextend-bash/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// severityOrder maps log level names to syslog numeric priority (0=highest, 7=lowest).
var severityOrder = map[string]int{
	"EMERG":   0,
	"ALERT":   1,
	"CRIT":    2,
	"ERR":     3,
	"ERROR":   3,
	"WARNING": 4,
	"WARN":    4,
	"NOTICE":  5,
	"INFO":    6,
	"DEBUG":   7,
}

// logMessage holds the fields we read from each JSON log entry.
type logMessage struct {
	Tag      string    `json:"tag"`
	Priority string    `json:"priority"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
}

func severityColor(priority string) func(string, ...any) string {
	switch strings.ToUpper(priority) {
	case "EMERG":
		return color.New(color.FgWhite, color.BgRed, color.Bold).Sprintf
	case "ALERT":
		return color.New(color.FgRed, color.Bold).Sprintf
	case "CRIT":
		return color.New(color.FgRed).Sprintf
	case "ERR", "ERROR":
		return color.New(color.FgHiRed).Sprintf
	case "WARNING", "WARN":
		return color.New(color.FgYellow).Sprintf
	case "NOTICE":
		return color.New(color.FgCyan).Sprintf
	case "INFO":
		return color.New(color.FgGreen).Sprintf
	case "DEBUG":
		return color.New(color.FgHiBlack).Sprintf
	default:
		return color.New(color.FgWhite).Sprintf
	}
}

func levelNum(level string) int {
	n, ok := severityOrder[strings.ToUpper(level)]
	if !ok {
		return severityOrder["DEBUG"]
	}
	return n
}

// formatLogMessage builds a colored, width-truncated display string from a parsed message.
// Truncation is applied before colorizing so escape sequences don't count toward width.
func formatLogMessage(msg logMessage, width int) string {
	prio := strings.ToUpper(msg.Priority)
	ts := msg.Time.Format("2006-01-02 15:04:05")

	var parts []string
	parts = append(parts, ts)
	parts = append(parts, fmt.Sprintf("%-7s", prio))
	if msg.Tag != "" {
		parts = append(parts, fmt.Sprintf("[%s]", msg.Tag))
	}
	parts = append(parts, msg.Message)

	line := strings.Join(parts, " ")
	if width > 0 && len(line) > width {
		line = line[:width]
	}

	return severityColor(prio)("%s", line)
}

// formatJSONLine parses raw JSON, applies the level filter, and returns a formatted
// colored string. Returns "" when the line is filtered out, empty, or unparseable as JSON
// (unparseable lines are returned truncated but uncolored).
func formatJSONLine(raw string, minLevel int, width int) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	var msg logMessage
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		if width > 0 && len(raw) > width {
			return raw[:width]
		}
		return raw
	}
	if levelNum(msg.Priority) > minLevel {
		return ""
	}
	return formatLogMessage(msg, width)
}

func printJSON(raw string, minLevel int, width int) {
	if s := formatJSONLine(raw, minLevel, width); s != "" {
		fmt.Println(s)
	}
}

func printLastN(filename string, n int, minLevel int, width int) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	for _, raw := range lines {
		printJSON(raw, minLevel, width)
	}
	return nil
}

func followFile(filename string, minLevel int, width int) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(200 * time.Millisecond)
				continue
			}
			return err
		}
		printJSON(strings.TrimRight(line, "\n"), minLevel, width)
	}
}

var tailCmd = &cobra.Command{
	Use:   "tail [file]",
	Short: "Display last lines of a structured log file",
	Long: `Display the last N lines of a JSON-structured log file with color-coded severity.

Severity levels (highest to lowest): EMERG, ALERT, CRIT, ERR, WARNING, NOTICE, INFO, DEBUG

When --level is set, only messages at that level or higher severity are shown.
Use --follow to watch the file for new entries in real time.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runTail,
}

var (
	tailLines  int
	tailFollow bool
	tailLevel  string
	tailFile   string
)

func runTail(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	filename := tailFile
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
	width := consoleWidth()
	minLevel := levelNum(tailLevel)

	if err := printLastN(filename, tailLines, minLevel, width); err != nil {
		log.Errorf("Error reading %s: %s", filename, err)
		os.Exit(1)
	}

	if tailFollow {
		if err := followFile(filename, minLevel, width); err != nil {
			log.Errorf("Error following %s: %s", filename, err)
			os.Exit(1)
		}
	}
}

func init() {
	rootCmd.AddCommand(tailCmd)
	tailCmd.Flags().IntVarP(&tailLines, "lines", "n", 10, "number of lines to show")
	tailCmd.Flags().BoolVarP(&tailFollow, "follow", "f", false, "follow the log file for new entries")
	tailCmd.Flags().StringVarP(&tailLevel, "level", "l", "DEBUG", "minimum severity level to display (EMERG|ALERT|CRIT|ERR|WARNING|NOTICE|INFO|DEBUG)")
	tailCmd.Flags().StringVar(&tailFile, "file", "", "log file to read (default: from config)")
}
