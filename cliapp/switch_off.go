package cliapp

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
	return "Switch off device(s)"
}

func (cmd *switchOffCommand) Run(args []string) int {
	assert.StringSliceHasAtLeast(args, 1, "Insufficient input: device name(s) expected.")
	f := fritz.New(clientLogin())
	err := f.SwitchOff(args...)
	assert.NoError(err, "Error switching off device(s):", err)
	return 0
}

func switchOffDevice() (cli.Command, error) {
	p := switchOffCommand{}
	return &p, nil
}
