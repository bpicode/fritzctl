package config

import (
	"errors"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestVersion asserts that a version is present.
func TestVersion(t *testing.T) {
	assert.NotNil(t, Version)
	assert.NotEmpty(t, Version)
	assert.NotContains(t, Version, " ")
}

// TestHomeDir tests that the running user's homedir can be resolved.
func TestHomeDir(t *testing.T) {
	home := homeDirOf(user.Current)
	f, err := home("abc.txt")
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

// TestErrorInDirOfUser tests that the error propagation.
func TestErrorInDirOfUser(t *testing.T) {
	home := homeDirOf(func() (*user.User, error) {
		return nil, errors.New("some error")
	})
	_, err := home("xyz.txt")
	assert.Error(t, err)
}
