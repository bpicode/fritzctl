package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

var temperatureCmd = &cobra.Command{
	Use:   "temperature [value in °C, on, off, sav, comf] [device/group names]",
	Short: "Set the temperature of HKR devices/groups or turn them on/off",
	Long: "Change the temperature of HKR devices/groups by supplying the desired value in °C. " +
		"When turning HKR devices on/off, replace the value by 'on'/'off' respectively." +
		"To reset each devices to its comfort/saving temperature, replace the value by 'comf'/'sav'." +
		"To increase/decrease temperatures relative to the current goal, supply '+' or '-' followed by space.",
	Example: `fritzctl temperature 21.0 HKR_1 HKR_2
fritzctl temperature off HKR_1
fritzctl temperature on HKR_2
fritzctl temperature comf HK1 HKR_2
fritzctl temperature sav HK1 HKR_2
fritzctl temperature + 1.5 HK1
fritzctl temperature - 2 HK1
`,
	RunE: changeTemperature,
}

func init() {
	RootCmd.AddCommand(temperatureCmd)
}

func changeTemperature(cmd *cobra.Command, args []string) error {
	assertMinLen(args, 2, "insufficient input: at least two parameters expected.\n\n", cmd.UsageString())
	val := args[0]
	action := changeAction(val)
	action(val, args[1:]...)
	logger.Info("It may take a few minutes until the changes propagate to the end device(s)")
	return nil
}

func changeAction(s string) func(val string, args ...string) {
	if strings.EqualFold(s, "sav") || strings.EqualFold(s, "saving") {
		return changeToSav
	}
	if strings.EqualFold(s, "comf") || strings.EqualFold(s, "comfort") {
		return changeToComf
	}
	if s == "+" || s == "-" {
		return changeBy
	}
	return changeTo
}

func changeToSav(_ string, args ...string) {
	changeByCallback(func(t fritz.Thermostat) string {
		return t.FmtComfortTemperature()
	}, args...)
}

func changeToComf(_ string, args ...string) {
	changeByCallback(func(t fritz.Thermostat) string {
		return t.FmtSavingTemperature()
	}, args...)
}

func changeBy(val string, args ...string) {
	assertMinLen(args, 2, "insufficient input: expected [+ or -] [amount] [devices]")
	delta, err := strconv.ParseFloat(val+args[0], 64)
	assertNoErr(err, "cannot parse temperature adjustment")
	changeByCallback(func(t fritz.Thermostat) string {
		cur, err := strconv.ParseFloat(t.FmtGoalTemperature(), 64)
		assertNoErr(err, "unable to parse the current temperature goal '%s'", t.FmtGoalTemperature())
		return strconv.FormatFloat(cur+delta, 'f', -1, 64)
	}, args[1:]...)
}

func changeTo(val string, devs ...string) {
	changeByValue(nil, val, devs...)
}

func changeByCallback(supplier func(t fritz.Thermostat) string, names ...string) {
	c := homeAutoClient(fritz.Caching(true))
	devices, err := c.List()
	assertNoErr(err, "cannot list available devices")
	for _, name := range names {
		device := deviceWithName(name, devices.Thermostats())
		assertTrue(device != nil, fmt.Sprintf("device with name '%s' not found", name))
		changeByValue(c, supplier(device.Thermostat), name)
	}
}

func deviceWithName(name string, list []fritz.Device) *fritz.Device {
	for _, d := range list {
		if d.Name == name {
			return &d
		}
	}
	return nil
}

func changeByValue(c fritz.HomeAuto, val string, names ...string) {
	temp, err := parseTemperature(val)
	assertNoErr(err, "cannot parse temperature value")
	if c == nil {
		c = homeAutoClient()
	}
	err = c.Temp(temp, names...)
	assertNoErr(err, "error setting temperature")
}

func parseTemperature(s string) (float64, error) {
	if strings.EqualFold(s, "off") {
		return 126.5, nil
	}
	if strings.EqualFold(s, "on") {
		return 127.0, nil
	}
	temp, errorParse := strconv.ParseFloat(s, 64)
	return temp, errorParse
}
