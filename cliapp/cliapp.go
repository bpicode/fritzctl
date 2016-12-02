package cliapp

import (
	"os"

	"github.com/bpicode/fritzctl/cmd"
	"github.com/bpicode/fritzctl/meta"
	"github.com/mitchellh/cli"
)

// New creates a new CLI application, that provides
// the commands implemented within this package.
func New() *cli.CLI {
	c := cli.NewCLI(meta.ApplicationName, meta.Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"configure":   cmd.Configure,
		"list":        delegating(pairOf("switches", cmd.ListSwitches), pairOf("thermostats", cmd.ListThermostats), pairOf("landevices", cmd.ListLandevices)),
		"ping":        cmd.Ping,
		"sessionid":   cmd.SessionID,
		"switch":      delegating(pairOf("on", cmd.SwitchOnDevice), pairOf("off", cmd.SwitchOffDevice)),
		"toggle":      cmd.ToggleDevice,
		"temperature": cmd.Temperature,
	}
	return c
}
