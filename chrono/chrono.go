package chrono

import (
	"strconv"
	"time"
)

// FormatEpochSecondString takes a string, parses to to an epoch second and formats it according to FormatCompact.
func FormatEpochSecondString(epoch string, ref time.Time) string {
	t, err := Unix(epoch)
	if err != nil {
		return ""
	}
	return FormatCompact(t, ref)
}

// Unix is equivalent to time.Unix with zero nanoseconds, where the string argument passed to this function is parsed
// as a base-10, 64-bit integer. It returns an error iff the argument could not be parsed.
func Unix(epoch string) (time.Time, error) {
	i, err := strconv.ParseInt(epoch, 10, 64)
	return time.Unix(i, 0), err
}

// FormatCompact formats a given time t to a short form, given a reference time ref. It particular:
// A simple time HH:MM:SS is displayed if t is in the same day as ref.
// Day, month and time is returned if t is in the same year as ref.
// Year, day, month and time is returned in all other cases.
func FormatCompact(t, ref time.Time) string {
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
