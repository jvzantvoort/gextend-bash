package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/mitchellh/go-wordwrap"
)

const (
	ConsoleWidth = 80
	MinBoxWidth  = 8
	TopLeft      = "╔"
	TopRight     = "╗"
	BottomLeft   = "╚"
	BottomRight  = "╝"
	Horizontal   = "═"
	Vertical     = "║"
	Rainbow      = "rainbow"
)

var namedColors = map[string]color.Attribute{
	"black":     color.FgBlack,
	"red":       color.FgRed,
	"green":     color.FgGreen,
	"yellow":    color.FgYellow,
	"blue":      color.FgBlue,
	"magenta":   color.FgMagenta,
	"cyan":      color.FgCyan,
	"white":     color.FgWhite,
	"hiblack":   color.FgHiBlack,
	"hired":     color.FgHiRed,
	"higreen":   color.FgHiGreen,
	"hiyellow":  color.FgHiYellow,
	"hiblue":    color.FgHiBlue,
	"himagenta": color.FgHiMagenta,
	"hicyan":    color.FgHiCyan,
	"hiwhite":   color.FgHiWhite,
}

var rainbowColors = []color.Attribute{
	color.FgRed,
	color.FgYellow,
	color.FgGreen,
	color.FgCyan,
	color.FgBlue,
	color.FgMagenta,
}

func openStdinOrFile(args []string) io.Reader {
	if len(args) > 0 {
		return strings.NewReader(strings.Join(args, " "))
	}
	return os.Stdin
}

// colorize applies the requested color to a single line of the banner.
// An unknown or empty color name leaves the line unmodified.
func colorize(line, colorName string) string {
	switch colorName {
	case "", "none":
		return line
	case Rainbow:
		return rainbow(line)
	default:
		if attr, ok := namedColors[colorName]; ok {
			return color.New(attr).Sprint(line)
		}
		return line
	}
}

// rainbow colors each character of a line, cycling through rainbowColors.
func rainbow(line string) string {
	var b strings.Builder
	runes := []rune(line)
	for i, r := range runes {
		attr := rainbowColors[i%len(rainbowColors)]
		b.WriteString(color.New(attr).Sprint(string(r)))
	}
	return b.String()
}

func buildLine(instr string, boxWidth int, left bool) string {
	width := boxWidth - 4
	if len(instr) > width {
		instr = instr[:width]
	}
	restlen := width - len(instr)
	leftpad, rightpad := restlen/2, restlen-restlen/2
	if left {
		leftpad, rightpad = 0, restlen
	}
	return Vertical + " " + strings.Repeat(" ", leftpad) + instr + strings.Repeat(" ", rightpad) + " " + Vertical
}

// PrintBanner prints instr inside a box, colored with colorName.
// hpad is the horizontal padding (columns) kept between the terminal width
// and the box; vpad is the vertical padding (blank lines) kept between the
// box border and the text. When left is true, text is left-aligned instead
// of centered.
func PrintBanner(instr, colorName string, hpad, vpad int, left bool) {
	hpad = max(hpad, 0)
	vpad = max(vpad, 0)

	boxWidth := max(ConsoleWidth-2*hpad, MinBoxWidth)
	margin := strings.Repeat(" ", hpad)

	instr = strings.TrimSpace(instr)
	wrapped := wordwrap.WrapString(instr, uint(boxWidth-4))
	lines := strings.Split(wrapped, "\n")

	top := TopLeft + strings.Repeat(Horizontal, boxWidth-2) + TopRight
	bottom := BottomLeft + strings.Repeat(Horizontal, boxWidth-2) + BottomRight
	blank := Vertical + strings.Repeat(" ", boxWidth-2) + Vertical

	print := func(s string) {
		fmt.Println(margin + colorize(s, colorName))
	}

	print(top)
	for i := 0; i < vpad; i++ {
		print(blank)
	}
	for _, line := range lines {
		print(buildLine(line, boxWidth, left))
	}
	for i := 0; i < vpad; i++ {
		print(blank)
	}
	print(bottom)
}

func usage() {
	names := make([]string, 0, len(namedColors)+1)
	for name := range namedColors {
		names = append(names, name)
	}
	names = append(names, Rainbow)
	sort.Strings(names)

	fmt.Fprintf(os.Stderr, "Usage: %s [-c|--color name] [-x|--hpad n] [-y|--vpad n] [-l|--left] [text...]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Reads text from the arguments, or from stdin when no arguments are given,\n")
	fmt.Fprintf(os.Stderr, "and prints it inside a box.\n\n")
	fmt.Fprintf(os.Stderr, "Available colors: %s\n\n", strings.Join(names, ", "))
	fmt.Fprintf(os.Stderr, "--hpad sets the padding (columns) kept between the terminal width and the box.\n")
	fmt.Fprintf(os.Stderr, "--vpad sets the padding (blank lines) kept between the box border and the text.\n")
	fmt.Fprintf(os.Stderr, "--left left-aligns the text instead of centering it.\n")
}

func main() {
	var colorName string
	var hpad, vpad int
	var left bool
	flag.StringVar(&colorName, "color", "", "border/text color, or 'rainbow'")
	flag.StringVar(&colorName, "c", "", "shorthand for --color")
	flag.IntVar(&hpad, "hpad", 0, "horizontal padding between the terminal width and the box")
	flag.IntVar(&hpad, "x", 0, "shorthand for --hpad")
	flag.IntVar(&vpad, "vpad", 0, "vertical padding between the box border and the text")
	flag.IntVar(&vpad, "y", 0, "shorthand for --vpad")
	flag.BoolVar(&left, "left", false, "left-align text instead of centering it")
	flag.BoolVar(&left, "l", false, "shorthand for --left")
	flag.Usage = usage
	flag.Parse()

	r := openStdinOrFile(flag.Args())
	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	PrintBanner(string(b), strings.ToLower(strings.TrimSpace(colorName)), hpad, vpad, left)
}

// vim: noexpandtab filetype=go
