package fritz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParsingFunctionBitMask verifies correctness of the function bit.mask interpretation.
func TestParsingFunctionBitMask(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mask   string
		fct    func(*Device) bool
		expect bool
	}{
		{name: "nonsense is not switch", mask: "nonsense", fct: (*Device).IsSwitch, expect: false},
		{name: "nonsense is not thermostat", mask: "nonsense", fct: (*Device).IsThermostat, expect: false},
		{name: "nonsense cannot measure temperature", mask: "nonsense", fct: (*Device).CanMeasureTemp, expect: false},
		{name: "nonsense cannot measure power", mask: "nonsense", fct: (*Device).CanMeasurePower, expect: false},
		{name: "2944 is switch", mask: "2944", fct: (*Device).IsSwitch, expect: true},
		{name: "2944 is not thermostat", mask: "2944", fct: (*Device).IsThermostat, expect: false},
		{name: "2944 can measure temperature", mask: "2944", fct: (*Device).CanMeasureTemp, expect: true},
		{name: "2944 can measure power", mask: "2944", fct: (*Device).CanMeasurePower, expect: true},
		{name: "320 is not switch", mask: "320", fct: (*Device).IsSwitch, expect: false},
		{name: "320 is thermostat", mask: "320", fct: (*Device).IsThermostat, expect: true},
		{name: "320 can measure temperature", mask: "320", fct: (*Device).CanMeasureTemp, expect: true},
		{name: "320 cannot measure power", mask: "320", fct: (*Device).CanMeasurePower, expect: false},
		{name: "320 cannot repeat dect", mask: "320", fct: (*Device).CanRepeatDECT, expect: false},
		{name: "320 cannot repeat trigger alerts", mask: "320", fct: (*Device).HasAlertSensor, expect: false},
		{name: "320 has no microphone", mask: "320", fct: (*Device).HasMicrophone, expect: false},
		{name: "320 has no hanfun unit", mask: "320", fct: (*Device).HasHANFUNUnit, expect: false},
		{name: "320 does not speak hanfun protocol", mask: "320", fct: (*Device).IsHANFUNCompatible, expect: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			device := &Device{Functionbitmask: tc.mask}
			assert.Equal(t, tc.expect, tc.fct(device))
		})
	}
}
