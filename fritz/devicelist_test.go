package fritz

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSwitchAndThermostatFiltering tests on the correctness of the switch/thermostat separation of a given device list.
func TestSwitchAndThermostatFiltering(t *testing.T) {
	l := mustUnmarshall(t, "../testdata/devicelist_fritzos06.83.xml")

	assert.Len(t, l.Thermostats(), 2)
	assert.Len(t, l.Switches(), 1)

	assert.Equal(t, len(l.Devices), len(l.Switches())+len(l.Thermostats()))
}

// TestSwitchAndThermostatFilteringIssue56 reproduces https://github.com/bpicode/fritzctl/issues/59.
func TestSwitchAndThermostatFilteringIssue56(t *testing.T) {
	l := mustUnmarshall(t, "../testdata/devicelist_issue_59.xml")

	assert.Len(t, l.Thermostats(), 4)
	assert.Len(t, l.Switches(), 8)

	assert.Equal(t, len(l.Devices), len(l.Switches())+len(l.Thermostats()))
}

// TestAlertSensorFilter tests filtering by alert sensor capabilities.
func TestAlertSensorFilter(t *testing.T) {
	l := mustUnmarshall(t, "../testdata/devicelist_fritzos06.83.xml")
	assert.Len(t, l.AlertSensors(), 0)
}

// TestGroupsIssue56 tests the group un-marshalling.
func TestGroupsIssue56(t *testing.T) {
	l := mustUnmarshall(t, "../testdata/devicelist_issue_59.xml")

	assertions := assert.New(t)

	groups := l.Groups
	assertions.Len(groups, 1)

	group := groups[0]
	assertions.True(group.MadeFromThermostats())
	assertions.False(group.MadeFromSwitches())
}

// TestGroupsSpec tests the group un-marshalling.
func TestGroupsSpec(t *testing.T) {
	l := mustUnmarshall(t, "../testdata/devicelist_from_spec.xml")

	groups := l.Groups

	assertions := assert.New(t)
	assertions.Len(groups, 1)

	group := groups[0]
	assertions.False(group.MadeFromThermostats())
	assertions.True(group.MadeFromSwitches())
}

// TestGroupMembersSpec tests the group joining.
func TestGroupMembersSpec(t *testing.T) {
	l := mustUnmarshall(t, "../testdata/devicelist_from_spec.xml")
	groups := l.DeviceGroups()
	assert.Len(t, groups, 1)
	group := groups[0]
	assert.Len(t, group.Devices, 1)
}

// TestGroupMembersSpec tests the group joining.
func TestGroupMembersIssue56(t *testing.T) {
	l := mustUnmarshall(t, "../testdata/devicelist_issue_59.xml")
	groups := l.DeviceGroups()
	assert.Len(t, groups, 1)
	group := groups[0]
	assert.Len(t, group.Devices, 4)
}

// TestDeviceWithId tests the search by ID.
func TestDeviceWithId(t *testing.T) {
	l := mustUnmarshall(t, "../testdata/devicelist_fritzos06.83.xml")

	d, ok := l.DeviceWithID("11")
	assert.True(t, ok)
	assert.Equal(t, d.ID, "11")

	_, ok = l.DeviceWithID("nonsense")
	assert.False(t, ok)
}

func mustUnmarshall(t *testing.T, fname string) Devicelist {
	assertions := assert.New(t)
	f, err := os.Open(fname)
	assertions.NoError(err)
	defer f.Close()
	var l Devicelist
	err = xml.NewDecoder(f).Decode(&l)
	assertions.NoError(err)
	return l
}
