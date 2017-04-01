package cliapp

import (
	"flag"

	"github.com/bpicode/fritzctl/cmd"
	"github.com/bpicode/fritzctl/config"
	"github.com/mitchellh/cli"
)

// New creates a new CLI application, that provides
// the commands implemented within this cmd package.
func New() *cli.CLI {
	c := cli.NewCLI(config.ApplicationName, config.Version)
	c.Args = flag.Args()
	c.Commands = map[string]cli.CommandFactory{
		"configure":        cmd.Configure,
		"list landevices":  cmd.ListLandevices,
		"list logs":        cmd.ListLogs,
		"list switches":    cmd.ListSwitches,
		"list thermostats": cmd.ListThermostats,
		"list inetstats":   cmd.ListInetstats,
		"ping":             cmd.Ping,
		"sessionid":        cmd.SessionID,
		"switch on":        cmd.SwitchOnDevice,
		"switch off":       cmd.SwitchOffDevice,
		"toggle":           cmd.ToggleDevice,
		"temperature":      cmd.Temperature,
	}
	return c
}
