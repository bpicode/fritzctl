package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegularLogDoNotPanic asserts that certain loggers do not panicLog.
func TestRegularLogDoNotPanic(t *testing.T) {
	informationalLPrinters := []*parameterizableLogger{
		red(),
		yellow(),
		green(),
		plain(),
		dark(),
	}
	for i, p := range informationalLPrinters {
		t.Run(fmt.Sprintf("Test printer %d", i), func(t *testing.T) {
			assert.NotPanics(t, func() {
				p.print("some", "random", "log")
			})
		})
	}
}

// TestPanicLogging asserts that the Panic logger panics.
func TestPanicLogging(t *testing.T) {
	assert.Panics(t, func() {
		Panic("I quit")
	})
}
