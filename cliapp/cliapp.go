package cliapp

import (
	"os"

	"github.com/bpicode/fritzctl/meta"
	"github.com/mitchellh/cli"
)

// New creates a new CLI application, that provides
// the commands implemented within this package.
func New() *cli.CLI {
	c := cli.NewCLI(meta.ApplicationName, meta.Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"configure":   configure,
		"list":        delegating(pairOf("switches", listSwitches), pairOf("thermostats", listThermostats), pairOf("landevices", listLandevices)),
		"ping":        ping,
		"sessionid":   sessionID,
		"switch":      delegating(pairOf("on", switchOnDevice), pairOf("off", switchOffDevice)),
		"toggle":      toggleDevice,
		"temperature": temperature,
	}
	return c
}
