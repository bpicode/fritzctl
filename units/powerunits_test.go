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
