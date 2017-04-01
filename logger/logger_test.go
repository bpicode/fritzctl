package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegularLogDoNotPanic asserts that certain loggers do not panicLog.
func TestRegularLogDoNotPanic(t *testing.T) {
	defer SetupLoggers()
	levels := []string{"debug", "info", "warn", "error", "none"}
	informationalLoggers := []func(v ...interface{}){
		Debug,
		Info,
		Success,
		Warn,
	}
	for _, level := range levels {
		configureLevelsByFlag(level)
		for i, l := range informationalLoggers {
			t.Run(fmt.Sprintf("Test logger %d on level %s", i, level), func(t *testing.T) {
				assert.NotPanics(t, func() {
					l("some", "random", "log")
				})
			})
		}
	}
}

// TestPanicLogging asserts that the Panic logger panics.
func TestPanicLogging(t *testing.T) {
	assert.Panics(t, func() {
		Panic("I quit")
	})
}

// TestApiUsage demonstrates simple usage.
func TestApiUsage(t *testing.T) {
	configureLevelsByFlag("debug")
	defer SetupLoggers()
	assert.NotPanics(t, func() {
		Info("One")
		Info("a", "message")
		Debug("debug", "message")
		Info()
	})

}
