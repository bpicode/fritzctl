package cliapp

import (
	"os"
	"strings"

	"github.com/bpicode/fritzctl/fatals"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/meta"
	"github.com/bpicode/fritzctl/units"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
)

func clientLogin() *fritz.Client {
	configFile, errConfigFile := meta.ConfigFile()
	fatals.AssertNoError(errConfigFile, "Unable to create FRITZ!Box client:", errConfigFile)
	fritzClient, errCreate := fritz.NewClient(configFile)
	fatals.AssertNoError(errCreate, "Unable to create FRITZ!Box client:", errCreate)
	fritzClient, errLogin := fritzClient.Login()
	fatals.AssertNoError(errLogin, "Unable to login:", errLogin)
	return fritzClient
}

type pingCommand struct {
}

func (cmd *pingCommand) Help() string {
	return "Attempt to contact the FRITZ!Box by trying to solve the login challenge"
}

func (cmd *pingCommand) Synopsis() string {
	return "Check if the FRITZ!Box responds"
}

func (cmd *pingCommand) Run(args []string) int {
	clientLogin()
	logger.Info("Success! FRITZ!Box seems to be alive!")
	return 0
}

func ping() (cli.Command, error) {
	p := pingCommand{}
	return &p, nil
}

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
	fatals.AssertNoError(err, "Cannot obtain device data:", err)
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
			units.ParseFloatAndScale(dev.Powermeter.Power, 0.001),
			dev.Powermeter.Energy,
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

type switchCommand struct {
}

func (cmd *switchCommand) Help() string {
	return "Switch on/off device. Example usage: fritzctl switch on mydevice"
}

func (cmd *switchCommand) Synopsis() string {
	return "Switch on/off device."
}

func (cmd *switchCommand) Run(args []string) int {
	f := fritz.UsingClient(clientLogin())
	res, err := f.Switch(args[1], args[0])
	fatals.AssertNoError(err, "Unable to switch device:", err)
	logger.Info("Success! FRITZ!Box answered: " + res)
	return 0
}

func switchDevice() (cli.Command, error) {
	p := switchCommand{}
	return &p, nil
}
