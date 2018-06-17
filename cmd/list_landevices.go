package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/internal/console"
	"github.com/bpicode/fritzctl/logger"
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

func listLanDevices(_ *cobra.Command, _ []string) error {
	c := clientLogin()
	f := fritz.NewInternal(c)
	devs, err := f.ListLanDevices()
	assertNoErr(err, "cannot obtain LAN devices data")
	logger.Success("Obtained LAN devices data:")

	table := lanDevicesTable()
	appendData(table, *devs)
	table.Print(os.Stdout)
	return nil
}

func lanDevicesTable() *console.Table {
	return console.NewTable(console.Headers("NAME", "IP", "MAC", "ACT/ONL", "SPEED [Mbit/s]"))
}

func appendData(table *console.Table, devs fritz.LanDevices) {
	for _, dev := range devs.Network {
		table.Append([]string{
			dev.Name,
			dev.IP,
			dev.Mac,
			console.StringToCheckmark(dev.Active) + "/" + console.StringToCheckmark(dev.Online),
			dev.Speed,
		})
	}
}
