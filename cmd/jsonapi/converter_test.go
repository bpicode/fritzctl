package jsonapi

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/stretchr/testify/assert"
)

// Test_thMapper_Convert probes the converter.
func Test_thMapper_Convert(t *testing.T) {
	m := NewMapper()
	devices := []fritz.Device{
		simpleHkr(),
		problematicPlungerHkr(),
		simpleSwitch(),
	}
	l := m.Convert(devices)
	bs, err := json.Marshal(l)
	assert.NoError(t, err)
	fmt.Println(string(bs))
	assert.Equal(t, l.NumberOfItems, len(devices))
}

func simpleHkr() fritz.Device {
	return fritz.Device{
		Name:            "myhkr",
		Functionbitmask: "320",
		Thermostat: fritz.Thermostat{
			NextChange: fritz.NextChange{TimeStamp: "121441515", Goal: "35"},
			BatteryLow: "0",
		},
	}
}

func problematicPlungerHkr() fritz.Device {
	return fritz.Device{
		Name:            "myhkrwitherr",
		Functionbitmask: "320",
		Thermostat:      fritz.Thermostat{ErrorCode: "2", DeviceLock: "1", BatteryLow: "1"},
	}
}

func simpleSwitch() fritz.Device {
	return fritz.Device{
		Name:            "myswitch",
		Functionbitmask: "2944",
		Switch:          fritz.Switch{State: "1", DeviceLock: "0", Lock: "0"},
		Powermeter:      fritz.Powermeter{Energy: "0", Power: "0"},
		Temperature:     fritz.Temperature{Celsius: "222"},
	}
}
