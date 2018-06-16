package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/bpicode/fritzctl/internal/errors"
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

//// TestStack exercises the stack traversal.
func TestStack(t *testing.T) {
	assert.Len(t, stack(nil), 0)
	assert.Equal(t, []string{"error"}, stack("error"))
	assert.Equal(t, []string{"error"}, stack(fmt.Errorf("error")))
	assert.Equal(t, []string{"parent:", "root"}, stack(errors.Wrapf(fmt.Errorf("root"), "parent")))
}
