package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBlankRun unit test.
func TestBlankRun(t *testing.T) {
	exitAdvice = func(code int) {}
	main()
	exitAdvice = os.Exit
}

// TestHelp unit test.
func TestHelp(t *testing.T) {
	osArgsBefore := os.Args
	exitAdvice = func(code int) {}
	defer func() {
		os.Args = osArgsBefore
		exitAdvice = os.Exit
	}()
	os.Args = []string{"app_called_in_test_scope", "--help"}
	main()

}

// TestVersion unit test.
func TestVersion(t *testing.T) {
	osArgsBefore := os.Args
	exitAdvice = func(code int) {}
	defer func() {
		os.Args = osArgsBefore
		exitAdvice = os.Exit
	}()
	os.Args = []string{"app_called_in_test_scope", "--version"}
	main()
}

// TestDetermineExitCode is a unit test.
func TestDetermineExitCode(t *testing.T) {
	assert.Equal(t, 0, determineExitCode(nil))
	assert.Equal(t, 1, determineExitCode("an error"))
}
