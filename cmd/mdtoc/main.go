package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jvzantvoort/gextend-bash/version"
)

// Heading represents a markdown heading
type Heading struct {
	Level  int
	Text   string
	Anchor string
}

var (
	headingPattern = regexp.MustCompile(`^(#{1,6})\s+(.+?)(?:\s+#+)?\s*$`) // ATX heading, optional closing #s
	fencePattern   = regexp.MustCompile("^ {0,3}(`{3,}|~{3,})")            // code fence open/close
	anchorStrip    = regexp.MustCompile(`[^\p{L}\p{M}\p{N}\p{Pc} -]`)      // what GitHub drops from anchors
)

// generateAnchor creates an anchor link based on the heading text, the way
// GitHub does: lowercase, punctuation removed (except '-' and '_'), and
// spaces replaced by dashes.
func generateAnchor(text string) string {
	anchor := strings.ToLower(text)
	anchor = anchorStrip.ReplaceAllString(anchor, "")
	anchor = strings.ReplaceAll(anchor, " ", "-")
	return anchor
}

// anchorSet hands out unique anchors: GitHub suffixes the second and later
// occurrences of an anchor with -1, -2, ...
type anchorSet map[string]int

func (a anchorSet) unique(anchor string) string {
	count := a[anchor]
	a[anchor] = count + 1
	if count == 0 {
		return anchor
	}
	return fmt.Sprintf("%s-%d", anchor, count)
}

// parseMarkdown extracts headings from a file, ignoring anything inside
// fenced code blocks.
func parseMarkdown(filename string, maxDepth int) ([]Heading, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var headings []Heading
	anchors := anchorSet{}
	fence := "" // the open fence, empty when outside a code block
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if m := fencePattern.FindStringSubmatch(line); m != nil {
			switch {
			case fence == "":
				fence = m[1]
			case m[1][0] == fence[0] && len(m[1]) >= len(fence) && strings.TrimSpace(line) == m[1]:
				fence = ""
			}
			continue
		}
		if fence != "" {
			continue
		}

		matches := headingPattern.FindStringSubmatch(line)
		if matches != nil {
			level := len(matches[1]) // Count the number of '#' to determine level
			text := matches[2]
			// register every heading so duplicate numbering matches GitHub,
			// even for headings deeper than maxDepth
			anchor := anchors.unique(generateAnchor(text))
			if level <= maxDepth {
				headings = append(headings, Heading{Level: level, Text: text, Anchor: anchor})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return headings, nil
}

// generateTOC creates a TOC string from headings
func generateTOC(headings []Heading) string {
	var sb strings.Builder
	for _, h := range headings {
		indent := strings.Repeat("  ", h.Level-1) // Indent based on level
		fmt.Fprintf(&sb, "%s- [%s](#%s)\n", indent, h.Text, h.Anchor)
	}
	return sb.String()
}

const (
	tocStart = "<!-- toc -->"
	tocStop  = "<!-- tocstop -->"
)

// replaceTOC returns content with the text between the tocStart and tocStop
// marker lines replaced by toc.
func replaceTOC(content, toc string) (string, error) {
	start := strings.Index(content, tocStart)
	if start < 0 {
		return "", fmt.Errorf("marker %s not found", tocStart)
	}
	bodyStart := start + len(tocStart)
	stop := strings.Index(content[bodyStart:], tocStop)
	if stop < 0 {
		return "", fmt.Errorf("marker %s not found after %s", tocStop, tocStart)
	}
	stop += bodyStart

	return content[:bodyStart] + "\n\n" + strings.TrimRight(toc, "\n") + "\n\n" + content[stop:], nil
}

// writeTOC replaces the TOC between the markers in filename.
func writeTOC(filename, toc string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	updated, err := replaceTOC(string(data), toc)
	if err != nil {
		return fmt.Errorf("%s: %w", filename, err)
	}
	if updated == string(data) {
		return nil
	}
	return os.WriteFile(filename, []byte(updated), info.Mode().Perm())
}

func main() {
	// Define command-line flags
	filename := flag.String("file", "README.md", "Markdown file to generate TOC for")
	maxDepth := flag.Int("depth", 3, "Maximum heading depth to include in TOC")
	inPlace := flag.Bool("w", false, "write the TOC into the file between "+tocStart+" and "+tocStop+" instead of printing it")
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		version.Print(os.Stdout, "mdtoc")
		return
	}

	// Parse markdown file
	headings, err := parseMarkdown(*filename, *maxDepth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading markdown file: %v\n", err)
		os.Exit(1)
	}

	// Generate and print TOC
	toc := generateTOC(headings)
	if *inPlace {
		if err := writeTOC(*filename, toc); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating markdown file: %v\n", err)
			os.Exit(1)
		}
		return
	}
	fmt.Println(toc)
}
