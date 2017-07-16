package flags

import (
	"fmt"

	"github.com/bpicode/fritzctl/logger"
)

// Loglevel represents a Value for different logging configs.
type Loglevel struct {
	level string
}

// Type returns a name for the value type.
func (l *Loglevel) Type() string {
	return "string"
}

// String represents the default value.
func (l *Loglevel) String() string {
	return "info"
}

// Set configures the loglevel for the application.
func (l *Loglevel) Set(val string) error {
	err := logger.ConfigureLogLevel(val)
	if err != nil {
		return fmt.Errorf("cannot apply loglevel configuration for value '%s': %v", val, err)
	}
	return nil
}
