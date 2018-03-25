package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	argsOrig := os.Args[:]
	defer func() { os.Args = argsOrig }()
	os.Args = []string{"notice2copyright", "../../", "MIT License (Expat)"}
	assert.NotPanics(t, main)
}
