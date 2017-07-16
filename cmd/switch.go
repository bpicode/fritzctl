package cmd

import (
	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use: "switch",
}

var switchOnCmd = &cobra.Command{
	Use:     "on [device names]",
	Short:   "Switch on device(s)",
	Long:    "Change the state of one ore more devices to \"on\".",
	Example: "fritzctl switch on SWITCH_1 SWITCH_2",
	RunE:    switchOn,
}

var switchOffCmd = &cobra.Command{
	Use:     "off [device names]",
	Short:   "Switch off device(s)",
	Long:    "Change the state of one ore more devices to \"off\".",
	Example: "fritzctl switch off SWITCH_1 SWITCH_2",
	RunE:    switchOff,
}

func init() {
	switchCmd.AddCommand(switchOnCmd, switchOffCmd)
	RootCmd.AddCommand(switchCmd)
}

func switchOn(cmd *cobra.Command, args []string) error {
	assert.StringSliceHasAtLeast(args, 1, "insufficient input: device name(s) expected")
	aha := fritz.HomeAutomation(clientLogin())
	err := fritz.ConcurrentHomeAutomation(aha).SwitchOn(args...)
	assert.NoError(err, "error switching on device(s):", err)
	return nil
}

func switchOff(cmd *cobra.Command, args []string) error {
	assert.StringSliceHasAtLeast(args, 1, "insufficient input: device name(s) expected.")
	aha := fritz.HomeAutomation(clientLogin())
	err := fritz.ConcurrentHomeAutomation(aha).SwitchOff(args...)
	assert.NoError(err, "error switching off device(s):", err)
	return nil
}
