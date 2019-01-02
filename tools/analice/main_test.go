package main

import (
	"errors"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_generate_notice(t *testing.T) {
	var argsOrig []string
	copy(argsOrig, os.Args)
	defer func() { os.Args = argsOrig }()
	os.Args = []string{"analice", "generate", "notice", "github.com/bpicode/fritzctl", "--tests=false"}
	assert.NotPanics(t, main)
}

func Test_generate_copyright(t *testing.T) {
	var argsOrig []string
	copy(argsOrig, os.Args)
	defer func() { os.Args = argsOrig }()
	os.Args = []string{"analice", "generate", "copyright", "github.com/bpicode/fritzctl"}
	assert.NotPanics(t, main)
}

func Test_generate_with_err(t *testing.T) {
	var argsOrig []string
	copy(argsOrig, os.Args)
	defer func() { os.Args = argsOrig }()

	var exOnErr = exitOnErr
	defer func() { exitOnErr = exOnErr }()
	exitOnErr = func(v ...interface{}) { panic(v) }

	rootCmdOrig := *rootCmd
	defer func() { rootCmd = &rootCmdOrig }()
	rootCmd = &cobra.Command{RunE: func(_ *cobra.Command, _ []string) error {
		return errors.New("an error")
	}}

	os.Args = []string{"analice", "generate", "notice", "./..."}
	assert.Panics(t, main)
}
