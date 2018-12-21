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

func apply(_ *cobra.Command, args []string) error {
	assertMinLen(args, 1, "insufficient input: path to input manifest expected")
	target := parseManifest(args[0])
	h := homeAutoClient(fritz.Caching(true))
	src := obtainSourcePlan(h)
	err := manifest.NewApplier(h).Apply(src, target)
	assertNoErr(err, "application of manifest was not successful")
	return nil
}
