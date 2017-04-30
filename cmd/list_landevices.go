package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
)

type listLandevicesCommand struct {
}

func (cmd *listLandevicesCommand) Help() string {
	return "List the available LAN devices along with several information like IP addresses, MAC addresses, etc."
}

func (cmd *listLandevicesCommand) Synopsis() string {
	return "list the available LAN devices"
}

func (cmd *listLandevicesCommand) Run(args []string) int {
	c := clientLogin()
	f := fritz.Internal(c)
	devs, err := f.ListLanDevices()
	assert.NoError(err, "cannot obtain LAN devices data:", err)
	logger.Success("Obtained LAN devices data:")

	table := cmd.table()
	cmd.appendData(table, *devs)
	table.Render()
	return 0
}

func (cmd *listLandevicesCommand) table() *tablewriter.Table {
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

func (cmd *listLandevicesCommand) appendData(table *tablewriter.Table, devs fritz.LanDevices) {
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

// ListLandevices is a factory creating commands for listing LAN devices.
func ListLandevices() (cli.Command, error) {
	p := listLandevicesCommand{}
	return &p, nil
}
