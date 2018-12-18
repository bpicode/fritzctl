package fritz

import "strings"

// Devicelist wraps a list of devices. This corresponds to
// the outer layer of the xml that the FRITZ!Box returns.
type Devicelist struct {
	Devices []Device `xml:"device"`
	Groups  []Group  `xml:"group"`
}

// DeviceGroup is an inflated version of one Group with multiple Device members.
type DeviceGroup struct {
	Group   Group
	Devices []Device
}

// Switches returns the devices which satisfy IsSwitch.
func (l *Devicelist) Switches() []Device {
	return l.filter(func(d Device) bool {
		return d.IsSwitch()
	})
}

// Thermostats returns the devices which satisfy IsThermostat.
func (l *Devicelist) Thermostats() []Device {
	return l.filter(func(d Device) bool {
		return d.IsThermostat()
	})
}

// AlertSensors returns the devices which satisfy HasAlertSensor.
func (l *Devicelist) AlertSensors() []Device {
	return l.filter(func(d Device) bool {
		return d.HasAlertSensor()
	})
}

// Buttons returns the devices which have a pressable button.
func (l *Devicelist) Buttons() []Device {
	return l.filter(func(d Device) bool {
		return d.Button.LastPressedTimestamp != ""
	})
}

func (l *Devicelist) filter(predicate func(Device) bool) []Device {
	var filtered []Device
	for _, d := range l.Devices {
		if predicate(d) {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

// DeviceGroups returns a slice of DeviceGroup by joining Group.Members() on Device.ID.
func (l *Devicelist) DeviceGroups() []DeviceGroup {
	var gs []DeviceGroup
	for _, g := range l.Groups {
		ds := l.devicesWithIDs(g.Members())
		gs = append(gs, DeviceGroup{Group: g, Devices: ds})
	}
	return gs
}

// DeviceWithID searches for a Device by its ID returns a the found/zero value and a flag true/false indicating
// whether the search was successful.
func (l *Devicelist) DeviceWithID(id string) (Device, bool) {
	for _, d := range l.Devices {
		if d.ID == id {
			return d, true
		}
	}
	return Device{}, false
}

// NamesAndAins returns a lookup name -> AIN.
func (l *Devicelist) NamesAndAins() map[string]string {
	ds := l.Devices
	gs := l.Groups
	table := make(map[string]string, len(ds)+len(gs))
	for _, grp := range gs {
		table[grp.Name] = strings.Replace(grp.Identifier, " ", "", -1)
	}
	for _, dev := range ds {
		table[dev.Name] = strings.Replace(dev.Identifier, " ", "", -1)
	}
	return table
}

func (l *Devicelist) devicesWithIDs(ids []string) []Device {
	var ds []Device
	for _, id := range ids {
		if d, ok := l.DeviceWithID(id); ok {
			ds = append(ds, d)
			continue
		}
	}
	return ds
}
