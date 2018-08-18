package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"

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
	assert.NoError(t, err)
	defer os.Remove(tempDir)
	defer func() { configReaderSrc = os.Stdin }()
	configReaderSrc = bytes.NewBufferString(path.Join(tempDir, "config.yml") + "\n")

	err = configureCmd.RunE(configureCmd, nil)
	assert.NoError(t, err)
}
