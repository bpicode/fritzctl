package cliapp

import (
	"strconv"
	"strings"

	"github.com/bpicode/fritzctl/fatals"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/meta"
	"github.com/mitchellh/cli"
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
	return "Lists the availble smart home devices and associated data."
}

func (cmd *listCommand) Synopsis() string {
	return "Lists the availble smart home devices"
}

func (cmd *listCommand) Run(args []string) int {
	c := clientLogin()
	f := fritz.UsingClient(c)
	devs, err := f.ListDevices()
	fatals.AssertNoError(err, "Cannot obtain device data:", err)
	logger.Info("Obtained device data:")
	logger.InfoNoTimestamp("+------------------+--------------+---------------------+---+---+---+---+")
	logger.InfoNoTimestamp("| NAME             | MANUFACTURER | PRODUCTNAME         | ? | S | M | L |")
	logger.InfoNoTimestamp("+------------------+--------------+---------------------+---+---+---+---+")
	for _, dev := range devs.Devices {
		logger.InfoNoTimestamp("| " + toSize(dev.Name, 16) +
			" | " + toSize(dev.Manufacturer, 12) +
			" | " + toSize(dev.Productname, 19) +
			" | " + toSize(strconv.Itoa(dev.Present), 1) +
			" | " + toSize(dev.Switch.State, 1) +
			" | " + toSize(dev.Switch.Mode, 1) +
			" | " + toSize(dev.Switch.Lock, 1) + " |")
	}
	logger.InfoNoTimestamp("+------------------+--------------+---------------------+---+---+---+---+")
	return 0
}

func toSize(s string, n int) string {
	str := strings.TrimSpace(s)
	if len(str) <= n {
		return str + strings.Repeat(" ", n-len(str))
	}
	return str[:n-3] + "..."
}

func list() (cli.Command, error) {
	p := listCommand{}
	return &p, nil
}
