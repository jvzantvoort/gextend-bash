// Package utils provides utility constants and functions for the gextend-bash project.
package utils

import (
	"github.com/fatih/color"
)

// Color and formatting constants used throughout the project.
const (
	// WIDTHSUBS defines the width for substitutions in output formatting.
	WIDTHSUBS          int             = 20
	// TitleColor is used for titles in output.
	TitleColor         color.Attribute = color.FgMagenta
	// InfoNameColor is used for informational name fields.
	InfoNameColor      color.Attribute = color.Bold
	// InfoValueColor is used for informational value fields.
	InfoValueColor     color.Attribute = color.FgYellow
	// BranchDefaultColor is used for default branch display.
	BranchDefaultColor color.Attribute = color.FgBlue
	// BranchChangedColor is used for changed branch display.
	BranchChangedColor color.Attribute = color.FgYellow

	// SuccessColor is used for success messages.
	SuccessColor color.Attribute = color.FgGreen
	// FailureColor is used for failure or error messages.
	FailureColor color.Attribute = color.FgRed
)
