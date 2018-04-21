package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_projectNameRegex(t *testing.T) {
	assert.True(t, projectNameRegex.MatchString("  name = \"github.com/spf13/pflag\""))
	assert.True(t, projectNameRegex.MatchString("name=a/b/c"))
	assert.False(t, projectNameRegex.MatchString("name="))
	assert.False(t, projectNameRegex.MatchString("name = "))
	assert.False(t, projectNameRegex.MatchString(`  analyzer-name = "dep"`))
	assert.False(t, projectNameRegex.MatchString(`  solver-name = "gps-cdcl"`))
}

func Test_parseName(t *testing.T) {
	dlp := depLockProjector{}
	assert.Equal(t, "github.com/spf13/pflag", dlp.parseName("  name = \"github.com/spf13/pflag\""))
}
