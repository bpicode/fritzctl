package cmd

import (
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

func listSwitches(cmd *cobra.Command, args []string) error {
	c := homeAutoClient()
	devs, err := c.List()
	assertNoError(err, "cannot obtain data for smart home switches:", err)
	logger.Success("Obtained device data:")

	table := switchTable()
	appendSwitches(devs, table)
	table.Print(os.Stdout)
	return nil
}

func switchTable() *console.Table {
	return console.NewTable(console.Headers(
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
		dev.Manufacturer,
		dev.Productname,
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
