package cmd

import (
	"fmt"
	"testing"

	"github.com/bpicode/fritzctl/config"
	"github.com/stretchr/testify/assert"
)

// TestConfigFileCannotBeDetermined asserts that default options are used if no config file can be found.
func TestConfigFileCannotBeDetermined(t *testing.T) {
	assertions := assert.New(t)
	opts := optsFromPlaces(config.InDir("", "asjnfasjfbq3.yml", config.YAML()))
	assertions.Empty(opts)
}

// TestConfigFiles walks through several config files and pipes them through the option determination.
func TestConfigFiles(t *testing.T) {
	for i, path := range []string{
		"../testdata/config/config_localhost_http_test.json",
		"../testdata/config/config_skip_tls.json",
		"../testdata/config/config_with_cert.json",
	} {
		t.Run(fmt.Sprintf("config file %d %s", i, path), func(t *testing.T) {
			assertions := assert.New(t)
			opts := optsFromPlaces(config.InDir("", path, config.JSON()))
			assertions.NotEmpty(opts)
		})
	}
}
