package cliapp

import "github.com/mitchellh/cli"
import "github.com/bpicode/fritzctl/logger"

type sessionIDCommand struct {
}

func (cmd *sessionIDCommand) Help() string {
	return "Obtain a session ID by solving the FRITZ!Box login challenge. The session ID can be used for subsequent requests until it gets invalidated."
}

func (cmd *sessionIDCommand) Synopsis() string {
	return "obtain a session ID"
}

func (cmd *sessionIDCommand) Run(args []string) int {
	client := clientLogin()
	logger.Success("Successfully obtained session ID: " + client.SessionInfo.SID)
	return 0
}

func sessionID() (cli.Command, error) {
	p := sessionIDCommand{}
	return &p, nil
}
