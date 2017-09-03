package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var temperatureCmd = &cobra.Command{
	Use:   "temperature [value in °C, on, off] [device/group names]",
	Short: "Set the temperature of HKR devices/groups or turn them on/off",
	Long: "Change the temperature of HKR devices/groups by supplying the desired value in °C. " +
		"When turning HKR devices on/off, replace the value by 'on'/'off' respectively.",
	Example: `fritzctl temperature 21.0 HKR_1 HKR_2
fritzctl temperature off HKR_1
fritzctl temperature on HKR_2
`,
	RunE: changeTemperature,
}

func init() {
	RootCmd.AddCommand(temperatureCmd)
}

func changeTemperature(cmd *cobra.Command, args []string) error {
	assertStringSliceHasAtLeast(args, 2, "insufficient input: at least two parameters expected.\n\n", cmd.UsageString())
	temp, errorParse := parseTemperature(args[0])
	assertNoError(errorParse, "cannot parse temperature value:", errorParse)
	c := homeAutoClient()
	err := c.Temp(temp, args[1:]...)
	assertNoError(err, "error setting temperature:", err)
	return nil
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
