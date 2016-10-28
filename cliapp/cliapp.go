package cliapp

import (
	"os"

	"github.com/bpicode/fritzctl/meta"
	"github.com/mitchellh/cli"
)

// Create creates a new CLI application.
func Create() *cli.CLI {
	c := cli.NewCLI(meta.ApplicationName, meta.Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"ping":   ping,
		"list":   list,
		"switch": switchDevice,
		"toggle": toggleDevice,
	}
	return c
}
