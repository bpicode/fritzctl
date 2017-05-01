package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/math"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
)

type listSwitchesCommand struct {
}

func (cmd *listSwitchesCommand) Help() string {
	return "List the available smart home devices [switches] and associated data."
}

func (cmd *listSwitchesCommand) Synopsis() string {
	return "list the available smart home switches"
}

func (cmd *listSwitchesCommand) Run(args []string) int {
	c := clientLogin()
	f := fritz.HomeAutomation(c)
	devs, err := f.ListDevices()
	assert.NoError(err, "cannot obtain data for smart home switches:", err)
	logger.Success("Obtained device data:")

	table := cmd.table()
	table = cmd.appendDevices(devs, table)
	table.Render()
	return 0
}

func (cmd *listSwitchesCommand) table() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"NAME",
		"MANUFACTURER",
		"PRODUCTNAME",
		"PRESENT",
		"STATE",
		"LOCK (BOX/DEV)",
		"MODE",
		"POWER [W]",
		"ENERGY [Wh]",
		"TEMP [°C]",
		"OFFSET [°C]",
	})
	return table
}

func (cmd *listSwitchesCommand) appendDevices(devs *fritz.Devicelist, table *tablewriter.Table) *tablewriter.Table {
	for _, dev := range devs.Switches() {
		table.Append(switchColumns(dev))
	}
	return table
}
func switchColumns(dev fritz.Device) []string {
	return []string{
		dev.Name,
		dev.Manufacturer,
		dev.Productname,
		console.IntToCheckmark(dev.Present),
		console.StringToCheckmark(dev.Switch.State),
		console.StringToCheckmark(dev.Switch.Lock) + "/" + console.StringToCheckmark(dev.Switch.DeviceLock),
		dev.Switch.Mode,
		math.ParseFloatAndScale(dev.Powermeter.Power, 0.001),
		dev.Powermeter.Energy,
		math.ParseFloatAndScale(dev.Temperature.Celsius, 0.1),
		math.ParseFloatAndScale(dev.Temperature.Offset, 0.1),
	}
}

// ListSwitches is a factory creating commands for commands listing switches.
func ListSwitches() (cli.Command, error) {
	p := listSwitchesCommand{}
	return &p, nil
}
