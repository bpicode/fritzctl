package fritz

import (
	"strings"

	"github.com/bpicode/fritzctl/internal/stringutils"
)

// Group models a grouping of smart home devices. This corresponds to
// the single entries of the xml that the FRITZ!Box returns.
// codebeat:disable[TOO_MANY_IVARS]
type Group struct {
	Identifier      string     `xml:"identifier,attr"`      // A unique ID like AIN, MAC address, etc.
	ID              string     `xml:"id,attr"`              // Internal device ID of the FRITZ!Box.
	Functionbitmask string     `xml:"functionbitmask,attr"` // Bitmask determining the functionality of the group: bit 6: Comet DECT, HKR, "thermostat", bit 7: energy measurment device, bit 8: temperature sensor, bit 9: switch, bit 10: AVM DECT repeater
	Fwversion       string     `xml:"fwversion,attr"`       // Firmware version, usually "1.0".
	Manufacturer    string     `xml:"manufacturer,attr"`    // Manufacturer, usually set to "AVM".
	Productname     string     `xml:"productname,attr"`     // Name of the product, usually set to "".
	Present         int        `xml:"present"`              // All devices connected (1) or not (0).
	Name            string     `xml:"name"`                 // The name of the group. Can be assigned in the web gui of the FRITZ!Box.
	Switch          Switch     `xml:"switch"`               // Only filled with sensible data when switches are contained in the group.
	Thermostat      Thermostat `xml:"hkr"`                  // Only filled with sensible data when thermostats are contained in the group.
	GroupInfo       GroupInfo  `xml:"groupinfo"`            // Core data. "What makes up the group".
}

// codebeat:enable[TOO_MANY_IVARS]

// GroupInfo contains the topological data of the grouping, in particular the members of the group.
type GroupInfo struct {
	MasterDeviceID string `xml:"masterdeviceid"` // Internal ID of the master-switch. "0" is no master is set.
	Members        string `xml:"members"`        // Internal IDs of the members of the group. Comma-separated values, references Device.ID.
}

// MadeFromSwitches returns true if the devices consist of switches.
func (g *Group) MadeFromSwitches() bool {
	return bitMasked{Functionbitmask: g.Functionbitmask}.hasMask(512)
}

// MadeFromThermostats returns true if the devices consist of thermostats.
func (g *Group) MadeFromThermostats() bool {
	return bitMasked{Functionbitmask: g.Functionbitmask}.hasMask(64)
}

// Members returns a slice of internal IDs corresponding to the device members that belong to this group.
func (g *Group) Members() []string {
	csv := strings.TrimSpace(g.GroupInfo.Members)
	tokens := strings.Split(csv, ",")
	return stringutils.Transform(tokens, strings.TrimSpace)
}
