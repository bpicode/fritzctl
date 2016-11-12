package logger

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
)

var (
	infoCol = color.New(color.Bold, color.FgGreen)

	// InfoSprintf can be used for colored formatting.
	InfoSprintf = infoCol.SprintfFunc()

	warnCol = color.New(color.Bold, color.FgYellow)

	// WarnSprintf can be used for colored formatting.
	WarnSprintf = warnCol.SprintfFunc()

	// WarnSprint can be used for colored formatting.
	WarnSprint = warnCol.SprintFunc()

	panicCol = color.New(color.Bold, color.FgRed)

	// PanicSprintf can be used for colored formatting.
	PanicSprintf = panicCol.SprintfFunc()
)

// Info logging in greeen.
func Info(v ...interface{}) {
	log.Printf("%s", InfoSprintf(strings.Repeat("%s ", len(v)), v...))
}

// Warn logging in yellow.
func Warn(v ...interface{}) {
	log.Printf("%s", WarnSprint(v...))
}

// InfoNoTimestamp logging in greeen, no timestamp.
func InfoNoTimestamp(v ...interface{}) {
	fmt.Printf("%s\n", InfoSprintf(strings.Repeat("%s ", len(v)), v...))
}

// Panic logging in red, followed by panic.
func Panic(v ...interface{}) {
	log.Panic(PanicSprintf(strings.Repeat("%s ", len(v)), v...))
}
