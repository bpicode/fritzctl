package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/manifest"
)

func parseManifest(filename string) *manifest.Plan {
	file, err := os.Open(filename)
	assert.NoError(err, "cannot open manifest file:", err)
	defer file.Close()
	p, err := manifest.Parse(file)
	assert.NoError(err, "cannot parse manifest file:", err)
	return p
}

func obtainSourcePlan(api fritz.HomeAutomationAPI) *manifest.Plan {
	l, err := api.ListDevices()
	assert.NoError(err, "cannot obtain device data:", err)
	return manifest.ConvertDevicelist(l)
}
