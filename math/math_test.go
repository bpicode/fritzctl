package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseRegularFloat test float parsing + scaling (regular).
func TestParseRegularFloat(t *testing.T) {
	str := ParseFloatAndScale("7500", 0.001)
	assert.Equal(t, "7.5", str)
}

// TestParseIrregularFloat test float parsing + scaling (malformed input).
func TestParseIrregularFloat(t *testing.T) {
	str := ParseFloatAndScale("xx", 0.001)
	assert.Equal(t, "", str)
}

// TestRounding tests rounding.
func TestRounding(t *testing.T) {
	assert.Equal(t, int64(1), Round(0.5))
	assert.Equal(t, int64(0), Round(0.4))
	assert.Equal(t, int64(0), Round(0.1))
	assert.Equal(t, int64(0), Round(-0.1))
	assert.Equal(t, int64(0), Round(-0.499))
	assert.Equal(t, int64(156), Round(156))
	assert.Equal(t, int64(3), Round(3.14))
	assert.Equal(t, int64(4), Round(3.54))
}
