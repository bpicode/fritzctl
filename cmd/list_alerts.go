package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/cmd/printer"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/internal/console"
	"github.com/bpicode/fritzctl/internal/stringutils"
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

var listAlertsCmd = &cobra.Command{
	Use:     "alerts",
	Short:   "List recent alert sensor reports",
	Long:    "List the most recently reported state of all devices with an alert sensor.",
	Example: "fritzctl list alerts",
	RunE:    listAlerts,
}

func init() {
	listAlertsCmd.Flags().StringP("output", "o", "", "specify output format")
	listCmd.AddCommand(listAlertsCmd)
}

func listAlerts(cmd *cobra.Command, _ []string) error {
	devs := mustList()
	data := selectFmt(cmd, devs.AlertSensors(), alertSensorsTable)
	logger.Success("Device data:")
	printer.Print(data, os.Stdout)
	return nil
}

func alertSensorsTable(devs []fritz.Device) interface{} {
	table := console.NewTable(console.Headers(
		"NAME",
		"ALERT",
	))
	symbols := stringutils.MapWithDefault(map[string]string{
		"0": console.Green("\u2714"),
		"1": console.Yellow("\u26a0"),
	}, console.Yellow("?"))
	for _, dev := range devs {
		columns := []string{dev.Name, symbols(dev.AlertSensor.State)}
		table.Append(columns)
	}
	return table
}
