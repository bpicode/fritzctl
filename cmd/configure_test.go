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
	command, err := Configure()
	assert.NoError(t, err)
	help := command.Help()
	assert.NotEmpty(t, help)
}

// TestConfigureHasSynopsis ensures that the command under test provides short a synopsis text.
func TestConfigureHasSynopsis(t *testing.T) {
	command, err := Configure()
	assert.NoError(t, err)
	syn := command.Synopsis()
	assert.NotEmpty(t, syn)
}

// TestConfigure tests the interactive configuration.
func TestConfigure(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_fritzctl")
	defer os.Remove(tempDir)
	assert.NoError(t, err)
	config.DefaultConfigDir = tempDir
	c := configureCommand{}
	i := c.Run([]string{})
	assert.Equal(t, 0, i)
}
