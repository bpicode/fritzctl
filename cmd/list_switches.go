package cmd

import (
	"fmt"
	"os"

	"github.com/bpicode/fritzctl/cmd/jsonapi"
	"github.com/bpicode/fritzctl/cmd/printer"
	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

var listSwitchesCmd = &cobra.Command{
	Use:   "switches",
	Short: "List the available smart home switches",
	Long:  "List the available smart home devices [switches] and associated data.",
	Example: `fritzctl list switches
fritzctl list switches --output=json`,
	RunE: listSwitches,
}

func init() {
	listSwitchesCmd.Flags().StringP("output", "o", "", "specify output format")
	listCmd.AddCommand(listSwitchesCmd)
}

func listSwitches(cmd *cobra.Command, _ []string) error {
	c := homeAutoClient()
	devs, err := c.List()
	assertNoErr(err, "cannot obtain data for smart home switches")
	logger.Success("Device data:")
	data := remapSwitches(cmd, devs.Switches())
	printer.Print(data, os.Stdout)
	return nil
}

func remapSwitches(cmd *cobra.Command, ds []fritz.Device) interface{} {
	switch cmd.Flag("output").Value.String() {
	case "json":
		return jsonapi.NewMapper().Convert(ds)
	default:
		return switchTable(ds)
	}
}

func switchTable(devs []fritz.Device) *console.Table {
	table := console.NewTable(console.Headers(
		"NAME",
		"PRODUCT",
		"PRESENT",
		"STATE",
		"LOCK (BOX/DEV)",
		"MODE",
		"POWER",
		"ENERGY",
		"TEMP",
		"OFFSET",
	))
	appendSwitches(devs, table)
	return table
}

func appendSwitches(devs []fritz.Device, table *console.Table) {
	for _, dev := range devs {
		table.Append(switchColumns(dev))
	}
}

func switchColumns(dev fritz.Device) []string {
	return []string{
		dev.Name,
		fmt.Sprintf("%s %s", dev.Manufacturer, dev.Productname),
		console.IntToCheckmark(dev.Present),
		console.StringToCheckmark(dev.Switch.State),
		console.StringToCheckmark(dev.Switch.Lock) + "/" + console.StringToCheckmark(dev.Switch.DeviceLock),
		dev.Switch.Mode,
		fmtUnit(dev.Powermeter.FmtPowerW, "W"),
		fmtUnit(dev.Powermeter.FmtEnergyWh, "Wh"),
		fmtUnit(dev.Temperature.FmtCelsius, "°C"),
		fmtUnit(dev.Temperature.FmtOffset, "°C"),
	}
}
