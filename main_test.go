package main

import (
	"os"
	"testing"
)

func TestBlankRun(t *testing.T) {
	main()
}

func TestHelp(t *testing.T) {
	os.Args = append(os.Args, "--help")
	main()
}

func TestVersion(t *testing.T) {
	os.Args = append(os.Args, "--version")
	main()
}
