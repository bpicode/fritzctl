package cliapp

import (
	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/mitchellh/cli"
)

type toggleCommand struct {
}

func (cmd *toggleCommand) Help() string {
	return "Toggle on/off state of a device. Example usage: fritzctl toggle mydevice"
}

func (cmd *toggleCommand) Synopsis() string {
	return "toggle on/off state of a device"
}

func (cmd *toggleCommand) Run(args []string) int {
	f := fritz.New(clientLogin())
	err := f.Toggle(args...)
	assert.NoError(err, "error toggling device(s):", err)
	return 0
}

func toggleDevice() (cli.Command, error) {
	p := toggleCommand{}
	return &p, nil
}
