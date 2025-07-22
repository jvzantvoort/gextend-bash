package colors

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Structure to hold the functionality of cprint
type Cprint struct {
	Format  string
	color   color.Color
	token   rune
	Colors  map[string]color.Attribute
	tokens  map[string]rune
	formats map[string]string
}

func (c *Cprint) InitColors() {
	colors := make(map[string]color.Attribute)

	// regular colors
	colors["black"] = color.FgBlack
	colors["blue"] = color.FgBlue
	colors["cyan"] = color.FgCyan
	colors["green"] = color.FgGreen
	colors["magenta"] = color.FgMagenta
	colors["red"] = color.FgRed
	colors["white"] = color.FgWhite
	colors["yellow"] = color.FgYellow

	// status
	for _, x := range []string{"nok"} {
		colors[x] = color.FgRed
	}

	for _, x := range []string{"ok", "oke", "profile", "platform", "workspace"} {
		colors[x] = color.FgGreen
	}

	for _, x := range []string{"debug"} {
		colors[x] = color.FgBlue
	}

	for _, x := range []string{"warn"} {
		colors[x] = color.FgYellow
	}
	c.Colors = colors
}

func (c *Cprint) InitTokens() {

	tokens := make(map[string]rune)

	for _, x := range []string{"ok", "oke"} {
		tokens[x] = HEAVY_CHECK_MARK
	}

	for _, x := range []string{"nok", "err", "error"} {
		tokens[x] = WARNING_SIGN
	}

	for _, x := range []string{"warn", "err", "error"} {
		tokens[x] = EXCLAMATION_MARK
	}

	for _, x := range []string{"profile", "platform", "workspace"} {
		tokens[x] = BLACK_DIAMOND_SUIT
	}
	c.tokens = tokens

}

func (c *Cprint) InitFormats() {

	formats := make(map[string]string)
	formats["defaults"] = "%s %s\n"
	formats["profile"] = "%s profile (%s) sourced\n"
	for _, x := range []string{"platform", "workspace"} {
		formats[x] = fmt.Sprintf("%%s %s %%s sourced\n", x)
	}

	c.formats = formats
}

func (c *Cprint) SetColor(carg string) {
	if lcolor, ok := c.Colors[carg]; ok {
		c.color.Add(lcolor)
		c.color.Add(color.Bold)
	}
}

func (c *Cprint) SetToken(carg string) {
	if ltoken, ok := c.tokens[carg]; ok {
		c.token = ltoken
	} else {
		c.token = BLACK_DIAMOND_SUIT
	}
}

func (c *Cprint) SetFormat(carg string) {
	if lformat, ok := c.formats[carg]; ok {
		c.Format = lformat
	} else {
		c.Format = c.formats["defaults"]
	}

}

func (c *Cprint) Print(args ...string) {
	carg := args[0]
	c.SetToken(carg)
	c.SetColor(carg)
	c.SetFormat(carg)
	fmt.Printf(c.Format, c.color.Sprintf("%c", c.token), strings.Join(args[1:], " "))
}

func NewCprint() *Cprint {
	retv := &Cprint{}
	retv.color = *color.New()
	retv.InitColors()
	retv.InitTokens()
	retv.InitFormats()
	return retv
}
