package files

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHomeDir unit test.
func TestHomeDir(t *testing.T) {
	f, err := InHomeDir("abc.txt")
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

// TestError unit test.
func TestError(t *testing.T) {
	_, err := inDirOfUser("xyz.txt", nil, errors.New("some error"))
	assert.Error(t, err)
}
