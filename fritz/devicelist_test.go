package fritz

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDeviceFiltering tests the filter API of DeviceList.
func TestDeviceFiltering(t *testing.T) {
	for _, tc := range []struct {
		name    string
		xmlfile string
		filter  func(Devicelist) []Device
		expect  int
	}{
		{name: "two thermostats", xmlfile: "../testdata/devicelist_fritzos06.83.xml", filter: func(d Devicelist) []Device { return d.Thermostats() }, expect: 2},
		{name: "one switch", xmlfile: "../testdata/devicelist_fritzos06.83.xml", filter: func(d Devicelist) []Device { return d.Switches() }, expect: 1},
		{name: "on alert sensors", xmlfile: "../testdata/devicelist_fritzos06.83.xml", filter: func(d Devicelist) []Device { return d.AlertSensors() }, expect: 0},
		{name: "on buttons", xmlfile: "../testdata/devicelist_fritzos06.83.xml", filter: func(d Devicelist) []Device { return d.Buttons() }, expect: 0},
		{name: "four thermostats (Issue #56)", xmlfile: "../testdata/devicelist_issue_59.xml", filter: func(d Devicelist) []Device { return d.Thermostats() }, expect: 4},
		{name: "eight switches (Issue #56)", xmlfile: "../testdata/devicelist_issue_59.xml", filter: func(d Devicelist) []Device { return d.Switches() }, expect: 8},
	} {
		t.Run(tc.name, func(t *testing.T) {
			list := mustUnmarshall(t, tc.xmlfile)
			assert.Len(t, tc.filter(list), tc.expect)
		})
	}
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
