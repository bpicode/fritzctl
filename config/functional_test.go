package config

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFirstWithoutErrorButNoOkFunction tests the situation of one function that returns an error.
func TestFirstWithoutErrorButNoOkFunction(t *testing.T) {
	_, err := firstWithoutError(func() (string, error) {
		return "", errors.New("not ok")
	})
	assert.Error(t, err)
}

// TestFirstWithoutErrorWithOkFunction tests the situation of two functions, both succeeding.
func TestFirstWithoutErrorWithOkFunction(t *testing.T) {
	_, err := firstWithoutError(func() (string, error) {
		return "", errors.New("not ok")
	}, func() (string, error) {
		return "ok", nil
	})
	assert.NoError(t, err)
}

// TestFirstWithoutErrorOneFunctionOk tests the situation of one function that succeeds.
func TestFirstWithoutErrorOneFunctionOk(t *testing.T) {
	_, err := firstWithoutError(func() (string, error) {
		return "XX", nil
	})
	assert.NoError(t, err)
}

// TestFirstWithoutErrorOneFunctionOkOneFunctionNotOk tests the situation of one function that succeeds and one that fails.
func TestFirstWithoutErrorOneFunctionOkOneFunctionNotOk(t *testing.T) {
	_, err := firstWithoutError(func() (string, error) {
		return "XX", nil
	}, func() (string, error) {
		return "", errors.New("ERR")
	})
	assert.NoError(t, err)
}

// TestFirstWithoutErrorOneFunctionNotOkOneFunctionOk tests the situation of one function that falis and one that succeeds.
func TestFirstWithoutErrorOneFunctionNotOkOneFunctionOk(t *testing.T) {
	_, err := firstWithoutError(func() (string, error) {
		return "", errors.New("ERR")
	}, func() (string, error) {
		return "XX", nil
	})
	assert.NoError(t, err)
}

// TestFirstWithoutErrorOneFunctionNotOkOneFunctionNotOk tests the situation of two functions, both failing.
func TestFirstWithoutErrorOneFunctionNotOkOneFunctionNotOk(t *testing.T) {
	_, err := firstWithoutError(func() (string, error) {
		return "", errors.New("ERR")
	}, func() (string, error) {
		return "XX", nil
	})
	assert.NoError(t, err)
}

// TestCurry tests the curry concept.
func TestCurry(t *testing.T) {
	upper := curry("arg", func(arg string) (string, error) {
		return strings.ToUpper(arg), nil
	})
	assert.NotNil(t, upper)
	asUppercase, err := upper()
	assert.NoError(t, err)
	assert.Equal(t, "ARG", asUppercase)
}

// TestCompose tests composition of two regular functions.
func TestCompose(t *testing.T) {
	composed := compose("arg", func(arg string) (string, error) {
		return "OK", nil
	}, func(arg string) (string, error) {
		return "OK", nil

	})
	assert.NotNil(t, composed)
	_, err := composed()
	assert.NoError(t, err)
}

// TestComposeWithError tests composition of one regular function and one that raises an error.
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
