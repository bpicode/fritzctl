package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/manifest"
	"github.com/spf13/cobra"
)

var exportManifestCmd = &cobra.Command{
	Use:     "export",
	Short:   "Export the current state of the FRITZ!Box in manifest format",
	Long:    "Export the current state of the FRITZ!Box in manifest format and print it to stdout.",
	Example: "fritzctl --loglevel=error manifest export > current_state.yml",
	RunE:    export,
}

func init() {
	manifestCmd.AddCommand(exportManifestCmd)
}

func export(cmd *cobra.Command, args []string) error {
	c := clientLogin()
	f := fritz.HomeAutomation(c)
	l, err := f.ListDevices()
	assertNoErr(err, "cannot obtain device data")
	plan := manifest.ConvertDevicelist(l)
	manifest.ExporterTo(os.Stdout).Export(plan)
	return nil
}
