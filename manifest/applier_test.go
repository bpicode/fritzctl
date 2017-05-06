package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDryRunSelfTransitionEmpty tests the dry-runner in the trivial sector.
func TestDryRunSelfTransitionEmpty(t *testing.T) {
	applier := DryRunner()
	err := applier.Apply(&Plan{}, &Plan{})
	assert.NoError(t, err)
}

// TestDryRunSwitchToggle tests the dry-runner for one switch on->off.
func TestDryRunSwitchToggle(t *testing.T) {
	applier := DryRunner()
	err := applier.Apply(&Plan{Switches: []Switch{{Name: "s", State: true}}}, &Plan{Switches: []Switch{{Name: "s", State: false}}})
	assert.NoError(t, err)
}

// TestDryRunSwitchNameNotFound tests the dry-runner for a switch that does not exist.
func TestDryRunSwitchNameNotFound(t *testing.T) {
	applier := DryRunner()
	err := applier.Apply(&Plan{Switches: []Switch{{Name: "s", State: true}}}, &Plan{Switches: []Switch{{Name: "x", State: false}}})
	assert.Error(t, err)
}

// TestDryRunSwitchToggleAndTemperatureChange tests the dry-runner.
func TestDryRunSwitchToggleAndTemperatureChange(t *testing.T) {
	applier := DryRunner()
	err := applier.Apply(
		&Plan{
			Switches:    []Switch{{Name: "s", State: true}},
			Thermostats: []Thermostat{{Name: "t", Temperature: 17.5}},
		},
		&Plan{
			Switches:    []Switch{{Name: "s", State: false}},
			Thermostats: []Thermostat{{Name: "t", Temperature: 20.5}},
		})
	assert.NoError(t, err)
}


// TestDryRunThermostatNameNotFound tests the dry-runner for a HKR that does not exist.
func TestDryRunThermostatNameNotFound(t *testing.T) {
	applier := DryRunner()
	err := applier.Apply(&Plan{Thermostats: []Thermostat{{Name: "XXX", Temperature: 24.5}}}, &Plan{Thermostats: []Thermostat{{Name: "YYY", Temperature: 20.5}}})
	assert.Error(t, err)
}