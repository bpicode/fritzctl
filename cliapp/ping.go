package cliapp

import (
	"github.com/bpicode/fritzctl/logger"
	"github.com/mitchellh/cli"
)

type pingCommand struct {
}

func (cmd *pingCommand) Help() string {
	return "Attempt to contact the FRITZ!Box by trying to solve the login challenge"
}

func (cmd *pingCommand) Synopsis() string {
	return "Check if the FRITZ!Box responds"
}

func (cmd *pingCommand) Run(args []string) int {
	clientLogin()
	logger.Info("Success! FRITZ!Box seems to be alive!")
	return 0
}

func ping() (cli.Command, error) {
	p := pingCommand{}
	return &p, nil
}
