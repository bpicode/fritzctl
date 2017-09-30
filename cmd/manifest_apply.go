package cmd

import (
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/manifest"
	"github.com/spf13/cobra"
)

var applyManifestCmd = &cobra.Command{
	Use:     "apply [manifest file]",
	Short:   "Apply a given manifest",
	Long:    "Apply a given manifest against the state of the FRITZ!Box.",
	Example: "fritzctl manifest apply /path/to/manifest.yml",
	RunE:    apply,
}

func init() {
	manifestCmd.AddCommand(applyManifestCmd)
}

func apply(cmd *cobra.Command, args []string) error {
	assertStringSliceHasAtLeast(args, 1, "insufficient input: path to input manifest expected.")
	target := parseManifest(args[0])
	api := fritz.HomeAutomation(clientLogin())
	src := obtainSourcePlan(api)
	err := manifest.AhaAPIApplier(api).Apply(src, target)
	assertNoErr(err, "application of manifest was not successful")
	return nil
}
