package chrono

import (
	"strconv"
	"time"
)

// FormatEpochSecondString takes a string, parses to to an epoch second and formats it according to the following rules:
// A simple time HH:MM:SS is displayed if the parsed date is today.
// Day, month and time is returned if the parsed date is the current year.
// Year, day, month and time is returned in all other cases.
func FormatEpochSecondString(epoch string, ref time.Time) string {
	i, err := strconv.ParseInt(epoch, 10, 64)
	if err != nil {
		return ""
	}
	return FormatEpochSecondI64(i, ref)
}

// FormatEpochSecondI64 follows the same conventions as FormatEpochSecondString.
func FormatEpochSecondI64(epoch int64, ref time.Time) string {
	t := time.Unix(epoch, 0)
	if ref.Day() == t.Day() {
		return t.Format("15:04:05")
	}
	if ref.Year() == t.Year() {
		return t.Format("Mon Jan 2 15:04:05")
	}
	return t.Format("Mon Jan 2 15:04:05 2006")
}

