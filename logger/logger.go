package logger

import (
	"log"
	"strings"

	"github.com/bpicode/fritzctl/flags"
	"github.com/fatih/color"
)

type printer func(v ...interface{})

type parameterizableLogger struct {
	print printer
}

func colored(color *color.Color) *parameterizableLogger {
	sprintLnFunc := color.SprintlnFunc()
	return &parameterizableLogger{
		print: func(v ...interface{}) {
			sprinted := sprintLnFunc(v...)
			trimmed := strings.TrimSpace(sprinted)
			log.Print(trimmed)
		},
	}
}

func plain() *parameterizableLogger {
	return &parameterizableLogger{
		print: func(v ...interface{}) {
			log.Println(v...)
		},
	}
}

func panicing(color *color.Color) *parameterizableLogger {
	sprintLnFunc := color.SprintlnFunc()
	return &parameterizableLogger{
		print: func(v ...interface{}) {
			sprinted := sprintLnFunc(v...)
			trimmed := strings.TrimSpace(sprinted)
			log.Panic(trimmed)
		},
	}
}

func nop() *parameterizableLogger {
	return &parameterizableLogger{
		print: func(v ...interface{}) {
		},
	}
}

var (
	debug, info, success, warn, panicLog *parameterizableLogger

	logLvlFlagPtr = flags.String("loglevel", "info", "set the loglevel during execution")
)

func init() {
	SetupLoggers()
}

// SetupLoggers configures the different loggers according to the
// log level flag. The default is to log info, success, warn and panic.
func SetupLoggers() {
	configureLevelsDefault()
	logLvl := strings.ToLower(*logLvlFlagPtr)
	configureLevelsByFlag(logLvl)
}

func configureLevelsDefault() {
	debug = nop()
	info = plain()
	success = colored(color.New(color.Bold, color.FgGreen))

	warn = colored(color.New(color.Bold, color.FgYellow))
	panicLog = panicing(color.New(color.Bold, color.FgRed))
}

func configureLevelsByFlag(logLvl string) {

	switch logLvl {
	case "debug":
		debug = colored(color.New(color.Bold, color.FgBlack))
	case "warn":
		info = nop()
		success = nop()
	case "error":
		info = nop()
		success = nop()
		warn = nop()
	case "none":
		info = nop()
		success = nop()
		warn = nop()
		panicLog = nop()
	}
}

// Debug logging.
func Debug(v ...interface{}) {
	debug.print(v...)
}

// Info logging.
func Info(v ...interface{}) {
	info.print(v...)
}

// Success logging in green.
func Success(v ...interface{}) {
	success.print(v...)
}

// Warn logging in yellow.
func Warn(v ...interface{}) {
	warn.print(v...)
}

// Panic logging in red, followed by panic.
func Panic(v ...interface{}) {
	panicLog.print(v...)
}
