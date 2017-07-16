package cliapp

import (
	"github.com/bpicode/fritzctl/cmd"
	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/flags"
	"github.com/mitchellh/cli"
)

// New creates a new CLI application, that provides
// the commands implemented within this cmd package.
func New() *cli.CLI {
	c := cli.NewCLI(config.ApplicationName, config.Version)
	c.Args = flags.Args()
	c.Commands = map[string]cli.CommandFactory{
		"completion bash":  cmd.CompletionBash(c),
		"configure":        cmd.Configure,
	}
	return c
}
