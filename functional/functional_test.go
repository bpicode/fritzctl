package functional

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFirstWithoutErrorButNoOkFunction unit test.
func TestFirstWithoutErrorButNoOkFunction(t *testing.T) {
	_, err := FirstWithoutError(func() (string, error) {
		return "", errors.New("not ok")
	})
	assert.Error(t, err)
}

// TestFirstWithoutErrorWithOkFunction unit test.
func TestFirstWithoutErrorWithOkFunction(t *testing.T) {
	_, err := FirstWithoutError(func() (string, error) {
		return "", errors.New("not ok")
	}, func() (string, error) {
		return "ok", nil
	})
	assert.NoError(t, err)
}

// TestCurry unit test.
func TestCurry(t *testing.T) {
	upper := Curry("arg", func(arg string) (string, error) {
		return strings.ToUpper(arg), nil
	})
	assert.NotNil(t, upper)
	asUppercase, err := upper()
	assert.NoError(t, err)
	assert.Equal(t, "ARG", asUppercase)
}

// TestCompose unit test.
func TestCompose(t *testing.T) {
	composed := Compose("arg", func(arg string) (string, error) {
		return "OK", nil
	}, func(arg string) (string, error) {
		return "OK", nil

	})
	assert.NotNil(t, composed)
	_, err := composed()
	assert.NoError(t, err)
}

// TestComposeWithError unit test.
func TestComposeWithError(t *testing.T) {
	composed := Compose("arg", func(arg string) (string, error) {
		return "OK", nil
	}, func(arg string) (string, error) {
		return "Not ok", errors.New("an error")

	})
	assert.NotNil(t, composed)
	_, err := composed()
	assert.Error(t, err)
}
