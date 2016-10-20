package cliapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCliCreate unit test.
func TestCliCreate(t *testing.T) {
	cli := Create()
	assert.NotNil(t, cli)
}

// TestCliCreate unit test.
func TestCommandsHaveHelp(t *testing.T) {
	c := Create()
	for _, command := range c.Commands {
		com, _ := command()
		help := com.Help()
		assert.NotEmpty(t, help)
	}
}

// TestCommandsHaveSynopsis unit test.
func TestCommandsHaveSynopsis(t *testing.T) {
	c := Create()
	for _, command := range c.Commands {
		com, _ := command()
		syn := com.Synopsis()
		assert.NotEmpty(t, syn)
	}
}

// TestCommandsHaveSaneCommandStrings unit test.
func TestCommandsHaveSaneCommandStrings(t *testing.T) {
	c := Create()
	for str, command := range c.Commands {
		com, _ := command()
		assert.NotNil(t, com)
		assert.NotEmpty(t, str)
		assert.NotContains(t, str, " ")
	}
}
