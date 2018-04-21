package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generate_notice(t *testing.T) {
	var argsOrig []string
	copy(argsOrig, os.Args)
	defer func() { os.Args = argsOrig }()
	os.Args = []string{"analice", "generate", "notice", "../.."}
	assert.NotPanics(t, main)
}

func Test_generate_copyright(t *testing.T) {
	var argsOrig []string
	copy(argsOrig, os.Args)
	defer func() { os.Args = argsOrig }()
	os.Args = []string{"analice", "generate", "copyright", "../.."}
	assert.NotPanics(t, main)
}

func Test_generate_with_err(t *testing.T) {
	var argsOrig []string
	copy(argsOrig, os.Args)
	defer func() { os.Args = argsOrig }()

	var exOnErr = exitOnErr
	defer func() { exitOnErr = exOnErr }()
	exitOnErr = func(v ...interface{}) { panic(v) }

	os.Args = []string{"analice", "generate", "notice", "/good/luck/with/that"}
	assert.Panics(t, main)

	os.Args = []string{"analice", "generate", "copyright", "/not/gonna/work"}
	assert.Panics(t, main)
}
