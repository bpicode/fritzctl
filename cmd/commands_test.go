package cmd

import (
	"fmt"
	"net"
	"net/http/httptest"
	"testing"

	"github.com/bpicode/fritzctl/config"
	"github.com/bpicode/fritzctl/mock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestCommands is a unit test that runs most commands.
func TestCommands(t *testing.T) {
	oldPlaces := defaultConfigPlaces
	defer func() { defaultConfigPlaces = oldPlaces }()
	defaultConfigPlaces = append([]config.Place{config.InDir("../testdata/config", "config_localhost_http_test.json", config.JSON())}, defaultConfigPlaces...)

	testCases := []struct {
		cmd  *cobra.Command
		args []string
		srv  *httptest.Server
	}{
		{cmd: versionCmd, srv: mock.New().UnstartedServer()},
		{cmd: toggleCmd, args: []string{"SWITCH_3"}, srv: mock.New().UnstartedServer()},
		{cmd: temperatureCmd, args: []string{"19.5", "HKR_1"}, srv: mock.New().UnstartedServer()},
		{cmd: temperatureCmd, args: []string{"comf", "HKR_1"}, srv: mock.New().UnstartedServer()},
		{cmd: temperatureCmd, args: []string{"sav", "HKR_1"}, srv: mock.New().UnstartedServer()},
		{cmd: switchOnCmd, args: []string{"SWITCH_1"}, srv: mock.New().UnstartedServer()},
		{cmd: switchOffCmd, args: []string{"SWITCH_2"}, srv: mock.New().UnstartedServer()},
		{cmd: sessionIDCmd, srv: mock.New().UnstartedServer()},
		{cmd: pingCmd, srv: mock.New().UnstartedServer()},
		{cmd: planManifestCmd, args: []string{"../testdata/devicelist_fritzos06.83_plan.yml"}, srv: mock.New().UnstartedServer()},
		{cmd: exportManifestCmd, srv: mock.New().UnstartedServer()},
		{cmd: applyManifestCmd, args: []string{"../testdata/devicelist_fritzos06.83_plan.yml"}, srv: mock.New().UnstartedServer()},
		{cmd: listGroupsCmd, srv: mock.New().UnstartedServer()},
		{cmd: listLanDevicesCmd, srv: mock.New().UnstartedServer()},
		{cmd: listLogsCmd, srv: mock.New().UnstartedServer()},
		{cmd: listCallsCmd, srv: mock.New().UnstartedServer()},
		{cmd: listSwitchesCmd, srv: mock.New().UnstartedServer()},
		{cmd: listSwitchesCmd, args: []string{"--output=json"}, srv: mock.New().UnstartedServer()},
		{cmd: listThermostatsCmd, srv: mock.New().UnstartedServer()},
		{cmd: listThermostatsCmd, args: []string{"--output=json"}, srv: mock.New().UnstartedServer()},
		{cmd: docManCmd, srv: mock.New().UnstartedServer()},
		{cmd: boxInfoCmd, srv: mock.New().UnstartedServer()},
		{cmd: aboutCmd, srv: mock.New().UnstartedServer()},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Test run command %s", testCase.cmd.Name()), func(t *testing.T) {
			var err error
			testCase.srv.Listener, err = net.Listen("tcp", ":61666")
			assert.NoError(t, err)
			testCase.srv.Start()
			defer testCase.srv.Close()
			err = testCase.cmd.ParseFlags(testCase.args)
			assert.NoError(t, err)
			err = testCase.cmd.RunE(testCase.cmd, testCase.args)
			assert.NoError(t, err)
		})
	}
}

// TestCommandsHaveHelp ensures that every command provides
// a help text.
func TestCommandsHaveHelp(t *testing.T) {
	for _, c := range allCommands(RootCmd) {
		t.Run(fmt.Sprintf("test long description of command %s", c.Name()), func(t *testing.T) {
			assert.NotEmpty(t, c.Long)
		})
	}
}

// TestCommandsHaveUsage tests that command have a usage pattern.
func TestCommandsHaveUsage(t *testing.T) {
	for _, c := range allCommands(RootCmd) {
		t.Run(fmt.Sprintf("test usage term of command %s", c.Name()), func(t *testing.T) {
			assert.NotEmpty(t, c.Use)
		})
	}
}

// TestCommandsHaveSynopsis ensures that every command provides
// short a synopsis text.
func TestCommandsHaveSynopsis(t *testing.T) {
	for _, c := range allCommands(RootCmd) {
		t.Run(fmt.Sprintf("test short description of command %s", c.Name()), func(t *testing.T) {
			assert.NotEmpty(t, c.Short)
		})
	}
}

func allCommands(cmd *cobra.Command) []*cobra.Command {
	var commands []*cobra.Command
	commands = append(commands, cmd)
	for _, sub := range cmd.Commands() {
		commands = append(commands, allCommands(sub)...)
	}
	return commands
}
