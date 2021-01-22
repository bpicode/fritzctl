package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/bpicode/fritzctl/cmd/printer"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/internal/console"
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
	listThermostatsCmd.Flags().BoolP("verbose", "v", false, "output all values")
	listCmd.AddCommand(listThermostatsCmd)
}

func listThermostats(cmd *cobra.Command, _ []string) error {
	devs := mustList()
	defaultF := func(devs []fritz.Device) interface{} {
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			verbose = false
		}
		return thermostatsTable(devs, verbose)
	}
	data := selectFmt(cmd, devs.Thermostats(), defaultF)
	logger.Success("Device data:")
	printer.Print(data, os.Stdout)
	return nil
}

func thermostatsTable(devs []fritz.Device, verbose bool) interface{} {
	headers := []string{
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
	}
	if verbose {
		headers = append(headers,
			"MODE (HOLIDAY/SUMMER)",
			"WINDOW (OPEN/UNTIL)",
			"BOOST (ACTIVE/UNTIL)",
		)
	}
	table := console.NewTable(console.Headers(headers...))
	appendThermostats(devs, table, verbose)
	return table
}

func appendThermostats(devs []fritz.Device, table *console.Table, verbose bool) {
	for _, dev := range devs {
		columns := thermostatColumns(dev, verbose)
		table.Append(columns)
	}
}

func thermostatColumns(dev fritz.Device, verbose bool) []string {
	var columnValues []string
	columnValues = appendMetadata(columnValues, dev)
	columnValues = appendRuntimeFlags(columnValues, dev)
	columnValues = appendTemperatureValues(columnValues, dev)
	columnValues = appendRuntimeWarnings(columnValues, dev)
	if verbose {
		columnValues = appendModeValues(columnValues, dev)
		columnValues = appendWindowValues(columnValues, dev)
		columnValues = appendBoostValues(columnValues, dev)
	}
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
	return append(cols, errorCode(dev.Thermostat.ErrorCode), batteryState(dev.Thermostat))
}

func appendModeValues(cols []string, dev fritz.Device) []string {
	return append(cols,
		fmt.Sprintf("%s/%s",
			console.Btoc(dev.Thermostat.Holiday).String(),
			console.Btoc(dev.Thermostat.Summer).String(),
		))
}

func appendWindowValues(cols []string, dev fritz.Device) []string {
	return append(cols,
		fmt.Sprintf("%s/%s",
			console.Stoc(dev.Thermostat.WindowOpen).String(),
			dev.Thermostat.FmtWindowOpenEndTimestamp(time.Now()),
		))
}

func appendBoostValues(cols []string, dev fritz.Device) []string {
	return append(cols,
		fmt.Sprintf("%s/%s",
			console.Btoc(dev.Thermostat.Boost).String(),
			dev.Thermostat.FmtBoostEndTimestamp(time.Now()),
		))
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

func batteryState(thermostat fritz.Thermostat) string {
	return fmt.Sprintf("%s%% %s", thermostat.BatteryChargeLevel, console.Stoc(thermostat.BatteryLow).Inverse().String())
}
