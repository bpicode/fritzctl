package main

import (
	"os"

	"github.com/bpicode/fritzctl/cmd"
	"github.com/bpicode/fritzctl/logger"
)

var (
	exitAdvice = os.Exit
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		logger.Error(err)
		exitAdvice(1)
	}
}
