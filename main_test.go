package main

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// TestBlankRun runs the app without args.
func TestBlankRun(t *testing.T) {
	defer func() {
		exit = os.Exit
	}()
	exit = func(code int) {}
	main()
}

// TestHelp runs the app with the help flag.
func TestHelp(t *testing.T) {
	osArgsBefore := os.Args
	exit = func(code int) {}
	defer func() {
		os.Args = osArgsBefore
		exit = os.Exit
	}()
	os.Args = []string{"app_called_in_test_scope", "--help"}
	main()

}

// TestVersion runs the app with the version flag.
func TestVersion(t *testing.T) {
	osArgsBefore := os.Args
	exit = func(code int) {}
	defer func() {
		os.Args = osArgsBefore
		exit = os.Exit
	}()
	os.Args = []string{"app_called_in_test_scope", "--version"}
	main()
}

// TestDetermineExitCode tests the exit code determination.
func TestDetermineExitCode(t *testing.T) {
	assert.Equal(t, 0, determineExitCode(nil))
	assert.Equal(t, 1, determineExitCode("an error"))
}

// TestSanitizeStack exercises the stack sanitize.
func TestSanitizeStack(t *testing.T) {
	assert.Len(t, sanitizedStack(nil), 0)
	assert.Equal(t, sanitizedStack("error"), []string{"error"})
	assert.Equal(t, sanitizedStack(errors.Errorf("error")), []string{"error"})
	assert.Equal(t, sanitizedStack(errors.Wrapf(errors.Errorf("root"), "parent")), []string{"parent:", "root"})
}
