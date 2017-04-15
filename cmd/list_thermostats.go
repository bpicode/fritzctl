package cmd

import (
	"os"
	"strings"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/math"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
)

type listThermostatsCommand struct {
}

var errorCodesVsDescriptions = map[string]string{
	"":  "",
	"0": "",
	"1": " Thermostat adjustment not possible. Is the device mounted corretly?",
	"2": " Valve plunger cannot be driven far enough. Possilbe solutions: Open and close the plunger a couple of times by hand. Check if the battery is too weak.",
	"3": " Valve plunger cannot be moved. Is it blocked?",
	"4": " Preparing installation.",
	"5": " Device in mode 'INSTALLATION'. It can be mounted now.",
	"6": " Device is adjusting to the valve plunger.",
}

func (cmd *listThermostatsCommand) Help() string {
	return "List the available smart home devices [thermostats] and associated data."
}

func (cmd *listThermostatsCommand) Synopsis() string {
	return "list the available smart home thermostats"
}

func (cmd *listThermostatsCommand) Run(args []string) int {
	c := clientLogin()
	f := fritz.New(c)
	devs, err := f.ListDevices()
	assert.NoError(err, "cannot obtain thermostats device data:", err)
	logger.Success("Obtained device data:")

	table := cmd.table()
	table = cmd.appendDevices(devs, table)
	table.Render()
	return 0
}

func (cmd *listThermostatsCommand) table() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"NAME",
		"MANUFACTURER",
		"PRODUCTNAME",
		"PRESENT",
		"LOCK (BOX/DEV)",
		"MEASURED [°C]",
		"OFFSET [°C]",
		"WANT [°C]",
		"SAVING [°C]",
		"COMFORT [°C]",
		"STATE",
	})
	return table
}

func (cmd *listThermostatsCommand) appendDevices(devs *fritz.Devicelist, table *tablewriter.Table) *tablewriter.Table {
	for _, dev := range devs.Devices {
		if dev.Thermostat.Measured != "" || dev.Thermostat.Goal != "" || dev.Thermostat.Saving != "" || dev.Thermostat.Comfort != "" || strings.Contains(dev.Productname, "Comet DECT") {
			table.Append(thermostatColumns(dev))
		}
	}
	return table
}
func thermostatColumns(dev fritz.Device) []string {
	return []string{
		dev.Name,
		dev.Manufacturer,
		dev.Productname,
		console.IntToCheckmark(dev.Present),
		console.StringToCheckmark(dev.Thermostat.Lock) + "/" + console.StringToCheckmark(dev.Thermostat.DeviceLock),
		math.ParseFloatAndScale(dev.Thermostat.Measured, 0.5),
		math.ParseFloatAndScale(dev.Temperature.Offset, 0.1),
		math.ParseFloatAndScale(dev.Thermostat.Goal, 0.5),
		math.ParseFloatAndScale(dev.Thermostat.Saving, 0.5),
		math.ParseFloatAndScale(dev.Thermostat.Comfort, 0.5),
		errorCode(dev.Thermostat.ErrorCode),
	}
}

func errorCode(ec string) string {
	checkMark := console.Stoc(ec).Inverse()
	return checkMark.String() + errorCodesVsDescriptions[ec]
}

// ListThermostats is a factory creating commands for listing thermostats.
func ListThermostats() (cli.Command, error) {
	p := listThermostatsCommand{}
	return &p, nil
}
