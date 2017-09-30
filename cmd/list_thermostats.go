package cmd

import (
	"os"
	"time"

	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/stringutils"
	"github.com/spf13/cobra"
)

var listThermostatsCmd = &cobra.Command{
	Use:     "thermostats",
	Short:   "List the available smart home thermostats",
	Long:    "List the available smart home devices [thermostats] and associated data.",
	Example: "fritzctl list thermostats",
	RunE:    listThermostats,
}

func init() {
	listCmd.AddCommand(listThermostatsCmd)
}

func listThermostats(cmd *cobra.Command, args []string) error {
	c := homeAutoClient()
	devs, err := c.List()
	assertNoErr(err, "cannot obtain thermostats device data")
	logger.Success("Device data:")

	table := thermostatsTable()
	appendThermostats(devs, table)
	table.Print(os.Stdout)
	return nil
}

var errorCodesVsDescriptions = map[string]string{
	"":  "",
	"0": "",
	"1": " Thermostat adjustment not possible. Is the device mounted correctly?",
	"2": " Valve plunger cannot be driven far enough. Possible solutions: Open and close the plunger a couple of times by hand. Check if the battery is too weak.",
	"3": " Valve plunger cannot be moved. Is it blocked?",
	"4": " Preparing installation.",
	"5": " Device in mode 'INSTALLATION'. It can be mounted now.",
	"6": " Device is adjusting to the valve plunger.",
}

func thermostatsTable() *console.Table {
	return console.NewTable(console.Headers(
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
	))
}

func appendThermostats(devs *fritz.Devicelist, table *console.Table) {
	for _, dev := range devs.Thermostats() {
		table.Append(thermostatColumns(dev))
	}
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
