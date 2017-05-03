package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/manifest"
	"github.com/mitchellh/cli"
)

type manifestExportCommand struct {
}

func (cmd *manifestExportCommand) Help() string {
	return "Export the current state of the FRITZ!Box in manifest format and print it to stdout. Example usage: fritzctl --loglevel=error manifest export > current_state.yml"
}

func (cmd *manifestExportCommand) Synopsis() string {
	return "export the current state of the FRITZ!Box in manifest format"
}

func (cmd *manifestExportCommand) Run(args []string) int {
	c := clientLogin()
	f := fritz.HomeAutomation(c)
	l, err := f.ListDevices()
	assert.NoError(err, "cannot obtain thermostats device data:", err)
	plan := manifest.ConvertDevicelist(l)
	manifest.ExporterTo(os.Stdout).Export(plan)
	return 0
}

// ManifestExport is a factory creating commands for exporting manifest files.
func ManifestExport() (cli.Command, error) {
	p := manifestExportCommand{}
	return &p, nil
}
