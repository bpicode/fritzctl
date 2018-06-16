package logger

import (
	"github.com/bpicode/fritzctl/internal/errors"
)

// Level represents a Value for different logging configs.
type Level struct {
}

// Type returns a name for the value type.
func (l *Level) Type() string {
	return "string"
}

// String represents the default value.
func (l *Level) String() string {
	return "info"
}

// Set configures the loglevel for the application.
func (l *Level) Set(val string) error {
	err := configureLogLevel(val)
	if err != nil {
		return errors.Wrapf(err, "cannot apply loglevel configuration for value '%s'", val)
	}
	return nil
}
