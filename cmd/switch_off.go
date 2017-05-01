package cmd

import (
	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/mitchellh/cli"
)

type switchOffCommand struct {
}

func (cmd *switchOffCommand) Help() string {
	return "Switch off device(s). Example usage: fritzctl switch off mydevice."
}

func (cmd *switchOffCommand) Synopsis() string {
	return "switch off device(s)"
}

func (cmd *switchOffCommand) Run(args []string) int {
	assert.StringSliceHasAtLeast(args, 1, "insufficient input: device name(s) expected.")
	f := fritz.HomeAutomation(clientLogin())
	err := f.ConcurrentSwitchOff(args...)
	assert.NoError(err, "error switching off device(s):", err)
	return 0
}

// SwitchOffDevice is a factory creating commands for switching off switches.
func SwitchOffDevice() (cli.Command, error) {
	p := switchOffCommand{}
	return &p, nil
}
