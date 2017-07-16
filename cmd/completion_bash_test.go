package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCompletionBashHasHelp ensures the command under test provides a help text.
func TestCompletionBashHasHelp(t *testing.T) {
	assert.NotEmpty(t, completionBashCmd.Long)
}

// TestCommandsHaveSynopsis ensures that the command under test provides short a synopsis text.
func TestCompletionBashHasSynopsis(t *testing.T) {
	assert.NotEmpty(t, completionBashCmd.Short)
}

// TestCompletionBash tests the bash completion export.
func TestCompletionBash(t *testing.T) {
	err := completionBashCmd.RunE(completionBashCmd, nil)
	assert.NoError(t, err)
}
