package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCompletionZshHasHelp ensures the command under test provides a help text.
func TestCompletionZshHasHelp(t *testing.T) {
	assert.NotEmpty(t, completionZshCmd.Long)
}

// TestCommandsHaveSynopsis ensures that the command under test provides short a synopsis text.
func TestCompletionZshHasSynopsis(t *testing.T) {
	assert.NotEmpty(t, completionZshCmd.Short)
}

// TestCompletionZsh tests the zsh completion export.
func TestCompletionZsh(t *testing.T) {
	err := completionZshCmd.RunE(completionBashCmd, nil)
	assert.NoError(t, err)
}
