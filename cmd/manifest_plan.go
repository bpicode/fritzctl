package cmd

import (
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/manifest"
	"github.com/spf13/cobra"
)

var planManifestCmd = &cobra.Command{
	Use:     "plan [manifest file]",
	Short:   "Plan a given manifest (dry-run)",
	Long:    "Plan/dry-run a given manifest against the state of the FRITZ!Box. No changes will be applied.",
	Example: "fritzctl manifest plan /path/to/manifest.yml",
	RunE:    plan,
}

func init() {
	manifestCmd.AddCommand(planManifestCmd)
}

func plan(cmd *cobra.Command, args []string) error {
	assertStringSliceHasAtLeast(args, 1, "insufficient input: path to input manifest expected.")
	target := parseManifest(args[0])
	api := fritz.HomeAutomation(clientLogin())
	src := obtainSourcePlan(api)
	err := manifest.DryRunner().Apply(src, target)
	assertNoErr(err, "plan (dry-run) of manifest was not successful")
	return nil
}
