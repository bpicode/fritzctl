package chrono

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFormatSimpleOnSameDay tests the FormatSimple function.
func TestFormatSimpleOnSameDay(t *testing.T) {
	formatted := FormatSimple("1", time.Unix(0, 0))
	assert.NotEmpty(t, formatted)
	assert.Regexp(t, `^\d{2}:\d{2}:\d{2}`, formatted)
}

// TestFormatSimpleOnSameYear tests the FormatSimple function.
func TestFormatSimpleOnSameYear(t *testing.T) {
	earliest := time.Unix(0, 0)
	later := earliest.AddDate(0, 1, 2)
	formatted := FormatSimple(strconv.FormatInt(later.Unix(), 10), earliest)
	assert.NotEmpty(t, formatted)
	assert.Regexp(t, `^[[:alpha:]]{3} [[:alpha:]]{3} \d{1,2} \d{2}:\d{2}:\d{2}`, formatted)
}

// TestFormatSimpleOnSomeOtherYear tests the FormatSimple function.
func TestFormatSimpleOnSomeOtherYear(t *testing.T) {
	earliest := time.Unix(0, 0)
	later := earliest.AddDate(20, 4, 17)
	formatted := FormatSimple(strconv.FormatInt(later.Unix(), 10), earliest)
	assert.NotEmpty(t, formatted)
	assert.Regexp(t, `^[[:alpha:]]{3} [[:alpha:]]{3} \d{1,2} \d{2}:\d{2}:\d{2} \d{4}`, formatted)
}

// TestFormatSimpleOnMalformedInput tests the FormatSimple function.
func TestFormatSimpleOnMalformedInput(t *testing.T) {
	earliest := time.Now()
	formatted := FormatSimple("nonsense", earliest)
	assert.Empty(t, formatted)
}
