package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTemperatureOf tests the temperatureOf method.
func TestTemperatureOf(t *testing.T) {
	plan, _ := ParseFile("../testdata/all_off.yml")
	assert.NotNil(t, plan)

	tmp, ok := plan.temperatureOf("ThermoOne")
	assert.Equal(t, float64(15), tmp)
	assert.Equal(t, true, ok)

	_, ok = plan.temperatureOf("DoesNotExist")
	assert.Equal(t, false, ok)

}

// TestSwitchStateOfOf tests the switchStateOf method.
func TestSwitchStateOfOf(t *testing.T) {
	plan, _ := ParseFile("../testdata/all_on.yml")
	assert.NotNil(t, plan)

	state, ok := plan.switchStateOf("SwitchOne")
	assert.Equal(t, true, state)
	assert.Equal(t, true, ok)

	_, ok = plan.switchStateOf("DoesNotExist")
	assert.Equal(t, false, ok)
}
