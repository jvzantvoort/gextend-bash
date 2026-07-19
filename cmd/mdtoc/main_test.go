package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateAnchor(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"simple", "Hello World", "hello-world"},
		{"strips special characters", "Foo & Bar!", "foo--bar"},
		{"collapses to lowercase", "MixedCase Heading", "mixedcase-heading"},
		{"keeps numbers", "Section 1.2", "section-12"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateAnchor(tt.text)
			if got != tt.want {
				t.Errorf("generateAnchor(%q) = %q, want %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestParseMarkdown(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "doc.md")
	content := `# Title

Some text.

## Section One

More text.

### Deep Section

#### Too Deep

## Section Two
`
	if err := os.WriteFile(target, []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	headings, err := parseMarkdown(target, 3)
	if err != nil {
		t.Fatalf("parseMarkdown() error = %v", err)
	}

	want := []Heading{
		{Level: 1, Text: "Title", Anchor: "title"},
		{Level: 2, Text: "Section One", Anchor: "section-one"},
		{Level: 3, Text: "Deep Section", Anchor: "deep-section"},
		{Level: 2, Text: "Section Two", Anchor: "section-two"},
	}

	if len(headings) != len(want) {
		t.Fatalf("parseMarkdown() returned %d headings, want %d: %+v", len(headings), len(want), headings)
	}
	for i := range want {
		if headings[i] != want[i] {
			t.Errorf("headings[%d] = %+v, want %+v", i, headings[i], want[i])
		}
	}
}

func TestParseMarkdownMissingFile(t *testing.T) {
	_, err := parseMarkdown(filepath.Join(t.TempDir(), "missing.md"), 3)
	if err == nil {
		t.Fatal("expected an error for a missing file")
	}
}

func TestGenerateTOC(t *testing.T) {
	headings := []Heading{
		{Level: 1, Text: "Title", Anchor: "title"},
		{Level: 2, Text: "Sub", Anchor: "sub"},
	}

	got := generateTOC(headings)
	want := "- [Title](#title)\n  - [Sub](#sub)\n"
	if got != want {
		t.Errorf("generateTOC() = %q, want %q", got, want)
	}
}
