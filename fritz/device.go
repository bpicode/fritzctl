package fritz

// Device models a smart home device. This corresponds to
// the single entries of the xml that the FRITZ!Box returns.
// codebeat:disable[TOO_MANY_IVARS]
type Device struct {
	Identifier      string      `xml:"identifier,attr"`      // A unique ID like AIN, MAC address, etc.
	ID              string      `xml:"id,attr"`              // Internal device ID of the FRITZ!Box.
	Functionbitmask string      `xml:"functionbitmask,attr"` // Bitmask determining the functionality of the device: bit 6: Comet DECT, HKR, "thermostat", bit 7: energy measurment device, bit 8: temperature sensor, bit 9: switch, bit 10: AVM DECT repeater
	Fwversion       string      `xml:"fwversion,attr"`       // Firmware version of the device.
	Manufacturer    string      `xml:"manufacturer,attr"`    // Manufacturer of the device, usually set to "AVM".
	Productname     string      `xml:"productname,attr"`     // Name of the product, empty for unknown or undefined devices.
	Present         int         `xml:"present"`              // Device connected (1) or not (0).
	Name            string      `xml:"name"`                 // The name of the device. Can be assigned in the web gui of the FRITZ!Box.
	Switch          Switch      `xml:"switch"`               // Only filled with sensible data for switch devices.
	Powermeter      Powermeter  `xml:"powermeter"`           // Only filled with sensible data for devices with an energy actuator.
	Temperature     Temperature `xml:"temperature"`          // Only filled with sensible data for devices with a temperature sensor.
	Thermostat      Thermostat  `xml:"hkr"`                  // Thermostat data, only filled with sensible data for HKR devices.
}

// codebeat:enable[TOO_MANY_IVARS]

// IsSwitch returns true if the device is recognized to be a switch and returns false otherwise.
func (d *Device) IsSwitch() bool {
	return bitMasked{Functionbitmask: d.Functionbitmask}.hasMask(512)
}

// IsThermostat returns true if the device is recognized to be a HKR device and returns false otherwise.
func (d *Device) IsThermostat() bool {
	return bitMasked{Functionbitmask: d.Functionbitmask}.hasMask(64)
}

// Returns is the device has a temperature sensor
func (d *Device) HasTemperatureSensor() bool {
	return bitMasked{Functionbitmask: d.Functionbitmask}.hasMask(256)
}
