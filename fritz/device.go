package fritz

// Capability enumerates the device capabilities.
type Capability int

// Known (specified) device capabilities.
// see https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf section 3.2 for full list
const (
	HANFUNCompatibility Capability = iota
	_
	_
	_
	AlertTrigger
	AVMButton
	HeatControl
	PowerSensor
	TemperatureSensor
	StateSwitch
	DECTRepeater
	Microphone
	_
	HANFUNUnit
	_
	SwitchableDevice
	DimmableDevice
	ColorSettableDevice
	_
	_
	HumiditySensor
)

// Device models a smart home device. This corresponds to
// the single entries of the xml that the FRITZ!Box returns.
// codebeat:disable[TOO_MANY_IVARS]
type Device struct {
	Identifier      string      `xml:"identifier,attr"`      // A unique ID like AIN, MAC address, etc.
	ID              string      `xml:"id,attr"`              // Internal device ID of the FRITZ!Box.
	Functionbitmask string      `xml:"functionbitmask,attr"` // Bitmask determining the functionality of the device: bit 6: Comet DECT, HKR, "thermostat", bit 7: energy measurement device, bit 8: temperature sensor, bit 9: switch, bit 10: AVM DECT repeater
	Fwversion       string      `xml:"fwversion,attr"`       // Firmware version of the device.
	Manufacturer    string      `xml:"manufacturer,attr"`    // Manufacturer of the device, usually set to "AVM".
	Productname     string      `xml:"productname,attr"`     // Name of the product, empty for unknown or undefined devices.
	Present         int         `xml:"present"`              // Device connected (1) or not (0).
	Name            string      `xml:"name"`                 // The name of the device. Can be assigned in the web gui of the FRITZ!Box.
	Switch          Switch      `xml:"switch"`               // Only filled with sensible data for switch devices.
	Powermeter      Powermeter  `xml:"powermeter"`           // Only filled with sensible data for devices with an energy actuator.
	Temperature     Temperature `xml:"temperature"`          // Only filled with sensible data for devices with a temperature sensor.
	Humidity        Humidity    `xml:"humidity"`             // Only filled with sensible data for devices with a humidity sensor.
	Thermostat      Thermostat  `xml:"hkr"`                  // Thermostat data, only filled with sensible data for HKR devices.
	AlertSensor     AlertSensor `xml:"alert"`                // Only filled with sensible data for devices with an alert sensor.
	Button          Button      `xml:"button"`               // Button data, only filled with sensible data for button devices.
}

// codebeat:enable[TOO_MANY_IVARS]

// IsHANFUNCompatible returns true if the device speaks the "Home Area Network FUNctional protocol".
func (d *Device) IsHANFUNCompatible() bool {
	return d.Has(HANFUNCompatibility)
}

// HasAlertSensor returns true if the device has a sensor that may trigger alerts.
func (d *Device) HasAlertSensor() bool {
	return d.Has(AlertTrigger)
}

// IsAVMButton returns true if the device is an AVM button like the FRITZ!DECT 440 and returns false otherwise.
func (d *Device) IsAVMButton() bool {
	return d.Has(AVMButton)
}

// IsThermostat returns true if the device is recognized to be a HKR device and returns false otherwise.
func (d *Device) IsThermostat() bool {
	return d.Has(HeatControl)
}

// CanMeasurePower returns true if the device has powermeter functionality. Returns false otherwise.
func (d *Device) CanMeasurePower() bool {
	return d.Has(PowerSensor)
}

// CanMeasureTemp returns true if the device has thermometer functionality. Returns false otherwise.
func (d *Device) CanMeasureTemp() bool {
	return d.Has(TemperatureSensor)
}

// IsSwitch returns true if the device is recognized to be a switch and returns false otherwise.
func (d *Device) IsSwitch() bool {
	return d.Has(StateSwitch)
}

// CanRepeatDECT returns true if the device is capable of relaying Digital Enhanced Cordless Telecommunications (DECT) signals.
func (d *Device) CanRepeatDECT() bool {
	return d.Has(DECTRepeater)
}

// HasMicrophone returns true if the device has a microphone.
func (d *Device) HasMicrophone() bool {
	return d.Has(Microphone)
}

// HasHANFUNUnit returns true if the device has a HAN FUN unit.
func (d *Device) HasHANFUNUnit() bool {
	return d.Has(HANFUNUnit)
}

// IsSwitchableDevice returns true if the device is a switchable device/power plug/actor.
func (d *Device) IsSwitchableDevice() bool {
	return d.Has(SwitchableDevice)
}

// CanBeDimmed returns true if the device can be dimmed somehow (e.g. light intensity, height level, etc.).
func (d *Device) CanBeDimmed() bool {
	return d.Has(DimmableDevice)
}

// CanSetColors returns true if the device can set colors.
func (d *Device) CanSetColors() bool {
	return d.Has(ColorSettableDevice)
}

// CanMeasureHumidity returns true if the device has humidity functionality. Returns false otherwise.
func (d *Device) CanMeasureHumidity() bool {
	return d.Has(HumiditySensor)
}

// Has checks the passed capabilities and returns true iff the device supports all capabilities.
func (d *Device) Has(cs ...Capability) bool {
	for _, c := range cs {
		b := bitMasked{Functionbitmask: d.Functionbitmask}.hasMask(1 << uint(c))
		if !b {
			return false
		}
	}
	return true
}
