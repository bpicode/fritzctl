package fritz

import "strconv"

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

// NextChange corresponds to the next HKR switch event.
type NextChange struct {
	TimeStamp string `xml:"endperiod"` // Timestamp (epoch time) when the next temperature switch is scheduled.
	Goal      string `xml:"tchange"`   // The temperature to switch to. Same unit convention as in Thermostat.Measured.
}

// codebeat:enable[TOO_MANY_IVARS]

// FmtMeasuredTemperature formats the value of t.Measured as obtained on the xml-over http interface to a floating
// point string, units in °C.
// If the value cannot be parsed an empty string is returned.
// If the value if 254 or 253, "ON" or "OFF" is returned, respectively.
// If the value is greater (less) than 56 (16) a cut-off "28" ("8") is returned.
func (t *Thermostat) FmtMeasuredTemperature() string {
	return t.fmtTemperature(t.Measured)
}

// FmtGoalTemperature formats the value of t.Goal as obtained on the xml-over http interface to a floating
// point string, units in °C.
// If the value cannot be parsed an empty string is returned.
// If the value if 254 or 253, "ON" or "OFF" is returned, respectively.
// If the value is greater (less) than 56 (16) a cut-off "28" ("8") is returned.
func (t *Thermostat) FmtGoalTemperature() string {
	return t.fmtTemperature(t.Goal)
}

// FmtSavingTemperature formats the value of t.Saving as obtained on the xml-over http interface to a floating
// point string, units in °C.
// If the value cannot be parsed an empty string is returned.
// If the value if 254 or 253, "ON" or "OFF" is returned, respectively.
// If the value is greater (less) than 56 (16) a cut-off "28" ("8") is returned.
func (t *Thermostat) FmtSavingTemperature() string {
	return t.fmtTemperature(t.Saving)
}

// FmtComfortTemperature formats the value of t.Comfort as obtained on the xml-over http interface to a floating
// point string, units in °C.
// If the value cannot be parsed an empty string is returned.
// If the value if 254 or 253, "ON" or "OFF" is returned, respectively.
// If the value is greater (less) than 56 (16) a cut-off "28" ("8") is returned.
func (t *Thermostat) FmtComfortTemperature() string {
	return t.fmtTemperature(t.Comfort)
}

func (t *Thermostat) fmtTemperature(th string) string {
	f, err := strconv.ParseFloat(th, 64)
	if err != nil {
		return ""
	}
	var str string
	switch {
	case f == 254:
		str = "ON"
	case f == 253:
		str = "OFF"
	case f < 16:
		str = "8"
	case f > 56:
		str = "28"
	default:
		str = strconv.FormatFloat(f*0.5, 'f', -1, 64)
	}
	return str
}
