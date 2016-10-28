package meta

import (
	"errors"
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
)

// TestVersion unit test.
func TestVersion(t *testing.T) {
	assert.NotNil(t, Version)
	assert.NotEmpty(t, Version)
	assert.NotContains(t, Version, " ")
}

// TestAppname unit test.
func TestAppname(t *testing.T) {
	assert.NotNil(t, ApplicationName)
	assert.NotEmpty(t, ApplicationName)
	assert.NotContains(t, ApplicationName, " ")
}

// TestConfigfilename unit test.
func TestConfigfilename(t *testing.T) {
	assert.NotNil(t, ConfigFilename)
	assert.NotEmpty(t, ConfigFilename)
	assert.NotContains(t, ConfigFilename, " ")
}

// TestConfigfile unit test.
func TestConfigfile(t *testing.T) {
	f, err := ConfigFile()
	assert.NoError(t, err)
	assert.NotNil(t, f)
	assert.NotEmpty(t, f)
}

// TestConfigfileWithSpecialDir unit test.
func TestConfigfileWithSpecialDir(t *testing.T) {
	ConfigDir = "./"
	f, err := ConfigFile()
	assert.NoError(t, err)
	assert.NotNil(t, f)
	assert.NotEmpty(t, f)
}

// TestFirstWithoutErrorButNoOkFunction unit test.
func TestFirstWithoutErrorButNoOkFunction(t *testing.T) {
	_, err := firstWithoutError(func() (string, error) {
		return "", errors.New("not ok")
	})
	assert.Error(t, err)
}

// TestCurry unit test.
func TestCurry(t *testing.T) {
	upper := curry("arg", func(arg string) (string, error) {
		return strings.ToUpper(arg), nil
	})
	assert.NotNil(t, upper)
	asUppercase, err := upper()
	assert.NoError(t, err)
	assert.Equal(t, "ARG", asUppercase)
}

// TestComposeWithError unit test.
func TestComposeWithError(t *testing.T) {
	composed := compose("arg", func(arg string) (string, error) {
		return "OK", nil
	}, func(arg string) (string, error) {
		return "Not ok", errors.New("an error")

	})
	assert.NotNil(t, composed)
	_, err := composed()
	assert.Error(t, err)
}
