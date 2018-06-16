package logger

import (
	"fmt"
	"strings"

	"github.com/bpicode/fritzctl/internal/errors"
)

type printers struct {
	debug, info, success, warn, error *parameterizableLogger
}

var ls printers

func init() {
	configureLogLevel("info")
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

// Error logging in red.
func Error(v ...interface{}) {
	ls.error.print(v...)
}

type levelLookupTable map[string]*printers

var levelNames = levelLookupTable{
	"debug": &printers{debug: dark(), info: plain(), success: green(), warn: yellow(), error: red()},
	"info":  &printers{debug: nop(), info: plain(), success: green(), warn: yellow(), error: red()},
	"warn":  &printers{debug: nop(), info: nop(), success: nop(), warn: yellow(), error: red()},
	"error": &printers{debug: nop(), info: nop(), success: nop(), warn: nop(), error: red()},
	"none":  &printers{debug: nop(), info: nop(), success: nop(), warn: nop(), error: nop()},
}

// configureLogLevel configures the loglevel identified by its name. It returns an error if the given name is unknown.
func configureLogLevel(name string) error {
	l, err := byName(name)
	if err != nil {
		return errors.Wrapf(err, "error determining loglevel details for name '%s'", name)
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
	var keys []string
	for k := range *l {
		keys = append(keys, k)
	}
	return keys
}
