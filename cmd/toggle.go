package cmd

import (
	"github.com/spf13/cobra"
)

var toggleCmd = &cobra.Command{
	Use:     "toggle [device names]",
	Short:   "Toggle on/off state of a device",
	Long:    "Change the on/off state of a device to the opposite of what it had before. Has no effect on devices that fo not support toggling.",
	Example: "fritzctl toggle dev1 dev2 dev3",
	RunE:    toggle,
}

func init() {
	RootCmd.AddCommand(toggleCmd)
}

func toggle(cmd *cobra.Command, args []string) error {
	c := homeAutoClient()
	err := c.Toggle(args...)
	assertNoError(err, "error toggling device(s):", err)
	return nil
}
