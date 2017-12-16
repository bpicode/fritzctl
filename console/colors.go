package console

import "github.com/fatih/color"

var (
	// Red can be used as Sprintf, where the output it wrapped in escape characters which will render the text red in terminals.
	Red = color.New(color.Bold, color.FgRed).SprintfFunc()
	// Green can be used as Sprintf, where the output it wrapped in escape characters which will render the text green in terminals.
	Green = color.New(color.Bold, color.FgGreen).SprintfFunc()
	// Yellow can be used as Sprintf, where the output it wrapped in escape characters which will render the text yellow in terminals.
	Yellow = color.New(color.Bold, color.FgYellow).SprintfFunc()
	// Cyan can be used as Sprintf, where the output it wrapped in escape characters which will render the text cyan in terminals.
	Cyan = color.New(color.FgCyan).SprintfFunc()
)
