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
		{cmd: &pingCommand{}, srv: mock.New().UnstartedServer()},
		{cmd: &listSwitchesCommand{}, srv: mock.New().UnstartedServer()},
		{cmd: &listThermostatsCommand{}, srv: mock.New().UnstartedServer()},
		{cmd: &switchOnCommand{}, args: []string{"SWITCH_1"}, srv: mock.New().UnstartedServer()},
		{cmd: &switchOffCommand{}, args: []string{"SWITCH_2"}, srv: mock.New().UnstartedServer()},
		{cmd: &temperatureCommand{}, args: []string{"19.5", "HKR_1"}, srv: mock.New().UnstartedServer()},
		{cmd: &toggleCommand{}, args: []string{"SWITCH_3"}, srv: mock.New().UnstartedServer()},
		{cmd: &sessionIDCommand{}, args: []string{}, srv: mock.New().UnstartedServer()},
		{cmd: &listLandevicesCommand{}, args: []string{}, srv: mock.New().UnstartedServer()},
		{cmd: &listLogsCommand{}, args: []string{}, srv: mock.New().UnstartedServer()},
		{cmd: &manifestExportCommand{}, srv: mock.New().UnstartedServer()},
		{cmd: &manifestPlanCommand{}, args: []string{"../testdata/devicelist_fritzos06.83_plan.yml"}, srv: mock.New().UnstartedServer()},
		{cmd: &manifestApplyCommand{}, args: []string{"../testdata/devicelist_fritzos06.83_plan.yml"}, srv: mock.New().UnstartedServer()},
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

// TestCommandsHaveHelp ensures that every command provides
// a help text.
func TestCommandsHaveHelp(t *testing.T) {
	c := cli.NewCLI(config.ApplicationName, config.Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"applymanifest":   ManifestApply,
		"exportmanifest":  ManifestExport,
		"planmanifest":    ManifestPlan,
		"listswitches":    ListSwitches,
		"listthermostats": ListThermostats,
		"listlandevices":  ListLandevices,
		"listlogs":        ListLogs,
		"ping":            Ping,
		"sessionid":       SessionID,
		"switchon":        SwitchOnDevice,
		"switchoff":       SwitchOffDevice,
		"toggle":          ToggleDevice,
		"temperature":     Temperature,
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
}

// TestCommandsHaveSynopsis ensures that every command provides
// short a synopsis text.
func TestCommandsHaveSynopsis(t *testing.T) {
	c := cli.NewCLI(config.ApplicationName, config.Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"applymanifest":   ManifestApply,
		"exportmanifest":  ManifestExport,
		"planmanifest":    ManifestPlan,
		"listswitches":    ListSwitches,
		"listthermostats": ListThermostats,
		"listlandevices":  ListLandevices,
		"listlogs":        ListLogs,
		"ping":            Ping,
		"sessionid":       SessionID,
		"switchon":        SwitchOnDevice,
		"switchoff":       SwitchOffDevice,
		"toggle":          ToggleDevice,
		"temperature":     Temperature,
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
}
