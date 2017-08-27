package fritz

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
	var switches []Device
	for _, d := range l.Devices {
		if d.IsSwitch() {
			switches = append(switches, d)
		}
	}
	return switches
}

// Thermostats returns the devices which satisfy IsThermostat.
func (l *Devicelist) Thermostats() []Device {
	var thermostats []Device
	for _, d := range l.Devices {
		if d.IsThermostat() {
			thermostats = append(thermostats, d)
		}
	}
	return thermostats
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
