package fritz

// Thermostat models the "HKR" device.
// codebeat:disable[TOO_MANY_IVARS]
type Thermostat struct {
	Measured   string `xml:"tist"`       // Measured temperature.
	Goal       string `xml:"tsoll"`      // Desired temperature, user controlled.
	Saving     string `xml:"absenk"`     // Energy saving temperature.
	Comfort    string `xml:"komfort"`    // Comfortable temperature.
	Lock       string `xml:"lock"`       // Switch locked (box defined)? 1/0 (empty if not known or if there was an error).
	DeviceLock string `xml:"devicelock"` // Switch locked (device defined)? 1/0 (empty if not known or if there was an error).
}

// codebeat:enable[TOO_MANY_IVARS]
