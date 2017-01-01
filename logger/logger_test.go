package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegularLogDoNotPanic asserts that certain loggers do not panic.
func TestRegularLogDoNotPanic(t *testing.T) {
	loggers := []func(v ...interface{}){
		Info,
		Success,
		SuccessNoTimestamp,
		Warn,
	}
	for i, l := range loggers {
		t.Run(fmt.Sprintf("Test logger %d", i), func(t *testing.T) {
			assert.NotPanics(t, func() {
				l("some", "random", "log")
			})
		})
	}
}

// TestPanicLogging asserts that the panic logger panics.
func TestPanicLogging(t *testing.T) {
	assert.Panics(t, func() {
		Panic("I quit")
	})
}
