package main

import (
	"os"
	"testing"
)

// TestBlankRun runs the app without args.
func TestBlankRun(t *testing.T) {
	defer func() {
		exitAdvice = os.Exit
	}()
	exitAdvice = func(code int) {}
	main()
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
