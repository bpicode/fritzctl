package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInfoLogging unit test.
func TestInfoLogging(t *testing.T) {
	Info("This is a log")
	InfoNoTimestamp("This is another log")
}

// TestPanicLogging unit test.
func TestPanicLogging(t *testing.T) {
	assert.Panics(t, func() {
		Panic("I quit")
	})
}
