package cliapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliCreate(t *testing.T) {
	cli := Create()
	assert.NotNil(t, cli)
}

func TestCommandsHaveHelp(t *testing.T) {
	c := Create()
	for _, command := range c.Commands {
		com, _ := command()
		help := com.Help()
		assert.NotEmpty(t, help)
	}
}

func TestCommandsHaveSynopsis(t *testing.T) {
	c := Create()
	for _, command := range c.Commands {
		com, _ := command()
		syn := com.Synopsis()
		assert.NotEmpty(t, syn)
	}
}

func TestCommandsHaveSaneCommandStrings(t *testing.T) {
	c := Create()
	for str, command := range c.Commands {
		com, _ := command()
		assert.NotNil(t, com)
		assert.NotEmpty(t, str)
		assert.NotContains(t, str, " ")
	}
}
