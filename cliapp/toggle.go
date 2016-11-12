package cliapp

import (
	"strings"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/mitchellh/cli"
)

type toggleCommand struct {
}

func (cmd *toggleCommand) Help() string {
	return "Toggle on/off state of a device. Example usage: fritzctl toggle mydevice"
}

func (cmd *toggleCommand) Synopsis() string {
	return "Toggle on/off state of a device"
}

func (cmd *toggleCommand) Run(args []string) int {
	f := fritz.UsingClient(clientLogin())
	res, err := f.Toggle(args[0])
	assert.NoError(err, "Unable to toggle device:", err)
	logger.Success("Success! FRITZ!Box answered:", strings.TrimSpace(res))
	return 0
}

func toggleDevice() (cli.Command, error) {
	p := toggleCommand{}
	return &p, nil
}
