package cmd

import (
	"os"
	"strings"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/internal/console"
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

var listLanDevicesCmd = &cobra.Command{
	Use:   "landevices",
	Short: "List the available LAN devices",
	Long:  "List the available LAN devices along with several information like IP addresses, MAC addresses, etc.",
	Example: `fritzctl list landevices
fritzctl list landevices --filters=active,online`,
	RunE: listLanDevices,
}

func init() {
	listLanDevicesCmd.Flags().StringP("filters", "", "", "filter device list")
	listCmd.AddCommand(listLanDevicesCmd)
}

func listLanDevices(cmd *cobra.Command, _ []string) error {
	c := clientLogin()
	f := fritz.NewInternal(c)
	devs, err := f.ListLanDevices()
	assertNoErr(err, "cannot obtain LAN devices data")
	devs = applyLanDevicesFilters(cmd.Flag("filters").Value.String(), *devs)
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

func applyLanDevicesFilters(filters string, devs fritz.LanDevices) *fritz.LanDevices {
	filteredDevices := devs
	filterIdentifiers := strings.Split(filters, ",")
	for _, filter := range filterIdentifiers {
		filteredDevices = applyLanDevicesFilter(filter, filteredDevices)
	}
	return &filteredDevices
}

func applyLanDevicesFilter(filter string, devs fritz.LanDevices) fritz.LanDevices {
	switch strings.TrimSpace(filter) {
	case "active":
		return filterActiveLanDevices(devs)
	case "online":
		return filterOnlineLanDevices(devs)
	default:
		return devs
	}
}

func filterActiveLanDevices(devs fritz.LanDevices) fritz.LanDevices {
	filtered := &fritz.LanDevices{}
	for _, device := range devs.Network {
		if device.Active == "1" {
			filtered.Network = append(filtered.Network, device)
		}
	}
	return *filtered
}

func filterOnlineLanDevices(devs fritz.LanDevices) fritz.LanDevices {
	filtered := &fritz.LanDevices{}
	for _, device := range devs.Network {
		if device.Online == "1" {
			filtered.Network = append(filtered.Network, device)
		}
	}
	return *filtered
}
