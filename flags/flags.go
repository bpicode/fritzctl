package flags

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bpicode/fritzctl/logger"
)

var set *flag.FlagSet

var (
	versionPresent *bool
	vPresent       *bool
	helpPresent    *bool
	hPresent       *bool
)

func init() {
	reset()
}

func reset() {
	set = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	set.SetOutput(ioutil.Discard)
	versionPresent = set.Bool("version", false, "--version")
	vPresent = set.Bool("v", false, "-v")
	helpPresent = set.Bool("help", false, "--help")
	hPresent = set.Bool("h", false, "-h")
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func String(name string, value string, usage string) *string {
	return set.String(name, value, usage)
}

// Parse parses the command-line flags from os.Args[1:].  Must be called
// after all flags are defined and before flags are accessed by the program.
func Parse(args []string) {
	// Ignore errors; CommandLine is set for ExitOnError.
	set.Parse(args)
}

// Args returns the non-flag command-line arguments.
func Args() []string {
	args := set.Args()
	if *vPresent || *versionPresent {
		args = append(args, "--version")
	}
	if *hPresent || *helpPresent {
		args = append(args, "--help")
	}
	return args
}

// Loglevel represents a Value for different logging configs.
type Loglevel struct {
	level string
}

// Type returns a name for the value type.
func (l *Loglevel) Type() string {
	return "string"
}

// String converts a Loglevel to human-readable format.
func (l *Loglevel) String() string {
	return l.level
}

// Set configures the loglevel for the application.
func (l *Loglevel) Set(val string) error {
	err := logger.ConfigureLogLevel(val)
	if err != nil {
		return fmt.Errorf("cannot apply loglevel configuration for value '%s': %v", val, err)
	}
	return nil
}
