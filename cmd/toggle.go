package cmd

import (
	"github.com/spf13/cobra"
)

var toggleCmd = &cobra.Command{
	Use:     "toggle [device/group names]",
	Short:   "Toggle on/off state of device(s) or group(s) of devices",
	Long:    "Change the on/off state of device(s) or group(s) of devices to the opposite of what it had before. Has no effect on devices that do not support toggling.",
	Example: "fritzctl toggle dev1 dev2 dev3",
	RunE:    toggle,
}

func init() {
	RootCmd.AddCommand(toggleCmd)
}

func toggle(_ *cobra.Command, args []string) error {
	assertMinLen(args, 1, "insufficient input: device/group name(s) expected (run with --help for more details)")
	c := homeAutoClient()
	err := c.Toggle(args...)
	assertNoErr(err, "error toggling device(s)")
	return nil
}
