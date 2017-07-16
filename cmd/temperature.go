package cmd

import (
	"strconv"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/spf13/cobra"
)

var temperatureCmd = &cobra.Command{
	Use:     "temperature [value in °C] [device names]",
	Short:   "Set the temperature of HKR devices",
	Long:    "Change the temperature of one or more HKR devices.",
	Example: "fritzctl temperature 21.0 HKR_1 HKR_2",
	RunE:    changeTemperature,
}

func init() {
	RootCmd.AddCommand(temperatureCmd)
}

func changeTemperature(cmd *cobra.Command, args []string) error {
	assert.StringSliceHasAtLeast(args, 2, "insufficient input: at least two parameters expected.")
	temp, errorParse := strconv.ParseFloat(args[0], 64)
	assert.NoError(errorParse, "cannot parse temperature value:", errorParse)
	aha := fritz.HomeAutomation(clientLogin())
	err := fritz.ConcurrentHomeAutomation(aha).ApplyTemperature(temp, args[1:]...)
	assert.NoError(err, "error setting temperature:", err)
	return nil
}
