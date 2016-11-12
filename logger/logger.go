package logger

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
)

var (
	successCol = color.New(color.Bold, color.FgGreen)

	// SuccessSprintf can be used for colored formatting.
	SuccessSprintf = successCol.SprintfFunc()

	warnCol = color.New(color.Bold, color.FgYellow)

	// WarnSprintf can be used for colored formatting.
	WarnSprintf = warnCol.SprintfFunc()

	// WarnSprint can be used for colored formatting.
	WarnSprint = warnCol.SprintFunc()

	panicCol = color.New(color.Bold, color.FgRed)

	// PanicSprintf can be used for colored formatting.
	PanicSprintf = panicCol.SprintfFunc()
)

// Success logging in green.
func Success(v ...interface{}) {
	log.Printf("%s", SuccessSprintf(strings.Repeat("%s ", len(v)), v...))
}

// Warn logging in yellow.
func Warn(v ...interface{}) {
	log.Printf("%s", WarnSprint(v...))
}

// SuccessNoTimestamp logging in green, no timestamp.
func SuccessNoTimestamp(v ...interface{}) {
	fmt.Printf("%s\n", SuccessSprintf(strings.Repeat("%s ", len(v)), v...))
}

// Panic logging in red, followed by panic.
func Panic(v ...interface{}) {
	log.Panic(PanicSprintf(strings.Repeat("%s ", len(v)), v...))
}
