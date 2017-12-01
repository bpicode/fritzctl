package cmd

import (
	"fmt"
	"os"

	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

var listSwitchesCmd = &cobra.Command{
	Use:     "switches",
	Short:   "List the available smart home switches",
	Long:    "List the available smart home devices [switches] and associated data.",
	Example: "fritzctl list switches",
	RunE:    listSwitches,
}

func init() {
	listCmd.AddCommand(listSwitchesCmd)
}

func listSwitches(_ *cobra.Command, _ []string) error {
	c := homeAutoClient()
	devs, err := c.List()
	assertNoErr(err, "cannot obtain data for smart home switches")
	logger.Success("Device data:")

	table := switchTable()
	appendSwitches(devs, table)
	table.Print(os.Stdout)
	return nil
}

func switchTable() *console.Table {
	return console.NewTable(console.Headers(
		"NAME",
		"PRODUCT",
		"PRESENT",
		"STATE",
		"LOCK (BOX/DEV)",
		"MODE",
		"POWER [W]",
		"ENERGY [Wh]",
		"TEMP [°C]",
		"OFFSET [°C]",
	))
}

func appendSwitches(devs *fritz.Devicelist, table *console.Table) {
	for _, dev := range devs.Switches() {
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
		dev.Powermeter.FmtPowerW(),
		dev.Powermeter.FmtEnergyWh(),
		dev.Temperature.FmtCelsius(),
		dev.Temperature.FmtOffset(),
	}
}
