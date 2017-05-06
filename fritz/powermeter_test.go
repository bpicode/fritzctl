package fritz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormattingOfEnergy tests formatting of values obtained by AHA interface.
func TestFormattingOfEnergy(t *testing.T) {
	assert.Equal(t, "2113", (&Powermeter{Energy: "2113"}).FmtEnergyWh())
}

// TestFormattingOfPower tests formatting of values obtained by AHA interface.
func TestFormattingOfPower(t *testing.T) {
	assert.Equal(t, "7", (&Powermeter{Power: "7000"}).FmtPowerW())
	assert.Zero(t, (&Powermeter{Power: "*kijeih14"}).FmtPowerW())
}
