package cmd

import (
	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/spf13/cobra"
)

var toggleCmd = &cobra.Command{
	Use:     "toggle mydevice",
	Short:   "toggle on/off state of a device",
	Long:    "Change the on/off state of a device to the opposite of what it had before. Has no effect on devices that fo not support toggling.",
	Example: "fritzctl toggle dev1 dev2 dev3",
	RunE:    toogle,
}

func init() {
	RootCmd.AddCommand(toggleCmd)
}

func toogle(cmd *cobra.Command, args []string) error {
	aha := fritz.HomeAutomation(clientLogin())
	err := fritz.ConcurrentHomeAutomation(aha).Toggle(args...)
	assert.NoError(err, "error toggling device(s):", err)
	return nil
}
