package cmd

import (
	"github.com/spf13/cobra"
)

var switchOffCmd = &cobra.Command{
	Use:     "off [device names]",
	Short:   "Switch off device(s)",
	Long:    "Change the state of one ore more devices to \"off\".",
	Example: "fritzctl switch off SWITCH_1 SWITCH_2",
	RunE:    switchOff,
}

func init() {
	switchCmd.AddCommand(switchOffCmd)
}

func switchOff(cmd *cobra.Command, args []string) error {
	assertStringSliceHasAtLeast(args, 1, "insufficient input: device name(s) expected.")
	c := homeAutoClient()
	err := c.Off(args...)
	assertNoError(err, "error switching off device(s):", err)
	return nil
}
