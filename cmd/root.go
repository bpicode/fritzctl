package cmd

import (
	"fmt"

	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any sub-commands.
var RootCmd = &cobra.Command{
	Use: "fritzctl",
}

type loglevelValue struct {
	level string
}

// Type returns a name for the value type.
func (l *loglevelValue) Type() string {
	return "loglevel"
}

// String converts a loglevelValue to human-readable format.
func (l *loglevelValue) String() string {
	return l.level
}

// Set configures the loglevel for the application.
func (l *loglevelValue) Set(val string) error {
	err := logger.ConfigureLogLevel(val)
	if err != nil {
		return fmt.Errorf("cannot apply loglevel configuration for value '%s': %v", val, err)
	}
	return nil
}

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().Var(&loglevelValue{}, "loglevel", "logging verbosity")
}
