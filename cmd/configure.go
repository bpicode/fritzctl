package cmd

import (
	"strings"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/config"
	"github.com/mitchellh/cli"
)

type configureCommand struct {
}

func (cmd *configureCommand) Help() string {
	return strings.Join([]string{
		"Walk through the configuration of fritzctl interactively.",
		"Configuration file is saved at the end of the survey.",
		"Run fritzctl with administrator privileges if the configuration file cannot be saved by a normal user.",
	}, "\n")
}

func (cmd *configureCommand) Synopsis() string {
	return "configure fritzctl"
}

func (cmd *configureCommand) Run(args []string) int {
	configurer := config.NewConfigurer()
	configurer.ApplyDefaults(config.Defaults())
	configurer.Greet()
	configurer.Obtain()
	err := configurer.Write()
	assert.NoError(err, "error writing configuration file:", err)
	return 0
}

// Configure is a factory creating commands for interactive fritzctl configuration.
func Configure() (cli.Command, error) {
	p := configureCommand{}
	return &p, nil
}
