package cmd

import (
	"github.com/bpicode/fritzctl/logger"
	"github.com/mitchellh/cli"
)

type pingCommand struct {
}

func (cmd *pingCommand) Help() string {
	return "Attempt to contact the FRITZ!Box by trying to solve the login challenge."
}

func (cmd *pingCommand) Synopsis() string {
	return "check if the FRITZ!Box responds"
}

func (cmd *pingCommand) Run(args []string) int {
	clientLogin()
	logger.Success("Success! FRITZ!Box seems to be alive!")
	return 0
}

// Ping is a factory creating commands for FRITZ!Box ping.
func Ping() (cli.Command, error) {
	p := pingCommand{}
	return &p, nil
}
