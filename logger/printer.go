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

func panicRed() *parameterizableLogger {
	return panicing(color.New(color.Bold, color.FgRed))
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
