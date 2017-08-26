package fritz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParsingFunctionBitMask verifies correctness of the function bit.mask interpretation.
func TestParsingFunctionBitMask(t *testing.T) {
	assertions := assert.New(t)
	assertions.False((&Device{Functionbitmask: "nonsense"}).IsSwitch())
	assertions.False((&Device{Functionbitmask: "nonsense"}).IsThermostat())
	assertions.True((&Device{Functionbitmask: "2944"}).IsSwitch())
	assertions.False((&Device{Functionbitmask: "2944"}).IsThermostat())
	assertions.False((&Device{Functionbitmask: "320"}).IsSwitch())
	assertions.True((&Device{Functionbitmask: "320"}).IsThermostat())
}
