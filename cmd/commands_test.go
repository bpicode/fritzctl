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
	config.ConfigDir = "../testdata/config"
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
		{cmd: listLanDevicesCmd, args: []string{}, srv: mock.New().UnstartedServer()},
		{cmd: listLogsCmd, args: []string{}, srv: mock.New().UnstartedServer()},
		{cmd: listSwitchesCmd, srv: mock.New().UnstartedServer()},
		{cmd: listThermostatsCmd, srv: mock.New().UnstartedServer()},
		{cmd: docManCmd, srv: mock.New().UnstartedServer()},
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
	for i, c := range coreCommands() {
		t.Run(fmt.Sprintf("test short description of command %d", i), func(t *testing.T) {
			assert.NotEmpty(t, c.Short)
		})
	}
}

func allCommands() []*cobra.Command {
	all := []*cobra.Command{
		versionCmd,
		switchCmd,
		manifestCmd,
		listCmd,
		docCmd,
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
		listLanDevicesCmd,
		listLogsCmd,
		listSwitchesCmd,
		listThermostatsCmd,
		docManCmd,
	}
}
