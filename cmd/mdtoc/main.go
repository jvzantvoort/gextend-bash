package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Heading represents a markdown heading
type Heading struct {
	Level int
	Text  string
	Anchor string
}

// generateAnchor creates an anchor link based on the heading text
func generateAnchor(text string) string {
	anchor := strings.ToLower(text)
	anchor = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(anchor, "") // Remove special characters
	anchor = strings.ReplaceAll(anchor, " ", "-")                             // Replace spaces with dashes
	return anchor
}

// parseMarkdown extracts headings from a file
func parseMarkdown(filename string, maxDepth int) ([]Heading, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var headings []Heading
	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`^(#{1,6})\s+(.+)$`) // Match markdown headings

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			level := len(matches[1]) // Count the number of '#' to determine level
			if level <= maxDepth {
				text := matches[2]
				anchor := generateAnchor(text)
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
		sb.WriteString(fmt.Sprintf("%s- [%s](#%s)\n", indent, h.Text, h.Anchor))
	}
	return sb.String()
}

func main() {
	// Define command-line flags
	filename := flag.String("file", "README.md", "Markdown file to generate TOC for")
	maxDepth := flag.Int("depth", 3, "Maximum heading depth to include in TOC")
	flag.Parse()

	// Parse markdown file
	headings, err := parseMarkdown(*filename, *maxDepth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading markdown file: %v\n", err)
		os.Exit(1)
	}

	// Generate and print TOC
	toc := generateTOC(headings)
	fmt.Println(toc)
}

