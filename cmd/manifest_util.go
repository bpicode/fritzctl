package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/manifest"
)

func parseManifest(filename string) *manifest.Plan {
	file, err := os.Open(filename)
	assertNoErr(err, "cannot open manifest file '%s'", filename)
	defer file.Close()
	p, err := manifest.Parse(file)
	assertNoErr(err, "cannot parse manifest file '%s'", filename)
	return p
}

func obtainSourcePlan(h fritz.HomeAuto) *manifest.Plan {
	l, err := h.List()
	assertNoErr(err, "cannot obtain device data")
	return manifest.ConvertDevicelist(l)
}
