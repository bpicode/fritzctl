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
	blankFree := sanitizedStack(r)
	logger.Error(strings.Join(blankFree, "\n  "))
}

func sanitizedStack(r interface{}) []string {
	stack := stack(r)
	duplicateFree := rmDuplicateMsgs(stack)
	blankFree := rmBlanks(duplicateFree)
	return blankFree
}

func stack(e interface{}) []string {
	if e == nil {
		return nil
	}
	type causer interface {
		error
		Cause() error
	}
	var frames []string
	csr, ok := e.(causer)
	if !ok {
		return []string{fmt.Sprint(e)}
	}
	frames = append(frames, csr.Error())
	frames = append(frames, stack(csr.Cause())...)
	return frames
}

func rmDuplicateMsgs(frames []string) []string {
	for i := len(frames) - 1; i >= 0; i-- {
		if i == 0 || frames[i] == "" {
			continue
		}
		for j := i - 1; j >= 0; j-- {
			split := strings.Split(frames[j], frames[i])
			frames[j] = split[0]
		}
	}
	return frames
}

func rmBlanks(strs []string) []string {
	var noBlanks []string
	for _, s := range strs {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			noBlanks = append(noBlanks, trimmed)
		}
	}
	return noBlanks
}
