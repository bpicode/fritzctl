package cmd

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

type listThermostatsCommand struct {
}

func (cmd *listThermostatsCommand) Help() string {
	return "List the available smart home devices [thermostats] and associated data."
}

func (cmd *listThermostatsCommand) Synopsis() string {
	return "list the available smart home thermostats"
}

func (cmd *listThermostatsCommand) Run(args []string) int {
	c := clientLogin()
	f := fritz.New(c)
	devs, err := f.ListDevices()
	assert.NoError(err, "cannot obtain thermostats device data:", err)
	logger.Success("Obtained device data:")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"NAME",
		"MANUFACTURER",
		"PRODUCTNAME",
		"PRESENT",
		"MEASURED [°C]",
		"OFFSET [°C]",
		"WANT [°C]",
		"SAVING [°C]",
		"COMFORT [°C]",
	})

	for _, dev := range devs.Devices {
		if dev.Thermostat.Measured != "" || dev.Thermostat.Goal != "" || dev.Thermostat.Saving != "" || dev.Thermostat.Comfort != "" || strings.Contains(dev.Productname, "Comet DECT") {
			table.Append([]string{
				dev.Name,
				dev.Manufacturer,
				dev.Productname,
				console.IntToCheckmark(dev.Present),
				math.ParseFloatAndScale(dev.Thermostat.Measured, 0.5),
				math.ParseFloatAndScale(dev.Temperature.Offset, 0.1),
				math.ParseFloatAndScale(dev.Thermostat.Goal, 0.5),
				math.ParseFloatAndScale(dev.Thermostat.Saving, 0.5),
				math.ParseFloatAndScale(dev.Thermostat.Comfort, 0.5),
			})
		}
	}
	table.Render()
	return 0
}

// ListThermostats is a factory creating commands for listing thermostats.
func ListThermostats() (cli.Command, error) {
	p := listThermostatsCommand{}
	return &p, nil
}
