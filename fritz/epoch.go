package fritz

import (
	"strconv"
	"time"
)

// FmtEpochSecond takes a int64, and formats it according to FmtCompact.
func FmtEpochSecond(t int64, ref time.Time) string {
	return FmtCompact(time.Unix(1, 0), ref)
}

// FmtEpochSecondString takes a string, parses to to an epoch second and formats it according to FmtCompact.
func FmtEpochSecondString(timeStamp string, ref time.Time) string {
	t, err := EpochToUnix(timeStamp)
	if err != nil {
		return ""
	}
	return FmtCompact(t, ref)
}

// EpochToUnix is equivalent to time.unix with zero nanoseconds, where the string argument passed to this function is parsed
// as a base-10, 64-bit integer. It returns an error iff the argument could not be parsed.
func EpochToUnix(epoch string) (time.Time, error) {
	i, err := strconv.ParseInt(epoch, 10, 64)
	return time.Unix(i, 0), err
}

// FmtCompact formats a given time t to a short form, given a reference time ref. It particular:
// A simple time HH:MM:SS is displayed if t is in the same day as ref.
// Day, month and time is returned if t is in the same year as ref.
// Year, day, month and time is returned in all other cases.
func FmtCompact(t, ref time.Time) string {
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
