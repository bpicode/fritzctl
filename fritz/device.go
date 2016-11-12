package fritz

// Devicelist wraps a list of devices. This corresponds to
// the outer layer of the xml that the FRITZ!Box returns.
type Devicelist struct {
	Devices []Device `xml:"device"`
}

// Device models a smart home device. This corresponds to
// the single entries of the xml that the FRITZ!Box returns.
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
