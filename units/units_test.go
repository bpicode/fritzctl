package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseRegularFloat unit test.
func TestParseRegularFloat(t *testing.T) {
	str := ParseFloatAndScale("7500", 0.001)
	assert.Equal(t, "7.5", str)
}

// TestParseIrregularFloat unit test.
func TestParseIrregularFloat(t *testing.T) {
	str := ParseFloatAndScale("xx", 0.001)
	assert.Equal(t, "", str)
}

// TestParseAndAddRegularFloat unit test.
func TestParseAndAddRegularFloat(t *testing.T) {
	str := ParseFloatAddAndScale("190", "0", 0.1)
	assert.Equal(t, "19", str)
}

// TestParseAndAddOneIrregularFloat unit test.
func TestParseAndAddOneIrregularFloat(t *testing.T) {
	str := ParseFloatAddAndScale("XX", "0", 0.1)
	assert.Equal(t, "", str)
}

// TestParseAndAddAnotherIrregularFloat unit test.
func TestParseAndAddAnotherIrregularFloat(t *testing.T) {
	str := ParseFloatAddAndScale("123", "ll", 0.1)
	assert.Equal(t, "", str)
}
