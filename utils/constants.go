// Package utils provides utility constants and functions for the gextend-bash project.
package utils

import (
	"github.com/fatih/color"
)

// Color and formatting constants used throughout the project.
const (
	// MAXWITH caps the terminal width used for status output, so lines stay
	// readable on very wide terminals.
	MAXWITH int = 160
	// WIDTHSUBS defines the width for substitutions in output formatting.
	WIDTHSUBS int = 20
	// TitleColor is used for titles in output.
	TitleColor color.Attribute = color.FgMagenta
	// InfoNameColor is used for informational name fields.
	InfoNameColor color.Attribute = color.Bold
	// InfoValueColor is used for informational value fields.
	InfoValueColor color.Attribute = color.FgYellow
	// BranchDefaultColor is used for default branch display.
	BranchDefaultColor color.Attribute = color.FgBlue
	// BranchChangedColor is used for changed branch display.
	BranchChangedColor color.Attribute = color.FgYellow

	// SuccessColor is used for success messages.
	SuccessColor color.Attribute = color.FgGreen
	// FailureColor is used for failure or error messages.
	FailureColor color.Attribute = color.FgRed
	// NoticeColor is used for informational notice messages.
	NoticeColor color.Attribute = color.FgWhite
	// WarningColor is used for warning messages.
	WarningColor color.Attribute = color.FgYellow
	// UnknownColor is used for messages with an unrecognized status.
	UnknownColor color.Attribute = color.FgYellow
)
