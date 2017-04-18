package chrono

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFormatEpochSecondStringOnSameDay tests the FormatEpochSecondString function.
func TestFormatEpochSecondStringOnSameDay(t *testing.T) {
	formatted := FormatEpochSecondString("1", time.Unix(0, 0))
	assert.NotEmpty(t, formatted)
	assert.Regexp(t, `^\d{2}:\d{2}:\d{2}`, formatted)
}

// TestFormatEpochSecondStringOnSameYear tests the FormatEpochSecondString function.
func TestFormatEpochSecondStringOnSameYear(t *testing.T) {
	earliest := time.Unix(0, 0)
	later := earliest.AddDate(0, 1, 2)
	formatted := FormatEpochSecondString(strconv.FormatInt(later.Unix(), 10), earliest)
	assert.NotEmpty(t, formatted)
	assert.Regexp(t, `^[[:alpha:]]{3} [[:alpha:]]{3} \d{1,2} \d{2}:\d{2}:\d{2}`, formatted)
}

// TestFormatEpochSecondStringOnSomeOtherYear tests the FormatEpochSecondString function.
func TestFormatEpochSecondStringOnSomeOtherYear(t *testing.T) {
	earliest := time.Unix(0, 0)
	later := earliest.AddDate(20, 4, 17)
	formatted := FormatEpochSecondString(strconv.FormatInt(later.Unix(), 10), earliest)
	assert.NotEmpty(t, formatted)
	assert.Regexp(t, `^[[:alpha:]]{3} [[:alpha:]]{3} \d{1,2} \d{2}:\d{2}:\d{2} \d{4}`, formatted)
}

// TestFormatEpochSecondStringOnMalformedInput tests the FormatEpochSecondString function.
func TestFormatEpochSecondStringOnMalformedInput(t *testing.T) {
	earliest := time.Now()
	formatted := FormatEpochSecondString("nonsense", earliest)
	assert.Empty(t, formatted)
}
