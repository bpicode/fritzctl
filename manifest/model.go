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
