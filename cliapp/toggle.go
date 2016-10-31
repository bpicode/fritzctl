package cliapp

import (
	"github.com/bpicode/fritzctl/fatals"
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
	return "Toggle on/off state of a device."
}

func (cmd *toggleCommand) Run(args []string) int {
	f := fritz.UsingClient(clientLogin())
	res, err := f.Toggle(args[0])
	fatals.AssertNoError(err, "Unable to toggle device:", err)
	logger.Info("Success! FRITZ!Box answered: " + res)
	return 0
}

func toggleDevice() (cli.Command, error) {
	p := toggleCommand{}
	return &p, nil
}
