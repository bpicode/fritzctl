package logger

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
)

var (
	infoCol     = color.New(color.Bold, color.FgGreen)
	infoSprintf = infoCol.SprintfFunc()

	panicCol     = color.New(color.Bold, color.FgRed)
	panicSprintf = panicCol.SprintfFunc()
)

// Info logging in greeen.
func Info(v ...interface{}) {
	log.Printf("%s", infoSprintf(strings.Repeat("%s ", len(v)), v...))
}

// InfoNoTimestamp logging in greeen, no timestamp.
func InfoNoTimestamp(v ...interface{}) {
	fmt.Printf("%s\n", infoSprintf(strings.Repeat("%s ", len(v)), v...))
}

// Panic logging in red, followed by panic.
func Panic(v ...interface{}) {
	log.Panic(panicSprintf(strings.Repeat("%s ", len(v)), v...))
}
