package fritz

import (
	"strconv"
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

// FmtTimestamp formats the epoch timestamp into a compact readable form. See fmtEpochSecondString.
func (n *NextChange) FmtTimestamp(ref time.Time) string {
	if n.TimeStamp == "0" {
		return ""
	}
	return n.fmtEpochSecondString(ref)
}

// fmtEpochSecondString takes a string, parses to to an epoch second and formats it according to fmtCompact.
func (n *NextChange) fmtEpochSecondString(ref time.Time) string {
	t, err := n.unix(n.TimeStamp)
	if err != nil {
		return ""
	}
	return n.fmtCompact(t, ref)
}

// unix is equivalent to time.unix with zero nanoseconds, where the string argument passed to this function is parsed
// as a base-10, 64-bit integer. It returns an error iff the argument could not be parsed.
func (n *NextChange) unix(epoch string) (time.Time, error) {
	i, err := strconv.ParseInt(epoch, 10, 64)
	return time.Unix(i, 0), err
}

// fmtCompact formats a given time t to a short form, given a reference time ref. It particular:
// A simple time HH:MM:SS is displayed if t is in the same day as ref.
// Day, month and time is returned if t is in the same year as ref.
// Year, day, month and time is returned in all other cases.
func (n *NextChange) fmtCompact(t, ref time.Time) string {
	refYear, refMonth, refDay := ref.Date()
	tYear, tMonth, tDay := t.Date()
	if refYear != tYear {
		return t.Format("Mon Jan 2 15:04:05 2006")
	}
	if refMonth != tMonth || refDay != tDay {
		return t.Format("Mon Jan 2 15:04:05")
	}
	return t.Format("15:04:05")
}
