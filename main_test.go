package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBlankRun runs the app without args.
func TestBlankRun(t *testing.T) {
	exitAdvice = func(code int) {}
	main()
	exitAdvice = os.Exit
}

// TestHelp runs the app with the help flag.
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

// TestVersion runs the app with the version flag.
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

// TestDetermineExitCode tests the exit code determination.
func TestDetermineExitCode(t *testing.T) {
	assert.Equal(t, 0, determineExitCode(nil))
	assert.Equal(t, 1, determineExitCode("an error"))
}
