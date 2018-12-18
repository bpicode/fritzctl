package fritz

import (
	"strconv"
	"time"
)

// Button collects data from devices that have a pressable button.
type Button struct {
	LastPressedTimestamp string `xml:"lastpressedtimestamp"` // Timestamp (in epoch seconds) when the button was last pressed. "0" or "" if unknown.
}

// LastPressed returns the time when the button was last pressed. It returns nil on absence of this information or upon parsing errors.
func (b *Button) LastPressed() *time.Time {
	unix, err := strconv.ParseInt(b.LastPressedTimestamp, 10, 64)
	if err != nil {
		return nil
	}
	t := time.Unix(unix, 0)
	return &t
}

// FmtLastPressedCompact returns a compact format of the last pressed timestamp.
func (b *Button) FmtLastPressedCompact(ref time.Time) string {
	t := b.LastPressed()
	if t == nil {
		return ""
	}
	ry, rm, rd := ref.Date()
	ty, tm, td := t.Date()
	if ry != ty {
		return t.Format("Mon Jan 2 15:04:05 2006")
	}
	if rm != tm || rd != td {
		return t.Format("Mon Jan 2 15:04:05")
	}
	return t.Format("15:04:05")
}
