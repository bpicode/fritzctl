package flags

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

// TestFlagsWithEmptyArgsArray is a unit test for this package.
func TestFlagsWithEmptyArgsArray(t *testing.T) {
	reset()
	Parse(os.Args[1:])
	args := Args()
	assert.NotNil(t, args)
}

// TestFlagsWithVersionFlag is a unit test for this package.
func TestFlagsWithVersionFlag(t *testing.T) {
	reset()
	Parse([]string{"-version"})
	args := Args()
	assert.Contains(t, args, "--version")
}

// TestFlagsWithVFlag is a unit test for this package.
func TestFlagsWithVFlag(t *testing.T) {
	reset()
	Parse([]string{"-v"})
	args := Args()
	assert.Contains(t, args, "--version")
}

// TestFlagsWithHelpFlag is a unit test for this package.
func TestFlagsWithHelpFlag(t *testing.T) {
	reset()
	Parse([]string{"--help"})
	args := Args()
	assert.Contains(t, args, "--help")
}

// TestFlagsWithHFlag is a unit test for this package.
func TestFlagsWithHFlag(t *testing.T) {
	reset()
	Parse([]string{"-h"})
	args := Args()
	assert.Contains(t, args, "--help")
}

// TestRegisterStringButValueNotPresent is a unit test for this package.
func TestRegisterStringButValueNotPresent(t *testing.T) {
	reset()
	sPtr := String("myvar", "xx", "my variable, default xx")
	Parse([]string{})
	assert.Equal(t, "xx", *sPtr)
}

// TestRegisterStringWithValueNotPresent is a unit test for this package.
func TestRegisterStringWithValueNotPresent(t *testing.T) {
	reset()
	sPtr := String("myvar", "xx", "my variable, default xx")
	Parse([]string{"--myvar=lol"})
	assert.Equal(t, "lol", *sPtr)
}
