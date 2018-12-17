package cmd

import (
	"os"

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

func export(_ *cobra.Command, _ []string) error {
	l := mustList()
	plan := manifest.ConvertDevicelist(l)
	manifest.ExporterTo(os.Stdout).Export(plan)
	return nil
}
