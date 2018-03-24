package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/bpicode/fritzctl/cmd/jsonapi"
	"github.com/bpicode/fritzctl/cmd/printer"
	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

var listThermostatsCmd = &cobra.Command{
	Use:   "thermostats",
	Short: "List the available smart home thermostats",
	Long:  "List the available smart home devices [thermostats] and associated data.",
	Example: `fritzctl list thermostats
fritzctl list thermostats --output=json`,
	RunE: listThermostats,
}

func init() {
	listThermostatsCmd.Flags().StringP("output", "o", "", "specify output format")
	listCmd.AddCommand(listThermostatsCmd)
}

func listThermostats(cmd *cobra.Command, _ []string) error {
	c := homeAutoClient()
	devs, err := c.List()
	assertNoErr(err, "cannot obtain thermostats device data")
	data := remapThermostats(cmd, devs.Thermostats())
	logger.Success("Device data:")
	printer.Print(data, os.Stdout)
	return nil
}

func remapThermostats(cmd *cobra.Command, ds []fritz.Device) interface{} {
	switch cmd.Flag("output").Value.String() {
	case "json":
		return jsonapi.NewMapper().Convert(ds)
	default:
		return thermostatsTable(ds)
	}
}

func thermostatsTable(devs []fritz.Device) *console.Table {
	table := console.NewTable(console.Headers(
		"NAME",
		"PRODUCT",
		"PRESENT",
		"LOCK (BOX/DEV)",
		"MEASURED",
		"OFFSET",
		"WANT",
		"SAVING",
		"COMFORT",
		"NEXT",
		"STATE",
		"BATTERY",
	))
	appendThermostats(devs, table)
	return table
}

func appendThermostats(devs []fritz.Device, table *console.Table) {
	for _, dev := range devs {
		columns := thermostatColumns(dev)
		table.Append(columns)
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
	return append(cols, dev.Name, fmt.Sprintf("%s %s", dev.Manufacturer, dev.Productname))
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
		fmtUnit(dev.Thermostat.FmtMeasuredTemperature, "°C"),
		fmtUnit(dev.Temperature.FmtOffset, "°C"),
		fmtUnit(dev.Thermostat.FmtGoalTemperature, "°C"),
		fmtUnit(dev.Thermostat.FmtSavingTemperature, "°C"),
		fmtUnit(dev.Thermostat.FmtComfortTemperature, "°C"),
		fmtNextChange(dev.Thermostat.NextChange))
}

func fmtNextChange(n fritz.NextChange) string {
	ts := n.FmtTimestamp(time.Now())
	if ts == "" {
		return "?"
	}
	return ts + " -> " + fmtUnit(n.FmtGoalTemperature, "°C")
}

func errorCode(ec string) string {
	checkMark := console.Stoc(ec).Inverse()
	return checkMark.String() + fritz.HkrErrorDescriptions[ec]
}
