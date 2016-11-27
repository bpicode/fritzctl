package files

import (
	"errors"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHomeDir tests that the running user's homedir can be resolved.
func TestHomeDir(t *testing.T) {
	f, err := InHomeDir("abc.txt")
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

// TestErrorInDirOfUser tests that the error propagation.
func TestErrorInDirOfUser(t *testing.T) {
	currUser = func() (*user.User, error) {
		return nil, errors.New("some error")
	}
	_, err := InHomeDir("xyz.txt")
	assert.Error(t, err)
}
