package fritz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormattingOfTemperaturesRegularRange tests formatting of the temperature values obtained by AHA interface.
func TestFormattingOfTemperaturesRegularRange(t *testing.T) {
	th := Thermostat{Measured: "47", Saving: "40", Goal: "42", Comfort: "44"}
	assert.Equal(t, "23.5", th.FmtMeasuredTemperature())
	assert.Equal(t, "20", th.FmtSavingTemperature())
	assert.Equal(t, "21", th.FmtGoalTemperature())
	assert.Equal(t, "22", th.FmtComfortTemperature())
}

// TestFormattingOfTemperaturesRegularRange tests formatting of the temperature values obtained by AHA interface.
func TestFormattingOfTemperaturesParseError(t *testing.T) {
	th := Thermostat{Measured: "assafsa", Saving: "dghdafhf", Goal: "dfahfh", Comfort: "rheeh"}
	assert.Equal(t, "", th.FmtMeasuredTemperature())
	assert.Equal(t, "", th.FmtSavingTemperature())
	assert.Equal(t, "", th.FmtGoalTemperature())
	assert.Equal(t, "", th.FmtComfortTemperature())
}

// TestFormattingOfTemperaturesSpecialValueOff tests formatting of the temperature values obtained by AHA interface.
func TestFormattingOfTemperaturesSpecialValueOff(t *testing.T) {
	th := Thermostat{Measured: "253", Saving: "253", Goal: "253", Comfort: "253"}
	assert.Equal(t, "OFF", th.FmtMeasuredTemperature())
	assert.Equal(t, "OFF", th.FmtSavingTemperature())
	assert.Equal(t, "OFF", th.FmtGoalTemperature())
	assert.Equal(t, "OFF", th.FmtComfortTemperature())
}

// TestFormattingOfTemperaturesSpecialValueOn tests formatting of the temperature values obtained by AHA interface.
func TestFormattingOfTemperaturesSpecialValueOn(t *testing.T) {
	th := Thermostat{Measured: "254", Saving: "254", Goal: "254", Comfort: "254"}
	assert.Equal(t, "ON", th.FmtMeasuredTemperature())
	assert.Equal(t, "ON", th.FmtSavingTemperature())
	assert.Equal(t, "ON", th.FmtGoalTemperature())
	assert.Equal(t, "ON", th.FmtComfortTemperature())
}

// TestFormattingOfTemperaturesOutOfRangeHigh tests formatting of the temperature values obtained by AHA interface.
func TestFormattingOfTemperaturesOutOfRangeHigh(t *testing.T) {
	th := Thermostat{Measured: "100", Saving: "110", Goal: "111", Comfort: "56"}
	assert.Equal(t, "28", th.FmtMeasuredTemperature())
	assert.Equal(t, "28", th.FmtSavingTemperature())
	assert.Equal(t, "28", th.FmtGoalTemperature())
	assert.Equal(t, "28", th.FmtComfortTemperature())
}

// TestFormattingOfTemperaturesOutOfRangeLow tests formatting of the temperature values obtained by AHA interface.
func TestFormattingOfTemperaturesOutOfRangeLow(t *testing.T) {
	th := Thermostat{Measured: "1", Saving: "2", Goal: "3", Comfort: "16"}
	assert.Equal(t, "8", th.FmtMeasuredTemperature())
	assert.Equal(t, "8", th.FmtSavingTemperature())
	assert.Equal(t, "8", th.FmtGoalTemperature())
	assert.Equal(t, "8", th.FmtComfortTemperature())
}
