package main

import (
	"os"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/cliapp"
	"github.com/bpicode/fritzctl/flags"
)

var (
	exitAdvice = os.Exit
)

func main() {
	defer func() {
		r := recover()
		exitCode := determineExitCode(r)
		exitAdvice(exitCode)
	}()
	flags.Parse(os.Args[1:])
	c := cliapp.New()
	_, err := c.Run()
	assert.NoError(err, "Error:", err)
}

func determineExitCode(v interface{}) int {
	if v == nil {
		return 0
	}
	return 1
}
