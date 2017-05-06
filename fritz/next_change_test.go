package fritz

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFormattingOfGoalNextChange tests formatting of the temperature values obtained by AHA interface.
func TestFormattingOfGoalNextChange(t *testing.T) {
	assert.Equal(t, "23.5", (&NextChange{Goal: "47"}).FmtGoalTemperature())
	assert.Equal(t, "?", (&NextChange{Goal: "255"}).FmtGoalTemperature())
	assert.Equal(t, "ON", (&NextChange{Goal: "254"}).FmtGoalTemperature())
	assert.Equal(t, "OFF", (&NextChange{Goal: "253"}).FmtGoalTemperature())
}

// TestFmtTimestampOnSameDay tests the FmtTimestamp function.
func TestFmtTimestampOnSameDay(t *testing.T) {
	formatted := (&NextChange{TimeStamp: "1"}).FmtTimestamp(time.Unix(0, 0))
	assert.NotEmpty(t, formatted)
	assert.Regexp(t, `^\d{2}:\d{2}:\d{2}`, formatted)
}

// TestFmtTimestampOnSameYear tests the FmtTimestamp function.
func TestFmtTimestampOnSameYear(t *testing.T) {
	earliest := time.Unix(0, 0)
	later := earliest.AddDate(0, 1, 2)
	formatted := (&NextChange{TimeStamp: strconv.FormatInt(later.Unix(), 10)}).FmtTimestamp(earliest)
	assert.NotEmpty(t, formatted)
	assert.Regexp(t, `^[[:alpha:]]{3} [[:alpha:]]{3} \d{1,2} \d{2}:\d{2}:\d{2}`, formatted)
}

// TestFmtTimestampOnSomeOtherYear tests the FmtTimestamp function.
func TestFmtTimestampOnSomeOtherYear(t *testing.T) {
	earliest := time.Unix(0, 0)
	later := earliest.AddDate(20, 4, 17)
	formatted := (&NextChange{TimeStamp: strconv.FormatInt(later.Unix(), 10)}).FmtTimestamp(earliest)
	assert.NotEmpty(t, formatted)
	assert.Regexp(t, `^[[:alpha:]]{3} [[:alpha:]]{3} \d{1,2} \d{2}:\d{2}:\d{2} \d{4}`, formatted)
}

// TestFmtTimestampOnMalformedInput tests the FmtTimestamp function.
func TestFmtTimestampOnMalformedInput(t *testing.T) {
	earliest := time.Now()
	formatted := (&NextChange{TimeStamp: "nonsense"}).FmtTimestamp(earliest)
	assert.Empty(t, formatted)
}

// TestFmtTimestampWithSpecialValueZero tests the FmtTimestamp function.
func TestFmtTimestampWithSpecialValueZero(t *testing.T) {
	formatted := (&NextChange{TimeStamp: "0"}).FmtTimestamp(time.Unix(0, 0))
	assert.Zero(t, formatted)
}
