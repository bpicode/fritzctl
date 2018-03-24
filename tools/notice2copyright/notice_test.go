package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_projectRegex(t *testing.T) {
	s := "- cobra (from https://github.com/spf13/cobra)"
	matches := projectRegex.FindStringSubmatch(s)
	assert.Len(t, matches, 3)
}
