package fritz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormattingOfTemperatures tests formatting of the temperature values obtained by AHA interface.
func TestFormattingOfTemperatures(t *testing.T) {
	temperature := Temperature{Offset: "0", Celsius: "235"}
	assert.Equal(t, "23.5", temperature.FmtCelsius())
	assert.Equal(t, "0", temperature.FmtOffset())
}

// TestFormattingOfTemperaturesParseError tests formatting of the temperature values obtained by AHA interface.
func TestFormattingOfTemperaturesWithNonsenseValues(t *testing.T) {
	temperature := Temperature{Offset: "xxxx", Celsius: "yyyy"}
	assert.Zero(t, temperature.FmtCelsius())
	assert.Zero(t, temperature.FmtOffset())
}
