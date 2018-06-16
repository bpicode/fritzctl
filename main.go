package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bpicode/fritzctl/cmd"
	"github.com/bpicode/fritzctl/logger"
)

var (
	exit = os.Exit
)

func main() {
	defer func() {
		r := recover()
		if r != nil {
			printErr(r)
		}
		exitCode := determineExitCode(r)
		exit(exitCode)
	}()
	err := cmd.RootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func determineExitCode(v interface{}) int {
	if v == nil {
		return 0
	}
	return 1
}

func printErr(r interface{}) {
	st := stack(r)
	logger.Error(strings.Join(st, "\n  "))
}

func stack(e interface{}) []string {
	if e == nil {
		return nil
	}
	type causer interface {
		error
		Cause() error
		Msg() string
	}
	var frames []string
	csr, ok := e.(causer)
	if !ok {
		return []string{fmt.Sprint(e)}
	}
	frames = append(frames, csr.Msg()+":")
	frames = append(frames, stack(csr.Cause())...)
	return frames
}
