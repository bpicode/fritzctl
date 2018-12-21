package logger

import (
	"log"
	"strings"

	"github.com/fatih/color"
)

type printer func(v ...interface{})

type parameterizableLogger struct {
	print printer
}

func red() *parameterizableLogger {
	return colored(color.New(color.Bold, color.FgRed))
}

func yellow() *parameterizableLogger {
	return colored(color.New(color.Bold, color.FgYellow))
}

func green() *parameterizableLogger {
	return colored(color.New(color.Bold, color.FgGreen))
}

func dark() *parameterizableLogger {
	return colored(color.New(color.Bold, color.FgBlack))
}

func colored(color *color.Color) *parameterizableLogger {
	sprint := color.SprintFunc()
	return &parameterizableLogger{
		print: func(v ...interface{}) {
			sprinted := sprint(v...)
			trimmed := strings.TrimSpace(sprinted)
			log.Println(trimmed)
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

func nop() *parameterizableLogger {
	return &parameterizableLogger{
		print: func(v ...interface{}) {
		},
	}
}
