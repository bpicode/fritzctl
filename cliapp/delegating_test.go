package cliapp

import (
	"testing"

	"log"

	"bytes"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

type testCommand struct {
	cli.Command
	help             string
	synopsis         string
	countInvocations int
}

func (cmd *testCommand) Help() string {
	return cmd.help
}

func (cmd *testCommand) Synopsis() string {
	return cmd.synopsis
}

func (cmd *testCommand) Run(args []string) int {
	cmd.countInvocations++
	log.Println("invovation #", cmd.countInvocations, cmd.help)
	return cmd.countInvocations
}

// TestDelegatingCommandsRoute1 is a unit test
func TestDelegatingCommandsRoute1(t *testing.T) {
	cmd1 := &testCommand{help: "help on cmd 1", synopsis: "synopis on cmd 1", countInvocations: 0}
	cmd2 := &testCommand{help: "help on cmd 2", synopsis: "synopis on cmd 2", countInvocations: 0}

	c1 := setupTestCli(cmd1, cmd2)
	c1.Args = []string{"cmd", "x"}
	_, err := c1.Run()

	assert.NoError(t, err)
	assert.Equal(t, 1, cmd1.countInvocations)
	assert.Equal(t, 0, cmd2.countInvocations)
}

// TestDelegatingCommandsRoute2 is a unit test
func TestDelegatingCommandsRoute2(t *testing.T) {
	cmd1 := &testCommand{help: "help on cmd 1", synopsis: "synopis on cmd 1", countInvocations: 0}
	cmd2 := &testCommand{help: "help on cmd 2", synopsis: "synopis on cmd 2", countInvocations: 0}

	c := setupTestCli(cmd1, cmd2)
	c.Args = []string{"cmd", "y"}
	c.Run()

	assert.Equal(t, 0, cmd1.countInvocations)
	assert.Equal(t, 1, cmd2.countInvocations)

}

// TestDelegatingCommandsHelp is a unit test
func TestDelegatingCommandsHelp(t *testing.T) {
	cmd1 := &testCommand{help: "help on cmd 1", synopsis: "synopis on cmd 1", countInvocations: 0}
	cmd2 := &testCommand{help: "help on cmd 2", synopsis: "synopis on cmd 2", countInvocations: 0}

	c := setupTestCli(cmd1, cmd2)
	buf3 := bytes.NewBufferString("")
	c.HelpWriter = buf3
	c.Args = []string{"--help", "cmd"}
	c.Run()

	log.Println("Help captured:\n", buf3.String())
	assert.Contains(t, buf3.String(), cmd1.Help())
	assert.Contains(t, buf3.String(), cmd2.Help())
}

// TestDelegatingCommandsSynopsis is a unit test
func TestDelegatingCommandsSynopsis(t *testing.T) {
	cmd1 := &testCommand{help: "help on cmd 1", synopsis: "synopis on cmd 1", countInvocations: 0}
	cmd2 := &testCommand{help: "help on cmd 2", synopsis: "synopis on cmd 2", countInvocations: 0}

	c := setupTestCli(cmd1, cmd2)
	buf := bytes.NewBufferString("")
	c.HelpWriter = buf
	c.Args = []string{}
	c.Run()

	log.Println("Synopsis captured:\n", buf.String())
	assert.Contains(t, buf.String(), cmd1.Synopsis())
	assert.Contains(t, buf.String(), cmd2.Synopsis())

}

func setupTestCli(cmd1, cmd2 *testCommand) *cli.CLI {
	cmd1Factory := func() (cli.Command, error) {
		return cmd1, nil
	}
	cmd2Factory := func() (cli.Command, error) {
		return cmd2, nil
	}
	factory := delegating(pairOf("x", cmd1Factory), pairOf("y", cmd2Factory))
	c := cli.NewCLI("test cli app", "0.1.124")
	c.Commands = map[string]cli.CommandFactory{
		"cmd": factory,
	}
	return c
}
