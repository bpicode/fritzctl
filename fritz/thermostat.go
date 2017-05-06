package fritz

// Thermostat models the "HKR" device.
// codebeat:disable[TOO_MANY_IVARS]
type Thermostat struct {
	Measured   string     `xml:"tist"`       // Measured temperature.
	Goal       string     `xml:"tsoll"`      // Desired temperature, user controlled.
	Saving     string     `xml:"absenk"`     // Energy saving temperature.
	Comfort    string     `xml:"komfort"`    // Comfortable temperature.
	NextChange NextChange `xml:"nextchange"` // The next scheduled temperature change.
	Lock       string     `xml:"lock"`       // Switch locked (box defined)? 1/0 (empty if not known or if there was an error).
	DeviceLock string     `xml:"devicelock"` // Switch locked (device defined)? 1/0 (empty if not known or if there was an error).
	ErrorCode  string     `xml:"errorcode"`  // Error codes: 0 = OK, 1 = ... see https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.
	BatteryLow string     `xml:"batterylow"` // "0" if the battery is OK, "1" if it is running low on capacity.
}

// codebeat:enable[TOO_MANY_IVARS]

// FmtMeasuredTemperature formats the value of t.Measured as obtained on the xml-over http interface to a floating
// point string, units in °C.
// If the value cannot be parsed an empty string is returned.
// If the value if 255, 254 or 253, "?", "ON" or "OFF" is returned, respectively.
// If the value is greater (less) than 56 (16) a cut-off "28" ("8") is returned.
func (t *Thermostat) FmtMeasuredTemperature() string {
	return fmtTemperatureHkr(t.Measured)
}

// FmtGoalTemperature formats the value of t.Goal as obtained on the xml-over http interface to a floating
// point string, units in °C.
// If the value cannot be parsed an empty string is returned.
// If the value if 255, 254 or 253, "?", "ON" or "OFF" is returned, respectively.
// If the value is greater (less) than 56 (16) a cut-off "28" ("8") is returned.
func (t *Thermostat) FmtGoalTemperature() string {
	return fmtTemperatureHkr(t.Goal)
}

// FmtSavingTemperature formats the value of t.Saving as obtained on the xml-over http interface to a floating
// point string, units in °C.
// If the value cannot be parsed an empty string is returned.
// If the value if 255, 254 or 253, "?", "ON" or "OFF" is returned, respectively.
// If the value is greater (less) than 56 (16) a cut-off "28" ("8") is returned.
func (t *Thermostat) FmtSavingTemperature() string {
	return fmtTemperatureHkr(t.Saving)
}

// FmtComfortTemperature formats the value of t.Comfort as obtained on the xml-over http interface to a floating
// point string, units in °C.
// If the value cannot be parsed an empty string is returned.
// If the value if 255, 254 or 253, "?", "ON" or "OFF" is returned, respectively.
// If the value is greater (less) than 56 (16) a cut-off "28" ("8") is returned.
func (t *Thermostat) FmtComfortTemperature() string {
	return fmtTemperatureHkr(t.Comfort)
}
