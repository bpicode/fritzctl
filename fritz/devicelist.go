package fritz

// Devicelist wraps a list of devices. This corresponds to
// the outer layer of the xml that the FRITZ!Box returns.
type Devicelist struct {
	Devices []Device `xml:"device"`
	Groups  []Group  `xml:"group"`
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
