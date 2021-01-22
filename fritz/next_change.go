package fritz

import (
	"time"
)

// NextChange corresponds to the next HKR switch event.
type NextChange struct {
	TimeStamp string `xml:"endperiod"` // Timestamp (epoch time) when the next temperature switch is scheduled.
	Goal      string `xml:"tchange"`   // The temperature to switch to. Same unit convention as in Thermostat.Measured.
}

// FmtGoalTemperature formats the value of t.Goal as obtained on the xml-over http interface to a floating
// point string, units in °C.
// If the value cannot be parsed an empty string is returned.
// If the value if 255, 254 or 253, "?", "ON" or "OFF" is returned, respectively.
// If the value is greater (less) than 56 (16) a cut-off "28" ("8") is returned.
func (n *NextChange) FmtGoalTemperature() string {
	return fmtTemperatureHkr(n.Goal)
}

// FmtTimestamp formats the epoch timestamp into a compact readable form. See FmtEpochSecondString.
func (n *NextChange) FmtTimestamp(ref time.Time) string {
	if n.TimeStamp == "0" {
		return ""
	}
	return FmtEpochSecondString(n.TimeStamp, ref)
}
