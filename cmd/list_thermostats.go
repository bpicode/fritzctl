package cmd

import (
	"os"
	"time"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/stringutils"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
)

type listThermostatsCommand struct {
}

var errorCodesVsDescriptions = map[string]string{
	"":  "",
	"0": "",
	"1": " Thermostat adjustment not possible. Is the device mounted correctly?",
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
	f := fritz.HomeAutomation(c)
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
		"NEXT",
		"STATE",
		"BATTERY",
	})
	return table
}

func (cmd *listThermostatsCommand) appendDevices(devs *fritz.Devicelist, table *tablewriter.Table) *tablewriter.Table {
	for _, dev := range devs.Thermostats() {
		table.Append(thermostatColumns(dev))
	}
	return table
}
func thermostatColumns(dev fritz.Device) []string {
	var columnValues []string
	columnValues = appendMetadata(columnValues, dev)
	columnValues = appendRuntimeFlags(columnValues, dev)
	columnValues = appendTemperatureValues(columnValues, dev)
	columnValues = appendRuntimeWarnings(columnValues, dev)
	return columnValues
}

func appendMetadata(cols []string, dev fritz.Device) []string {
	return append(cols, dev.Name, dev.Manufacturer, dev.Productname)
}

func appendRuntimeFlags(cols []string, dev fritz.Device) []string {
	return append(cols,
		console.IntToCheckmark(dev.Present),
		console.StringToCheckmark(dev.Thermostat.Lock)+"/"+console.StringToCheckmark(dev.Thermostat.DeviceLock))
}

func appendRuntimeWarnings(cols []string, dev fritz.Device) []string {
	return append(cols, errorCode(dev.Thermostat.ErrorCode),
		console.Stoc(dev.Thermostat.BatteryLow).Inverse().String())
}

func appendTemperatureValues(cols []string, dev fritz.Device) []string {
	return append(cols,
		dev.Thermostat.FmtMeasuredTemperature(),
		dev.Temperature.FmtOffset(),
		dev.Thermostat.FmtGoalTemperature(),
		dev.Thermostat.FmtSavingTemperature(),
		dev.Thermostat.FmtComfortTemperature(),
		fmtNextChange(dev.Thermostat.NextChange))
}
func fmtNextChange(n fritz.NextChange) string {
	return stringutils.DefaultIfEmpty(n.FmtTimestamp(time.Now()), "?") +
		" -> " +
		stringutils.DefaultIfEmpty(n.FmtGoalTemperature(), "?") +
		"°C"
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
