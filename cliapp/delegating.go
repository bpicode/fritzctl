package cliapp

import (
	"strings"

	"github.com/bpicode/fritzctl/assert"
	"github.com/mitchellh/cli"
)

type pair struct {
	cmd     string
	factory cli.CommandFactory
}

func pairOf(s string, l cli.CommandFactory) pair {
	return pair{cmd: s, factory: l}
}

func delegating(pairs ...pair) cli.CommandFactory {
	factories := make(map[string]cli.CommandFactory, len(pairs))
	for _, p := range pairs {
		factories[p.cmd] = p.factory
	}
	return (&delegatingCommandFactory{commandFactories: factories}).command
}

type delegatingCommandFactory struct {
	commandFactories map[string]cli.CommandFactory
}

type delegatingCommand struct {
	commandFactories map[string]cli.CommandFactory
}

func (cmd *delegatingCommand) Help() string {
	return join("Available subcommands:\n", cmd.commandFactories, func(c cli.Command) string {
		return c.Help()
	}, "\n")
}

func (cmd *delegatingCommand) Synopsis() string {
	return join("available subcommands: ", cmd.commandFactories, func(c cli.Command) string {
		return c.Synopsis()
	}, "; ")
}

func (cmd *delegatingCommand) Run(args []string) int {
	assert.StringSliceHasAtLeast(args, 1, "Insufficient input: subcommand required")
	firstArg := args[0]
	subCmdFactory, ok := cmd.commandFactories[firstArg]
	assert.IsTrue(ok, "Cannot find subcommand", firstArg)
	subcmd, err := subCmdFactory()
	assert.NoError(err)
	return subcmd.Run(args[1:])
}

func (delegating *delegatingCommandFactory) command() (cli.Command, error) {
	p := delegatingCommand{commandFactories: delegating.commandFactories}
	return &p, nil
}

func join(pre string, m map[string]cli.CommandFactory, f func(cli.Command) string, lineSep string) string {
	joined := pre
	for k, v := range m {
		subCmd, _ := v()
		joined += k + ": " + f(subCmd) + lineSep
	}
	return strings.TrimSpace(joined)
}
