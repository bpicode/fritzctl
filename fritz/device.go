package fritz

// Devicelist wraps a list of devices.
type Devicelist struct {
	Devices []Device `xml:"device"`
}

// Device models a smart home device.
type Device struct {
	Identifier      string      `xml:"identifier,attr"`
	ID              string      `xml:"id,attr"`
	Functionbitmask string      `xml:"functionbitmask,attr"`
	Fwversion       string      `xml:"fwversion,attr"`
	Manufacturer    string      `xml:"manufacturer,attr"`
	Productname     string      `xml:"productname,attr"`
	Present         int         `xml:"present"`
	Name            string      `xml:"name"`
	Switch          Switch      `xml:"switch"`
	Powermeter      Powermeter  `xml:"powermeter"`
	Temperature     Temperature `xml:"temperature"`
}
