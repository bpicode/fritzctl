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

// TestCommandsHaveHelp ensures that every command provides
// a help text.
func TestCommandsHaveHelp(t *testing.T) {
	c := New()
	for i, command := range c.Commands {
		t.Run(fmt.Sprintf("Test help of command %s", i), func(t *testing.T) {
			com, err := command()
			assert.NoError(t, err)
			help := com.Help()
			fmt.Printf("Help on command %s: '%s'\n", i, help)
			assert.NotEmpty(t, help)
		})
	}
}

// TestCommandsHaveSynopsis ensures that every command provides
// short a synopsis text.
func TestCommandsHaveSynopsis(t *testing.T) {
	c := New()
	for i, command := range c.Commands {
		t.Run(fmt.Sprintf("Test synopsis of command %s", i), func(t *testing.T) {
			com, err := command()
			assert.NoError(t, err)
			syn := com.Synopsis()
			fmt.Printf("Synopsis on command '%s': '%s'\n", i, syn)
			assert.NotEmpty(t, syn)
		})
	}
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
			assert.NotContains(t, str, " ")
		})
	}
}
