package cmd

import (
	"fmt"
	"net"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/mock"
	"github.com/mitchellh/cli"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestCommands is a unit test that runs most commands.
func TestCommands(t *testing.T) {

	config.ConfigDir = "../testdata"
	config.ConfigFilename = "config_localhost_https_test.json"

	testCases := []struct {
		cmd  cli.Command
		args []string
		srv  *httptest.Server
	}{
		{cmd: &listSwitchesCommand{}, srv: mock.New().UnstartedServer()},
		{cmd: &listThermostatsCommand{}, srv: mock.New().UnstartedServer()},
		{cmd: &listLandevicesCommand{}, args: []string{}, srv: mock.New().UnstartedServer()},
		{cmd: &listLogsCommand{}, args: []string{}, srv: mock.New().UnstartedServer()},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test run command %d", i), func(t *testing.T) {
			l, err := net.Listen("tcp", ":61666")
			assert.NoError(t, err)
			testCase.srv.Listener = l
			testCase.srv.Start()
			defer testCase.srv.Close()
			exitCode := testCase.cmd.Run(testCase.args)
			assert.Equal(t, 0, exitCode)
		})
	}
}

// TestCommandsCobra is a unit test that runs most commands.
func TestCommandsCobra(t *testing.T) {

	config.ConfigDir = "../testdata"
	config.ConfigFilename = "config_localhost_https_test.json"

	testCases := []struct {
		cmd  *cobra.Command
		args []string
		srv  *httptest.Server
	}{
		{cmd: versionCmd, srv: mock.New().UnstartedServer()},
		{cmd: toggleCmd, args: []string{"SWITCH_3"}, srv: mock.New().UnstartedServer()},
		{cmd: temperatureCmd, args: []string{"19.5", "HKR_1"}, srv: mock.New().UnstartedServer()},
		{cmd: switchOnCmd, args: []string{"SWITCH_1"}, srv: mock.New().UnstartedServer()},
		{cmd: switchOffCmd, args: []string{"SWITCH_2"}, srv: mock.New().UnstartedServer()},
		{cmd: sessionIDCmd, srv: mock.New().UnstartedServer()},
		{cmd: pingCmd, srv: mock.New().UnstartedServer()},
		{cmd: planManifestCmd, args: []string{"../testdata/devicelist_fritzos06.83_plan.yml"}, srv: mock.New().UnstartedServer()},
		{cmd: exportManifestCmd, srv: mock.New().UnstartedServer()},
		{cmd: applyManifestCmd, args: []string{"../testdata/devicelist_fritzos06.83_plan.yml"}, srv: mock.New().UnstartedServer()},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test run command %d", i), func(t *testing.T) {
			l, err := net.Listen("tcp", ":61666")
			assert.NoError(t, err)
			testCase.srv.Listener = l
			testCase.srv.Start()
			defer testCase.srv.Close()
			err = testCase.cmd.RunE(testCase.cmd, testCase.args)
			assert.NoError(t, err)
		})
	}
}

// TestCommandsHaveHelp ensures that every command provides
// a help text.
func TestCommandsHaveHelp(t *testing.T) {
	c := cli.NewCLI(config.ApplicationName, config.Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"listswitches":    ListSwitches,
		"listthermostats": ListThermostats,
		"listlandevices":  ListLandevices,
		"listlogs":        ListLogs,
	}
	for i, command := range c.Commands {
		t.Run(fmt.Sprintf("Test help of command %s", i), func(t *testing.T) {
			com, err := command()
			assert.NoError(t, err)
			help := com.Help()
			fmt.Printf("Help on command %s: '%s'\n", i, help)
			assert.NotEmpty(t, help)
		})
	}

	for i, c := range coreCommands() {
		t.Run(fmt.Sprintf("test long description of command %d", i), func(t *testing.T) {
			assert.NotEmpty(t, c.Long)
		})
	}
}

// TestCommandsHaveUsage tests that command have a usage pattern.
func TestCommandsHaveUsage(t *testing.T) {
	for i, c := range allCommands() {
		t.Run(fmt.Sprintf("test usage term of command %d", i), func(t *testing.T) {
			assert.NotEmpty(t, c.Use)
		})
	}
}

// TestCommandsHaveSynopsis ensures that every command provides
// short a synopsis text.
func TestCommandsHaveSynopsis(t *testing.T) {
	c := cli.NewCLI(config.ApplicationName, config.Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"listswitches":    ListSwitches,
		"listthermostats": ListThermostats,
		"listlandevices":  ListLandevices,
		"listlogs":        ListLogs,
	}
	for i, command := range c.Commands {
		t.Run(fmt.Sprintf("Test synopsis of command %s", i), func(t *testing.T) {
			com, err := command()
			assert.NoError(t, err)
			syn := com.Synopsis()
			fmt.Printf("Synopsis on command '%s': '%s'\n", i, syn)
			assert.NotEmpty(t, syn)
		})
	}

	for i, c := range coreCommands() {
		t.Run(fmt.Sprintf("test short description of command %s", i), func(t *testing.T) {
			assert.NotEmpty(t, c.Short)
		})
	}
}

func allCommands() []*cobra.Command {
	all := []*cobra.Command{
		versionCmd,
		switchCmd,
		manifestCmd,
	}
	core := coreCommands()
	all = append(all, core...)
	return all
}

func coreCommands() []*cobra.Command {
	return []*cobra.Command{
		versionCmd,
		toggleCmd,
		temperatureCmd,
		switchOnCmd,
		switchOffCmd,
		sessionIDCmd,
		pingCmd,
		planManifestCmd,
		exportManifestCmd,
		applyManifestCmd,
	}
}
