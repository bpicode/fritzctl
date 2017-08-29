package manifest

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fritzAlwaysSuccess struct {
}

// SwitchOn always succeeds.
func (f *fritzAlwaysSuccess) SwitchOn(ain string) (string, error) {
	return "1", nil
}

// SwitchOff always succeeds.
func (f *fritzAlwaysSuccess) SwitchOff(ain string) (string, error) {
	return "0", nil
}

// ApplyTemperature always succeeds.
func (f *fritzAlwaysSuccess) ApplyTemperature(value float64, ain string) (string, error) {
	return strconv.FormatFloat(value*2.0, 'f', -1, 64), nil
}

// TestApplyViaAha tests the http interface applier.
func TestApplyViaAha(t *testing.T) {
	applier := AhaAPIApplier(&fritzAlwaysSuccess{})
	err := applier.Apply(
		&Plan{
			Switches:    []Switch{{Name: "s", State: true}},
			Thermostats: []Thermostat{{Name: "t", Temperature: 17.5}},
		},
		&Plan{
			Switches:    []Switch{{Name: "s", State: false}},
			Thermostats: []Thermostat{{Name: "t", Temperature: 20.5}},
		})
	assert.NoError(t, err)
}

// TestApplyViaAhaLargeSystem tests the http interface applier.
func TestApplyViaAhaLargeSystem(t *testing.T) {
	applier := AhaAPIApplier(&fritzAlwaysSuccess{})
	err := applier.Apply(
		&Plan{
			Switches: []Switch{
				{Name: "s1", State: true},
				{Name: "s2", State: true},
				{Name: "s3", State: true},
				{Name: "s4", State: false},
				{Name: "s5", State: false},
				{Name: "s6", State: true},
				{Name: "s7", State: true},
				{Name: "s8", State: true},
				{Name: "s9", State: false},
				{Name: "s10", State: true},
				{Name: "s11", State: true},
				{Name: "s12", State: false},
				{Name: "s13", State: false},
			},
			Thermostats: []Thermostat{
				{Name: "t1", Temperature: 17.5},
				{Name: "t2", Temperature: 18.5},
				{Name: "t3", Temperature: 19.5},
				{Name: "t4", Temperature: 21.5},
				{Name: "t5", Temperature: 22.5},
				{Name: "t6", Temperature: 23.0},
				{Name: "t7", Temperature: 34.0},
				{Name: "t8", Temperature: 26.0},
				{Name: "t9", Temperature: 27.5},
			},
		},
		&Plan{
			Switches: []Switch{
				{Name: "s1", State: false},
				{Name: "s2", State: false},
				{Name: "s3", State: true},
				{Name: "s4", State: true},
				{Name: "s5", State: true},
				{Name: "s6", State: false},
				{Name: "s7", State: false},
				{Name: "s8", State: true},
				{Name: "s9", State: true},
				{Name: "s10", State: true},
				{Name: "s11", State: false},
				{Name: "s12", State: false},
				{Name: "s13", State: true},
			},
			Thermostats: []Thermostat{
				{Name: "t1", Temperature: 27.5},
				{Name: "t2", Temperature: 19.5},
				{Name: "t3", Temperature: 17.5},
				{Name: "t4", Temperature: 25.5},
				{Name: "t5", Temperature: 21.5},
				{Name: "t6", Temperature: 24.0},
				{Name: "t7", Temperature: 24.0},
				{Name: "t8", Temperature: 16.0},
				{Name: "t9", Temperature: 17.5},
			},
		})
	assert.NoError(t, err)
}

type fritzAlwaysError struct {
}

// SwitchOn always returns an error.
func (f *fritzAlwaysError) SwitchOn(ain string) (string, error) {
	return "", errors.New("That didn't work")
}

// SwitchOff always returns an error.
func (f *fritzAlwaysError) SwitchOff(ain string) (string, error) {
	return "", errors.New("That didn't work")
}

// ApplyTemperature always returns an error.
func (f *fritzAlwaysError) ApplyTemperature(value float64, ain string) (string, error) {
	return "", errors.New("That didn't work")
}

// TestApplyViaAhaErrorByThermostat tests the http interface applier.
func TestApplyViaAhaErrorByThermostat(t *testing.T) {
	applier := AhaAPIApplier(&fritzAlwaysError{})
	err := applier.Apply(
		&Plan{
			Switches:    []Switch{{Name: "s", State: true}},
			Thermostats: []Thermostat{{Name: "t", Temperature: 17.5}},
		},
		&Plan{
			Switches:    []Switch{{Name: "s", State: true}},
			Thermostats: []Thermostat{{Name: "t", Temperature: 20.5}},
		})
	assert.Error(t, err)
}

// TestApplyViaAhaErrorBySwitch tests the http interface applier.
func TestApplyViaAhaErrorBySwitch(t *testing.T) {
	applier := AhaAPIApplier(&fritzAlwaysError{})
	err := applier.Apply(
		&Plan{
			Switches:    []Switch{{Name: "s", State: false}},
			Thermostats: []Thermostat{{Name: "t", Temperature: 20.5}},
		},
		&Plan{
			Switches:    []Switch{{Name: "s", State: true}},
			Thermostats: []Thermostat{{Name: "t", Temperature: 20.5}},
		})
	assert.Error(t, err)
}

// TestApplyViaAhaErrorBySwitch tests the http interface applier.
func TestApplyViaAhaErrorByMalformedPlan(t *testing.T) {
	applier := AhaAPIApplier(&fritzAlwaysError{})
	err := applier.Apply(
		&Plan{
			Switches:    []Switch{{Name: "s", State: false}},
			Thermostats: []Thermostat{{Name: "t", Temperature: 20.5}},
		},
		&Plan{
			Switches:    []Switch{{Name: "CCCCC", State: true}},
			Thermostats: []Thermostat{{Name: "YYYYYY", Temperature: 20.5}},
		})
	assert.Error(t, err)
}

// TestApplyViaAhaLargeSystemWithErrors tests the http interface applier.
func TestApplyViaAhaLargeSystemWithErrors(t *testing.T) {
	applier := AhaAPIApplier(&fritzAlwaysError{})
	err := applier.Apply(
		&Plan{
			Switches: []Switch{
				{Name: "s1", State: true},
				{Name: "s2", State: true},
				{Name: "s3", State: true},
				{Name: "s4", State: false},
				{Name: "s5", State: false},
				{Name: "s6", State: true},
				{Name: "s7", State: true},
				{Name: "s8", State: true},
				{Name: "s9", State: false},
				{Name: "s10", State: true},
				{Name: "s11", State: true},
				{Name: "s12", State: false},
				{Name: "s13", State: false},
			},
			Thermostats: []Thermostat{
				{Name: "t1", Temperature: 17.5},
				{Name: "t2", Temperature: 18.5},
				{Name: "t3", Temperature: 19.5},
				{Name: "t4", Temperature: 21.5},
				{Name: "t5", Temperature: 22.5},
				{Name: "t6", Temperature: 23.0},
				{Name: "t7", Temperature: 34.0},
				{Name: "t8", Temperature: 26.0},
				{Name: "t9", Temperature: 27.5},
			},
		},
		&Plan{
			Switches: []Switch{
				{Name: "s1", State: false},
				{Name: "s2", State: false},
				{Name: "s3", State: true},
				{Name: "s4", State: true},
				{Name: "s5", State: true},
				{Name: "s6", State: false},
				{Name: "s7", State: false},
				{Name: "s8", State: true},
				{Name: "s9", State: true},
				{Name: "s10", State: true},
				{Name: "s11", State: false},
				{Name: "s12", State: false},
				{Name: "s13", State: true},
			},
			Thermostats: []Thermostat{
				{Name: "t1", Temperature: 27.5},
				{Name: "t2", Temperature: 19.5},
				{Name: "t3", Temperature: 17.5},
				{Name: "t4", Temperature: 25.5},
				{Name: "t5", Temperature: 21.5},
				{Name: "t6", Temperature: 24.0},
				{Name: "t7", Temperature: 24.0},
				{Name: "t8", Temperature: 16.0},
				{Name: "t9", Temperature: 17.5},
			},
		})
	assert.Error(t, err)
}
