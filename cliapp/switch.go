package cliapp

import (
	"github.com/bpicode/fritzctl/fatals"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/mitchellh/cli"
)

type switchCommand struct {
}

func (cmd *switchCommand) Help() string {
	return "Switch on/off device. Example usage: fritzctl switch on mydevice"
}

func (cmd *switchCommand) Synopsis() string {
	return "Switch on/off device."
}

func (cmd *switchCommand) Run(args []string) int {
	fatals.AssertStringSliceHasAtLeast(args, 2, "Insufficient input: two parameters expected.")
	f := fritz.UsingClient(clientLogin())
	res, err := f.Switch(args[1], args[0])
	fatals.AssertNoError(err, "Unable to switch device:", err)
	logger.Info("Success! FRITZ!Box answered: " + res)
	return 0
}

func switchDevice() (cli.Command, error) {
	p := switchCommand{}
	return &p, nil
}
