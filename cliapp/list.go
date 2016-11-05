package cliapp

import (
	"os"
	"strings"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/math"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
)

type listCommand struct {
}

func (cmd *listCommand) Help() string {
	return "Lists the available smart home devices and associated data."
}

func (cmd *listCommand) Synopsis() string {
	return "Lists the available smart home devices"
}

func (cmd *listCommand) Run(args []string) int {
	c := clientLogin()
	f := fritz.UsingClient(c)
	devs, err := f.ListDevices()
	assert.NoError(err, "Cannot obtain device data:", err)
	logger.Info("Obtained device data:")

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
		"TEMPERATURE [°C]",
	})

	for _, dev := range devs.Devices {
		table.Append([]string{
			dev.Name,
			dev.Manufacturer,
			dev.Productname,
			checkMarkFromInt(dev.Present),
			checkMarkFromString(dev.Switch.State),
			checkMarkFromString(dev.Switch.Lock),
			dev.Switch.Mode,
			math.ParseFloatAndScale(dev.Powermeter.Power, 0.001),
			dev.Powermeter.Energy,
			math.ParseFloatAddAndScale(dev.Temperature.Celsius, dev.Temperature.Offset, 0.1),
		})
	}
	table.Render()
	return 0
}

func checkMarkFromInt(i int) string {
	if i == 0 {
		return logger.PanicSprintf("\u2718")
	}
	return logger.InfoSprintf("\u2714")
}

func checkMarkFromString(s string) string {
	str := strings.TrimSpace(s)
	if str == "" {
		return logger.WarnSprintf("?")
	} else if str == "0" {
		return logger.PanicSprintf("\u2718")
	} else {
		return logger.InfoSprintf("\u2714")
	}
}

func list() (cli.Command, error) {
	p := listCommand{}
	return &p, nil
}
