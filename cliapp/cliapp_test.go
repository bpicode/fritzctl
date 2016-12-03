package cliapp

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

// TestCliCreate tests that the creation of a cli returns
// a sensible object.
func TestCliCreate(t *testing.T) {
	cli := New()
	assert.NotNil(t, cli)
	assert.NotNil(t, cli.Commands)
	assert.NotNil(t, cli.HelpFunc)
	assert.NotNil(t, cli.Name)
	assert.NotNil(t, cli.Version)
}

// TestCommandsHaveSaneCommandStrings ensures that the command
// string are sane (not empty, etc.).
func TestCommandsHaveSaneCommandStrings(t *testing.T) {
	c := New()
	for str, command := range c.Commands {
		t.Run(fmt.Sprintf("Test command string of command %s", str), func(t *testing.T) {
			com, err := command()
			assert.NoError(t, err)
			assert.NotNil(t, com)
			assert.NotEmpty(t, str)
		})
	}
}
