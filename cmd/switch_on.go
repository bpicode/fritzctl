package cmd

import (
	"github.com/spf13/cobra"
)

var switchOnCmd = &cobra.Command{
	Use:     "on [device names]",
	Short:   "Switch on device(s)",
	Long:    "Change the state of one ore more devices to \"on\".",
	Example: "fritzctl switch on SWITCH_1 SWITCH_2",
	RunE:    switchOn,
}

func init() {
	switchCmd.AddCommand(switchOnCmd)
}

func switchOn(cmd *cobra.Command, args []string) error {
	assertStringSliceHasAtLeast(args, 1, "insufficient input: device name(s) expected")
	c := homeAutoClient()
	err := c.On(args...)
	assertNoError(err, "error switching on device(s):", err)
	return nil
}
