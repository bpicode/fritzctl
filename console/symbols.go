package console

import (
	"strings"

	"github.com/fatih/color"
)

const (
	checkX = "\u2718"
	checkV = "\u2714"
	checkQ = "?"
)

var (
	red    = color.New(color.Bold, color.FgRed).SprintfFunc()
	green  = color.New(color.Bold, color.FgGreen).SprintfFunc()
	yellow = color.New(color.Bold, color.FgYellow).SprintfFunc()
)

// IntToCheckmark returns a string with ansi color instruction
// characters:
// a red ✘ if the argument is zero,
// a green ✔ otherwise.
func IntToCheckmark(i int) string {
	if i == 0 {
		return redX()
	}
	return greenV()
}

// StringToCheckmark returns a string with ansi color instruction
// characters:
// a red ✘ if the argument is "0",
// a yellow ? if the argument is "",
// a green ✔ otherwise.
func StringToCheckmark(s string) string {
	str := strings.TrimSpace(s)
	if str == "" {
		return yellowQ()
	} else if str == "0" {
		return redX()
	} else {
		return greenV()
	}
}

func redX() string {
	return red(checkX)
}

func greenV() string {
	return green(checkV)
}

func yellowQ() string {
	return yellow(checkQ)
}
