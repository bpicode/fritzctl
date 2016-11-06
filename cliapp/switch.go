package cliapp

import (
	"strings"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
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
	res, err := f.SwitchOn(args[0])
	assert.NoError(err, "Unable to switch on device:", err)
	logger.Info("Success! FRITZ!Box answered:", strings.TrimSpace(res))
	return 0
}

func switchOnDevice() (cli.Command, error) {
	p := switchOnCommand{}
	return &p, nil
}

type switchOffCommand struct {
}

func (cmd *switchOffCommand) Help() string {
	return "Switch off device. Example usage: fritzctl switch on mydevice."
}

func (cmd *switchOffCommand) Synopsis() string {
	return "Switch off device"
}

func (cmd *switchOffCommand) Run(args []string) int {
	assert.StringSliceHasAtLeast(args, 1, "Insufficient input: device name expected.")
	f := fritz.UsingClient(clientLogin())
	res, err := f.SwitchOff(args[0])
	assert.NoError(err, "Unable to switch off device:", err)
	logger.Info("Success! FRITZ!Box answered:", strings.TrimSpace(res))
	return 0
}

func switchOffDevice() (cli.Command, error) {
	p := switchOffCommand{}
	return &p, nil
}
