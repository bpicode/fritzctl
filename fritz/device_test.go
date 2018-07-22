package fritz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParsingFunctionBitMask verifies correctness of the function bit.mask interpretation.
func TestParsingFunctionBitMask(t *testing.T) {
	isSwitch := func(d *Device) bool { return d.IsSwitch() }
	isThermo := func(d *Device) bool { return d.IsThermostat() }
	canMeasT := func(d *Device) bool { return d.CanMeasureTemp() }
	canMeasP := func(d *Device) bool { return d.CanMeasurePower() }
	for _, tc := range []struct {
		name   string
		mask   string
		fct    func(*Device) bool
		expect bool
	}{
		{name: "nonsense is not switch", mask: "nonsense", fct: isSwitch, expect: false},
		{name: "nonsense is not thermostat", mask: "nonsense", fct: isThermo, expect: false},
		{name: "nonsense cannot measure temperature", mask: "nonsense", fct: canMeasT, expect: false},
		{name: "nonsense cannot measure power", mask: "nonsense", fct: canMeasP, expect: false},
		{name: "2944 is switch", mask: "2944", fct: isSwitch, expect: true},
		{name: "2944 is not thermostat", mask: "2944", fct: isThermo, expect: false},
		{name: "2944 can measure temperature", mask: "2944", fct: canMeasT, expect: true},
		{name: "2944 can measure power", mask: "2944", fct: canMeasP, expect: true},
		{name: "320 is not switch", mask: "320", fct: isSwitch, expect: false},
		{name: "320 is thermostat", mask: "320", fct: isThermo, expect: true},
		{name: "320 can measure temperature", mask: "320", fct: canMeasT, expect: true},
		{name: "320 cannot measure power", mask: "320", fct: canMeasP, expect: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			device := &Device{Functionbitmask: tc.mask}
			assert.Equal(t, tc.expect, tc.fct(device))
		})
	}
}
