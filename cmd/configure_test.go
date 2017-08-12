package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bpicode/fritzctl/config"
	"github.com/stretchr/testify/assert"
)

// TestConfigureHasHelp ensures the command under test provides a help text.
func TestConfigureHasHelp(t *testing.T) {
	assert.NotEmpty(t, configureCmd.Long)
}

// TestConfigureHasSynopsis ensures that the command under test provides short a synopsis text.
func TestConfigureHasSynopsis(t *testing.T) {
	assert.NotEmpty(t, configureCmd.Short)
}

// TestConfigure tests the interactive configuration.
func TestConfigure(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_fritzctl")
	defer os.Remove(tempDir)
	assert.NoError(t, err)
	config.DefaultDir = tempDir

	err = configureCmd.RunE(configureCmd, nil)
	assert.NoError(t, err)
}
