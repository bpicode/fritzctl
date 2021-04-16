package fritz

import "time"

// HkrErrorDescriptions has a translation of error code to a warning/error/status description.
var HkrErrorDescriptions = map[string]string{
	"":  "",
	"0": "",
	"1": " Thermostat adjustment not possible. Is the device mounted correctly?",
	"2": " Valve plunger cannot be driven far enough. Possible solutions: Open and close the plunger a couple of times by hand. Check if the battery is too weak.",
	"3": " Valve plunger cannot be moved. Is it blocked?",
	"4": " Preparing installation.",
	"5": " Device in mode 'INSTALLATION'. It can be mounted now.",
	"6": " Device is adjusting to the valve plunger.",
}

// Thermostat models the "HKR" device.
// codebeat:disable[TOO_MANY_IVARS]
type Thermostat struct {
	Measured           string     `xml:"tist"`                    // Measured temperature.
	Goal               string     `xml:"tsoll"`                   // Desired temperature, user controlled.
	Saving             string     `xml:"absenk"`                  // Energy saving temperature.
	Comfort            string     `xml:"komfort"`                 // Comfortable temperature.
	NextChange         NextChange `xml:"nextchange"`              // The next scheduled temperature change.
	Lock               string     `xml:"lock"`                    // Switch locked (box defined)? 1/0 (empty if not known or if there was an error).
	DeviceLock         string     `xml:"devicelock"`              // Switch locked (device defined)? 1/0 (empty if not known or if there was an error).
	ErrorCode          string     `xml:"errorcode"`               // Error codes: 0 = OK, 1 = ... see https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.
	BatteryLow         string     `xml:"batterylow"`              // "0" if the battery is OK, "1" if it is running low on capacity.
	WindowOpen         string     `xml:"windowopenactiv"`         // "1" if detected an open window (usually turns off heating), "0" if not.
	BatteryChargeLevel string     `xml:"battery"`                 // Battery charge level in percent.
	WindowOpenEnd      int64      `xml:"windowopenactiveendtime"` // Scheduled end of window-open state (seconds since 1970)
	Boost              bool       `xml:"boostactive"`             // true if boost mode is active, false if not.
	BoostEnd           int64      `xml:"boostactiveendtime"`      // Scheduled end of boost time (seconds since 1970)
	Holiday            bool       `xml:"holidayactive"`           // true if device is in holiday-mode, false if not.
	Summer             bool       `xml:"summeractive"`            // true if device is in summer mode (heating off), false if not.
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

// FmtWindowOpenEndTimestamp formats the epoch timestamp into a compact readable form. See FmtEpochSecondString.
func (t *Thermostat) FmtWindowOpenEndTimestamp(ref time.Time) string {
	if t.WindowOpenEnd == 0 {
		return ""
	}
	return FmtEpochSecond(t.WindowOpenEnd, ref)
}

// FmtBoostEndTimestamp formats the epoch timestamp into a compact readable form. See FmtEpochSecondString.
func (t *Thermostat) FmtBoostEndTimestamp(ref time.Time) string {
	if t.BoostEnd == 0 {
		return ""
	}
	return FmtEpochSecond(t.WindowOpenEnd, ref)
}
