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

func switchOff(_ *cobra.Command, args []string) error {
	assertMinLen(args, 1, "insufficient input: device/group name(s) expected (run with --help for more details)")
	c := homeAutoClient()
	err := c.Off(args...)
	assertNoErr(err, "error switching off device(s)")
	return nil
}
