package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listLanDevicesCmd = &cobra.Command{
	Use:     "landevices",
	Short:   "List the available LAN devices",
	Long:    "List the available LAN devices along with several information like IP addresses, MAC addresses, etc.",
	Example: "fritzctl list landevices",
	RunE:    listLanDevices,
}

func init() {
	listCmd.AddCommand(listLanDevicesCmd)
}

func listLanDevices(cmd *cobra.Command, args []string) error {
	c := clientLogin()
	f := fritz.Internal(c)
	devs, err := f.ListLanDevices()
	assert.NoError(err, "cannot obtain LAN devices data:", err)
	logger.Success("Obtained LAN devices data:")

	table := table()
	appendData(table, *devs)
	table.Render()
	return nil
}

func table() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"NAME",
		"IP",
		"MAC",
		"ACT/ONL",
		"SPEED [Mbit/s]",
		"COMPLIANCE",
	})
	return table
}

func appendData(table *tablewriter.Table, devs fritz.LanDevices) {
	for _, dev := range devs.Network {
		table.Append([]string{
			dev.Name,
			dev.IP,
			dev.Mac,
			console.StringToCheckmark(dev.Active) + "/" + console.StringToCheckmark(dev.Online),
			dev.Speed,
			console.Stoc(dev.ParentalControlAbuse).Inverse().String(),
		})
	}
}
