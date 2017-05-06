package manifest

// Plan represents the data model of an absolute state of the fritz smart home.
type Plan struct {
	Switches    []Switch     // The power switches.
	Thermostats []Thermostat // The HKR devices.
}

// Switch represents the state of a switch.
type Switch struct {
	Name  string // Name of the switch.
	State bool   // On (true) or off (false).
}

// Thermostat represents the state of a HKR device.
type Thermostat struct {
	Name        string  // Name of the device.
	Temperature float64 // The temperature in °C.
}

func (plan *Plan) switchNamed(name string) (sw Switch, ok bool) {
	for _, s := range plan.Switches {
		if name == s.Name {
			sw = s
			ok = true
			return sw, ok
		}
	}
	return sw, ok
}

func (plan *Plan) switchStateOf(name string) (bool, bool) {
	if sw, ok := plan.switchNamed(name); ok {
		return sw.State, true
	}
	return false, false
}

func (plan *Plan) thermostatNamed(name string) (th Thermostat, ok bool) {
	for _, t := range plan.Thermostats {
		if name == t.Name {
			th = t
			ok = true
			return th, ok
		}
	}
	return th, ok
}

func (plan *Plan) temperatureOf(name string) (float64, bool) {
	if th, ok := plan.thermostatNamed(name); ok {
		return th.Temperature, true
	}
	return 0, false
}
