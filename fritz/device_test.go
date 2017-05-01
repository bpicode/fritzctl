package fritz

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSwitchAndThermostatFiltering tests on the correctness of the switch/thermostat separation of a given device list.
func TestSwitchAndThermostatFiltering(t *testing.T) {
	f, err := os.Open("../testdata/devicelist_fritzos06.83.xml")
	defer f.Close()
	assert.NoError(t, err)

	var l Devicelist
	err = xml.NewDecoder(f).Decode(&l)
	assert.NoError(t, err)

	assert.Len(t, l.Thermostats(), 2)
	assert.Len(t, l.Switches(), 1)

	assert.Equal(t, len(l.Devices), len(l.Switches()) + len(l.Thermostats()))
}