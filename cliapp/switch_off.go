package cliapp

import (
	"strings"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/mitchellh/cli"
)

type switchOffCommand struct {
}

func (cmd *switchOffCommand) Help() string {
	return "Switch off device. Example usage: fritzctl switch off mydevice."
}

func (cmd *switchOffCommand) Synopsis() string {
	return "Switch off device"
}

func (cmd *switchOffCommand) Run(args []string) int {
	assert.StringSliceHasAtLeast(args, 1, "Insufficient input: device name expected.")
	f := fritz.UsingClient(clientLogin())
	res, err := f.SwitchOff(args[0])
	assert.NoError(err, "Unable to switch off device:", err)
	logger.Success("Success! FRITZ!Box answered:", strings.TrimSpace(res))
	return 0
}

func switchOffDevice() (cli.Command, error) {
	p := switchOffCommand{}
	return &p, nil
}
