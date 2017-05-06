package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/assert"
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
