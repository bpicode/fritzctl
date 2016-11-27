package cliapp

import (
	"os"
	"strings"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/math"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
)

type listSwitchesCommand struct {
}

func (cmd *listSwitchesCommand) Help() string {
	return "Lists the available smart home devices [switches] and associated data."
}

func (cmd *listSwitchesCommand) Synopsis() string {
	return "lists the available smart home switches"
}

func (cmd *listSwitchesCommand) Run(args []string) int {
	c := clientLogin()
	f := fritz.New(c)
	devs, err := f.ListDevices()
	assert.NoError(err, "cannot obtain data for smart home switches:", err)
	logger.Success("Obtained device data:")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"NAME",
		"MANUFACTURER",
		"PRODUCTNAME",
		"PRESENT",
		"STATE",
		"LOCK",
		"MODE",
		"POWER [W]",
		"ENERGY [Wh]",
		"TEMP [°C]",
	})

	for _, dev := range devs.Devices {
		if dev.Powermeter.Power != "" || dev.Powermeter.Energy != "" || strings.Contains(dev.Productname, "FRITZ!DECT") {
			table.Append([]string{
				dev.Name,
				dev.Manufacturer,
				dev.Productname,
				console.IntToCheckmark(dev.Present),
				console.StringToCheckmark(dev.Switch.State),
				console.StringToCheckmark(dev.Switch.Lock),
				dev.Switch.Mode,
				math.ParseFloatAndScale(dev.Powermeter.Power, 0.001),
				dev.Powermeter.Energy,
				math.ParseFloatAddAndScale(dev.Temperature.Celsius, dev.Temperature.Offset, 0.1),
			})
		}
	}
	table.Render()
	return 0
}

func listSwitches() (cli.Command, error) {
	p := listSwitchesCommand{}
	return &p, nil
}
