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
		{"keeps underscores", "Foo_bar", "foo_bar"},
		{"keeps hyphens", "pre-commit hooks", "pre-commit-hooks"},
		{"keeps non-ASCII letters", "Über Café", "über-café"},
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

func TestParseMarkdownEdgeCases(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "doc.md")
	content := "# Title\n" +
		"\n" +
		"```bash\n" +
		"# not a heading\n" +
		"```\n" +
		"\n" +
		"~~~~\n" +
		"## also not a heading\n" +
		"~~~\n" +
		"## still inside the block\n" +
		"~~~~\n" +
		"\n" +
		"## Usage\n" +
		"## Usage\n" +
		"### Closed heading ###\n"
	if err := os.WriteFile(target, []byte(content), 0644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	headings, err := parseMarkdown(target, 3)
	if err != nil {
		t.Fatalf("parseMarkdown() error = %v", err)
	}

	want := []Heading{
		{Level: 1, Text: "Title", Anchor: "title"},
		{Level: 2, Text: "Usage", Anchor: "usage"},
		{Level: 2, Text: "Usage", Anchor: "usage-1"},
		{Level: 3, Text: "Closed heading", Anchor: "closed-heading"},
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

func TestReplaceTOC(t *testing.T) {
	toc := "- [A](#a)\n"

	t.Run("replaces the text between the markers", func(t *testing.T) {
		in := "# Doc\n\n<!-- toc -->\nold\n<!-- tocstop -->\n\n## A\n"
		got, err := replaceTOC(in, toc)
		if err != nil {
			t.Fatalf("replaceTOC() error = %v", err)
		}
		want := "# Doc\n\n<!-- toc -->\n\n- [A](#a)\n\n<!-- tocstop -->\n\n## A\n"
		if got != want {
			t.Errorf("replaceTOC() = %q, want %q", got, want)
		}

		again, err := replaceTOC(got, toc)
		if err != nil || again != got {
			t.Errorf("replaceTOC() is not idempotent: %q, %v", again, err)
		}
	})

	t.Run("markers inside other text are ignored", func(t *testing.T) {
		in := "Use `<!-- toc -->` and `<!-- tocstop -->`.\n<!-- toc -->\n<!-- tocstop -->\n"
		got, err := replaceTOC(in, toc)
		if err != nil {
			t.Fatalf("replaceTOC() error = %v", err)
		}
		want := "Use `<!-- toc -->` and `<!-- tocstop -->`.\n<!-- toc -->\n\n- [A](#a)\n\n<!-- tocstop -->\n"
		if got != want {
			t.Errorf("replaceTOC() = %q, want %q", got, want)
		}
	})

	t.Run("missing markers are an error", func(t *testing.T) {
		for _, in := range []string{"no markers", "<!-- toc -->\nonly start\n", "<!-- tocstop -->\n<!-- toc -->\n", "<!-- toc -->", "text <!-- toc -->\n<!-- tocstop -->\n"} {
			if _, err := replaceTOC(in, toc); err == nil {
				t.Errorf("replaceTOC(%q) expected an error", in)
			}
		}
	})
}

func TestWriteTOC(t *testing.T) {
	target := filepath.Join(t.TempDir(), "doc.md")
	if err := os.WriteFile(target, []byte("<!-- toc -->\n<!-- tocstop -->\n"), 0640); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	if err := writeTOC(target, "- [A](#a)\n"); err != nil {
		t.Fatalf("writeTOC() error = %v", err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(data) != "<!-- toc -->\n\n- [A](#a)\n\n<!-- tocstop -->\n" {
		t.Errorf("file content = %q", data)
	}
	if info, _ := os.Stat(target); info.Mode().Perm() != 0640 {
		t.Errorf("mode = %v, want 0640", info.Mode().Perm())
	}
}
