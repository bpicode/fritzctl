package cliapp

import (
	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/mitchellh/cli"
)

type switchOnCommand struct {
}

func (cmd *switchOnCommand) Help() string {
	return "Switch on device. Example usage: fritzctl switch on mydevice"
}

func (cmd *switchOnCommand) Synopsis() string {
	return "Switch on device"
}

func (cmd *switchOnCommand) Run(args []string) int {
	assert.StringSliceHasAtLeast(args, 1, "Insufficient input: device name expected")
	f := fritz.UsingClient(clientLogin())
	err := f.SwitchOn(args...)
	assert.NoError(err, "Error switching off device(s):", err)
	return 0
}

func switchOnDevice() (cli.Command, error) {
	p := switchOnCommand{}
	return &p, nil
}
