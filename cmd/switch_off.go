package cmd

import (
	"github.com/spf13/cobra"
)

var switchOffCmd = &cobra.Command{
	Use:   "off [device/group names]",
	Short: "Switch off devices or groups of devices",
	Long:  "Change the state of devices/groups to \"off\".",
	Example: `fritzctl switch off SWITCH_1 SWITCH_2
fritzctl switch off GROUP_1`,
	RunE: switchOff,
}

func init() {
	switchCmd.AddCommand(switchOffCmd)
}

func switchOff(cmd *cobra.Command, args []string) error {
	assertStringSliceHasAtLeast(args, 1, "insufficient input: device/group name(s) expected.")
	c := homeAutoClient()
	err := c.Off(args...)
	assertNoError(err, "error switching off:", err)
	return nil
}
