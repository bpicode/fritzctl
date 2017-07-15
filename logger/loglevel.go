package logger

import (
	"fmt"
	"strings"
)

type printers struct {
	debug, info, success, warn, panicLog *parameterizableLogger
}

var ls printers

func init() {
	ConfigureLogLevel("info")
}

// Debug logging.
func Debug(v ...interface{}) {
	ls.debug.print(v...)
}

// Info logging.
func Info(v ...interface{}) {
	ls.info.print(v...)
}

// Success logging in green.
func Success(v ...interface{}) {
	ls.success.print(v...)
}

// Warn logging in yellow.
func Warn(v ...interface{}) {
	ls.warn.print(v...)
}

// Panic logging in red, followed by panic.
func Panic(v ...interface{}) {
	ls.panicLog.print(v...)
}

type levelLookupTable map[string]*printers

var levelNames = levelLookupTable{
	"debug": &printers{debug: dark(), info: plain(), success: green(), warn: yellow(), panicLog: panicRed()},
	"info":  &printers{debug: nop(), info: plain(), success: green(), warn: yellow(), panicLog: panicRed()},
	"warn":  &printers{debug: nop(), info: nop(), success: nop(), warn: yellow(), panicLog: panicRed()},
	"error": &printers{debug: nop(), info: nop(), success: nop(), warn: nop(), panicLog: panicRed()},
	"none":  &printers{debug: nop(), info: nop(), success: nop(), warn: nop(), panicLog: nop()},
}

// ConfigureLogLevel configures the loglevel identified by its name. It returns an error if the given name is unknown.
func ConfigureLogLevel(name string) error {
	l, err := byName(name)
	if err != nil {
		return fmt.Errorf("error determining loglevel details for name '%s': %v", name, err)
	}
	ls = *l
	return nil
}

func byName(name string) (*printers, error) {
	level, ok := levelNames[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("'%s' is not a valid loglevel, possible values are %s", name, levelNames.keys())
	}
	return level, nil
}

func (l *levelLookupTable) keys() []string {
	keys := []string{}
	for k := range *l {
		keys = append(keys, k)
	}
	return keys
}
