package cliapp

import "github.com/mitchellh/cli"
import "github.com/bpicode/fritzctl/configurer"
import "github.com/bpicode/fritzctl/assert"

type configureCommand struct {
}

func (cmd *configureCommand) Help() string {
	return "Walk through the configuration of fritzctl"
}

func (cmd *configureCommand) Synopsis() string {
	return "Configure fritzctl"
}

func (cmd *configureCommand) Run(args []string) int {
	cli := configurer.New()
	cli.ApplyDefaults(configurer.Defaults())
	cli.Obtain()
	err := cli.Write()
	assert.NoError(err, "Error writing configuration file:", err)
	return 0
}

func configure() (cli.Command, error) {
	p := configureCommand{}
	return &p, nil
}
