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
