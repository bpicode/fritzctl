package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/manifest"
	"github.com/mitchellh/cli"
)

type manifestPlanCommand struct {
}

func (cmd *manifestPlanCommand) Help() string {
	return "Plan/dry-run a given manifest against the state of the FRITZ!Box. No changes will be applied."
}

func (cmd *manifestPlanCommand) Synopsis() string {
	return "plan a given manifest (dry-run)."
}

func (cmd *manifestPlanCommand) Run(args []string) int {
	assert.StringSliceHasAtLeast(args, 1, "insufficient input: path to input manifest expected.")
	target := cmd.parseManifest(args[0])

	l, err := fritz.HomeAutomation(clientLogin()).ListDevices()
	assert.NoError(err, "cannot obtain device data:", err)
	src := manifest.ConvertDevicelist(l)

	err = manifest.DryRunner().Apply(src, target)
	assert.NoError(err, "plan (dry-run) of manifest was not successful:", err)
	return 0
}

// ManifestPlan is a factory creating commands for dry-running manifest files.
func ManifestPlan() (cli.Command, error) {
	p := manifestPlanCommand{}
	return &p, nil
}

func (cmd *manifestPlanCommand) parseManifest(filename string) *manifest.Plan {
	file, err := os.Open(filename)
	assert.NoError(err, "cannot open manifest file:", err)
	defer file.Close()
	p, err := manifest.Parse(file)
	assert.NoError(err, "cannot parse manifest file:", err)
	return p
}