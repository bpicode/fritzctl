package main

import (
	"os"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/cmd"
	"github.com/bpicode/fritzctl/logger"
)

var (
	exitAdvice = os.Exit
)

func main() {
	defer func() {
		r := recover()
		logger.Error(r)
		exitCode := determineExitCode(r)
		exitAdvice(exitCode)
	}()
	err := cmd.RootCmd.Execute()
	assert.NoError(err, "error:", err)
}

func determineExitCode(v interface{}) int {
	if v == nil {
		return 0
	}
	return 1
}
