package colors

import (
	"fmt"
)

const (
	defaultMainColor string = "lightcyan"
	defaultOSColor   string = "green"
	defaultDirColor  string = "yellow"
	ColorEnd         string = "0m"
	ColorBlack       string = "0;30m"
	ColorRed         string = "0;31m"
	ColorGreen       string = "0;32m"
	ColorBrown       string = "0;33m"
	ColorBlue        string = "0;34m"
	ColorPurple      string = "0;35m"
	ColorCyan        string = "0;36m"
	ColorLightGray   string = "0;37m"
	ColorDarkGray    string = "1;30m"
	ColorGray        string = "1:30m"
	ColorLightBlue   string = "1;34m"
	ColorLightCyan   string = "1;36m"
	ColorLightGreen  string = "1;32m"
	ColorLightPurpl  string = "1;35m"
	ColorLightRed    string = "1;31m"
	ColorWhite       string = "1;37m"
	ColorYellow      string = "1;33m"
)

// printc wraps the color definition in escape strings needed.
func printc(color string) string {
	return fmt.Sprintf("\\[\033[%s\\]", color)
}

func ColornameToColorvalue(name string) string {
	retv := printc(ColorWhite)
	if name == "black" {
		retv = printc(ColorBlack)
	} else if name == "blue" {
		retv = printc(ColorBlue)
	} else if name == "brown" {
		retv = printc(ColorBrown)
	} else if name == "cyan" {
		retv = printc(ColorCyan)
	} else if name == "darkgray" {
		retv = printc(ColorDarkGray)
	} else if name == "gray" {
		retv = printc(ColorGray)
	} else if name == "green" {
		retv = printc(ColorGreen)
	} else if name == "lightblue" {
		retv = printc(ColorLightBlue)
	} else if name == "lightcyan" {
		retv = printc(ColorLightCyan)
	} else if name == "lightgray" {
		retv = printc(ColorLightGray)
	} else if name == "lightgreen" {
		retv = printc(ColorLightGreen)
	} else if name == "lightpurple" {
		retv = printc(ColorLightPurpl)
	} else if name == "lightred" {
		retv = printc(ColorLightRed)
	} else if name == "purple" {
		retv = printc(ColorPurple)
	} else if name == "red" {
		retv = printc(ColorRed)
	} else if name == "white" {
		retv = printc(ColorWhite)
	} else if name == "yellow" {
		retv = printc(ColorYellow)
	} else if name == "end" {
		retv = printc(ColorEnd)
	}
	return retv
}
