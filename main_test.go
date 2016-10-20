package main

import (
	"os"
	"testing"
)

// TestBlankRun unit test.
func TestBlankRun(t *testing.T) {
	main()
}

// TestHelp unit test.
func TestHelp(t *testing.T) {
	os.Args = append(os.Args, "--help")
	main()
}

// TestVersion unit test.
func TestVersion(t *testing.T) {
	os.Args = append(os.Args, "--version")
	main()
}
