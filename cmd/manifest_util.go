package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/manifest"
)

func parseManifest(filename string) *manifest.Plan {
	file, err := os.Open(filename)
	assertNoError(err, "cannot open manifest file:", err)
	defer file.Close()
	p, err := manifest.Parse(file)
	assertNoError(err, "cannot parse manifest file:", err)
	return p
}

func obtainSourcePlan(api fritz.HomeAutomationAPI) *manifest.Plan {
	l, err := api.ListDevices()
	assertNoError(err, "cannot obtain device data:", err)
	return manifest.ConvertDevicelist(l)
}
