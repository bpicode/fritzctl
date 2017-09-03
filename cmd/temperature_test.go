package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseTemperature tests the user-supplied temperature interpretation.
func TestParseTemperature(t *testing.T) {
	assertions := assert.New(t)

	twelve, err := parseTemperature("12")
	assertions.NoError(err)
	assertions.Equal(float64(12), twelve)

	on, err := parseTemperature("on")
	assertions.NoError(err)
	assertions.Equal(float64(127), on)

	off, err := parseTemperature("off")
	assertions.NoError(err)
	assertions.Equal(float64(126.5), off)
}
