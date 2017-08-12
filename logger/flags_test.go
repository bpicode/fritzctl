package logger

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegisterLoglevel is a unit test for setting the loglevel via CLI flag.
func TestRegisterLoglevel(t *testing.T) {
	f := flag.NewFlagSet("TestFlagSet", flag.PanicOnError)
	ll := Loglevel{}
	f.Var(&ll, "loglevel", "logging verbosity specified as "+ll.Type())
	err := f.Parse([]string{"--loglevel=info"})
	assert.NoError(t, err)
}

// TestRegisterLoglevel is a unit test for setting the loglevel via CLI flag.
func TestRegisterLoglevelInvalid(t *testing.T) {
	f := flag.NewFlagSet("TestFlagSet", flag.ContinueOnError)
	ll := Loglevel{}
	f.Var(&ll, "loglevel", "logging verbosity specified as "+ll.Type())
	err := f.Parse([]string{"--loglevel=whistleblower"})
	assert.Error(t, err)
}
