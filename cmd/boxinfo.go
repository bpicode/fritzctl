package cmd

import (
	"fmt"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/internal/console"
	"github.com/spf13/cobra"
)

var boxInfoCmd = &cobra.Command{
	Use:     "boxinfo",
	Short:   "Display information about the FRITZ!Box",
	Long:    "Show information about the FRITZ!Box like firmware version, uptime, etc.",
	Example: "fritzctl boxinfo",
	RunE:    boxInfo,
}

func init() {
	RootCmd.AddCommand(boxInfoCmd)
}

func boxInfo(_ *cobra.Command, _ []string) error {
	c := clientLogin()
	f := fritz.NewInternal(c)
	info, err := f.BoxInfo()
	assertNoErr(err, "cannot obtain FRITZ!Box data")
	printBoxInfo(info)
	return nil
}

func printBoxInfo(boxData *fritz.BoxData) {
	fmt.Printf("%s %s\n", console.Cyan("Model:    "), boxData.Model)
	fmt.Printf("%s %s\n", console.Cyan("Firmware: "), boxData.FirmwareVersion)
	fmt.Printf("%s %s\n", console.Cyan("Running:  "), boxData.Runtime)
}
