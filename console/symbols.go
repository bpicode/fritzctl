package console

import (
	"strings"
)

const (
	checkX = "\u2718"
	checkV = "\u2714"
	checkQ = "?"
)

// Checkmark type is a string with some functions attached.
type Checkmark string

func (c Checkmark) String() string {
	return string(c)
}

// Inverse returns the opposite of the given Checkmark:
// a red ✘ if the argument is a green ✔ and vice versa,
// a yellow ? otherwise.
func (c Checkmark) Inverse() Checkmark {
	switch c {
	case redX():
		return greenV()
	case greenV():
		return redX()
	}
	return c
}

// IntToCheckmark returns a string with ansi color instruction
// characters:
// a red ✘ if the argument is zero,
// a green ✔ otherwise.
func IntToCheckmark(i int) string {
	return Itoc(i).String()
}

// Itoc returns a Checkmark from an int.
func Itoc(i int) Checkmark {
	if i == 0 {
		return redX()
	}
	return greenV()
}

// Btoc returns a Checkmark from a boolean, red ✘ if the argument is false, a green ✔ otherwise.
func Btoc(b bool) Checkmark {
	if b {
		return greenV()
	}
	return redX()
}

// StringToCheckmark returns a string with ansi color instruction
// characters:
// a red ✘ if the argument is "0",
// a yellow ? if the argument is "",
// a green ✔ otherwise.
func StringToCheckmark(s string) string {
	return Stoc(s).String()
}

// Stoc returns a Checkmark from a string.
func Stoc(s string) Checkmark {
	str := strings.TrimSpace(s)
	if str == "" {
		return yellowQ()
	} else if str == "0" {
		return redX()
	} else {
		return greenV()
	}
}

func redX() Checkmark {
	return Checkmark(Red(checkX))
}

func greenV() Checkmark {
	return Checkmark(Green(checkV))
}

func yellowQ() Checkmark {
	return Checkmark(Yellow(checkQ))
}
