package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInfoLogging unit test.
func TestInfoLogging(t *testing.T) {
	assert.NotPanics(t, func() {
		Info("some", "random", "log")
	})
}

// TestSuccessLogging unit test.
func TestSuccessLogging(t *testing.T) {
	assert.NotPanics(t, func() {
		Success("This is a log")
		SuccessNoTimestamp("This is another log")
	})
}

// TestPanicLogging unit test.
func TestPanicLogging(t *testing.T) {
	assert.Panics(t, func() {
		Panic("I quit")
	})
}

// TestWarnLogging should not panic.
func TestWarnLogging(t *testing.T) {
	assert.NotPanics(t, func() {
		Warn("A warning")
	})
}
