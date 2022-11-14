package cmd

import (
	"github.com/bpicode/fritzctl/cmd/printer"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/internal/console"
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
	"os"
)

var listSmarthomeCmd = &cobra.Command{
	Use:   "smarthome",
	Short: "List the available smart home devices",
	Long:  "List the available smart home devices and associated data.",
	Example: `fritzctl list smarthome
fritzctl list smarthome --output=json`,
	RunE: listSmarthome,
}

func init() {
	listSmarthomeCmd.Flags().StringP("output", "o", "", "specify output format")
	listCmd.AddCommand(listSmarthomeCmd)
}

func listSmarthome(cmd *cobra.Command, _ []string) error {
	devs := mustList()
	data := selectFmt(cmd, devs.Smarthome(), smarthomeTable)
	logger.Success("Device data:")
	printer.Print(data, os.Stdout)
	return nil
}

func smarthomeTable(devs []fritz.Device) interface{} {
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
		"HUMIDITY",
		"STATE",
		"BATTERY",
	))
	appendSmarthome(devs, table)
	return table
}

func appendSmarthome(devs []fritz.Device, table *console.Table) {
	for _, dev := range devs {
		columns := smarthomeColumns(dev)
		table.Append(columns)
	}
}

func smarthomeColumns(dev fritz.Device) []string {
	var columnValues []string
	columnValues = appendMetadata(columnValues, dev)
	columnValues = appendSmarthomeRuntimeFlags(columnValues, dev)
	columnValues = appendSmarthomeTemperatureValues(columnValues, dev)
	columnValues = appendSmarthomeHumidityValues(columnValues, dev)
	columnValues = appendSmarthomeRuntimeWarnings(columnValues, dev)
	return columnValues
}

func appendSmarthomeRuntimeFlags(cols []string, dev fritz.Device) []string {
	if dev.IsThermostat() {
		return append(cols,
			console.IntToCheckmark(dev.Present),
			console.StringToCheckmark(dev.Thermostat.Lock)+"/"+console.StringToCheckmark(dev.Thermostat.DeviceLock))
	} else {
		return append(cols,
			console.IntToCheckmark(dev.Present), "")
	}
}

func appendSmarthomeRuntimeWarnings(cols []string, dev fritz.Device) []string {
	if dev.IsThermostat() {
		return append(cols, errorCode(dev.Thermostat.ErrorCode), batteryState(dev.Thermostat))
	} else {
		return append(cols, "", "")
	}
}

func appendSmarthomeHumidityValues(cols []string, dev fritz.Device) []string {
	if dev.CanMeasureHumidity() {
		return append(cols, fmtUnit(dev.Humidity.FmtRelativeHumidity, "%"))
	} else {
		return append(cols, "")
	}
}

func appendSmarthomeTemperatureValues(cols []string, dev fritz.Device) []string {
	var measured func() string
	var nextChange string
	if dev.IsThermostat() {
		measured = dev.Thermostat.FmtMeasuredTemperature
		nextChange = fmtNextChange(dev.Thermostat.NextChange)
	} else {
		measured = dev.Temperature.FmtCelsius
	}

	return append(cols,
		fmtUnit(measured, "°C"),
		fmtUnit(dev.Temperature.FmtOffset, "°C"),
		fmtUnit(dev.Thermostat.FmtGoalTemperature, "°C"),
		fmtUnit(dev.Thermostat.FmtSavingTemperature, "°C"),
		fmtUnit(dev.Thermostat.FmtComfortTemperature, "°C"),
		nextChange)
}
