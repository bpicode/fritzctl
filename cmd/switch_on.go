package cmd

import (
	"github.com/spf13/cobra"
)

var switchOnCmd = &cobra.Command{
	Use:   "on [device/group names]",
	Short: "Switch on devices or groups of devices",
	Long:  "Change the state of devices/groups to \"on\".",
	Example: `fritzctl switch on SWITCH_1 SWITCH_2
fritzctl switch on GROUP_1`,
	RunE: switchOn,
}

func init() {
	switchCmd.AddCommand(switchOnCmd)
}

func switchOn(_ *cobra.Command, args []string) error {
	assertMinLen(args, 1, "insufficient input: device/group name(s) expected (run with --help for more details)")
	c := homeAutoClient()
	err := c.On(args...)
	assertNoErr(err, "error switching on device(s)")
	return nil
}
