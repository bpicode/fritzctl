package cmd

import (
	"os"
	"testing"

	"github.com/bpicode/fritzctl/config"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

// TestCompletionBashHasHelp ensures the command under test provides a help text.
func TestCompletionBashHasHelp(t *testing.T) {
	commandFactory := CompletionBash(cli.NewCLI(config.ApplicationName, config.Version))
	command, err := commandFactory()
	assert.NoError(t, err)
	help := command.Help()
	assert.NotEmpty(t, help)
}

// TestCommandsHaveSynopsis ensures that the command under test provides short a synopsis text.
func TestCompletionBashHasSynopsis(t *testing.T) {
	commandFactory := CompletionBash(cli.NewCLI(config.ApplicationName, config.Version))
	command, err := commandFactory()
	assert.NoError(t, err)
	syn := command.Synopsis()
	assert.NotEmpty(t, syn)
}

// TestCompletionBash tests the bash completion export.
func TestCompletionBash(t *testing.T) {
	c := cli.NewCLI(config.ApplicationName, config.Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"complete": CompletionBash(c),
	}
	completionBashFactory := CompletionBash(c)
	command, err := completionBashFactory()
	assert.NoError(t, err)
	exitCode := command.Run([]string{})
	assert.Equal(t, 0, exitCode)
}
