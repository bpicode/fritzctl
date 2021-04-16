package jsonapi

import (
	"strconv"
	"time"

	"github.com/bpicode/fritzctl/fritz"
)

// Mapper maps the XML to JSON model.
type Mapper interface {
	Convert([]fritz.Device) DeviceList
}

// NewMapper instantiates a Mapper.
func NewMapper() Mapper {
	return &mapper{}
}

type mapper struct {
}

// Convert translates a slice of Devices into a ThermostatList
func (m *mapper) Convert(ds []fritz.Device) DeviceList {
	l := DeviceList{Devices: []Device{}}
	l.NumberOfItems = len(ds)
	for _, d := range ds {
		l.Devices = append(l.Devices, m.convertOne(d))
	}
	return l
}

func (m *mapper) convertOne(src fritz.Device) Device {
	target := Device{}
	m.mapIdentifiers(&target, &src)
	m.mapProperties(&target, &src)
	m.mapMeasurements(&target, &src)
	m.mapState(&target, &src)
	return target
}

func (m *mapper) mapIdentifiers(target *Device, src *fritz.Device) {
	target.ID = src.Identifier
	target.InternalID = src.ID
	target.Name = src.Name
}

func (m *mapper) mapProperties(target *Device, src *fritz.Device) {
	props := &Properties{}
	mapVendor(props, src)
	mapLock(props, src)
	mapWarnings(props, src)
	target.Properties = props
}

func mapWarnings(target *Properties, src *fritz.Device) {
	if src.Thermostat.BatteryLow == "1" {
		target.Warnings = append(target.Warnings, "Battery is running on low capacity")
	}
	if fritz.HkrErrorDescriptions[src.Thermostat.ErrorCode] != "" {
		target.Warnings = append(target.Warnings, fritz.HkrErrorDescriptions[src.Thermostat.ErrorCode])
	}
}

func mapVendor(target *Properties, src *fritz.Device) {
	target.Vendor = &Vendor{
		Manufacturer:    src.Manufacturer,
		ProductName:     src.Productname,
		FirmwareVersion: src.Fwversion,
	}
}

var lockStateLookUp = map[string]string{
	"0": "UNLOCKED",
	"1": "LOCKED",
}

func mapLock(target *Properties, src *fritz.Device) {

	if src.Switch.Lock != "" || src.Switch.DeviceLock != "" {
		target.Lock = &Lock{
			HwLock: lockStateLookUp[src.Switch.DeviceLock],
			SwLock: lockStateLookUp[src.Switch.Lock],
		}
	}
	if src.Thermostat.Lock != "" || src.Thermostat.DeviceLock != "" {
		target.Lock = &Lock{
			HwLock: lockStateLookUp[src.Thermostat.DeviceLock],
			SwLock: lockStateLookUp[src.Thermostat.Lock],
		}
	}
}

var alertSignalLookup = map[string]string{
	"0": "OFF",
	"1": "ON",
}

func (m *mapper) mapMeasurements(target *Device, src *fritz.Device) {
	meas := &Measurements{}
	if src.Temperature.Celsius != "" {
		meas.Temperature = src.Temperature.FmtCelsius()
	}
	if src.Powermeter.Power != "" {
		meas.PowerConsumption = src.Powermeter.FmtPowerW()
	}
	if src.Powermeter.Energy != "" {
		meas.EnergyConsumption = src.Powermeter.FmtEnergyWh()
	}
	meas.AlertSignal = alertSignalLookup[src.AlertSensor.State]
	meas.ButtonLastPressed = src.Button.LastPressed()
	target.Measurements = meas
}

var switchStateLookup = map[string]string{
	"0": "OFF",
	"1": "ON",
}

func (m *mapper) mapState(target *Device, src *fritz.Device) {
	st := &State{}
	st.Connected = src.Present == 1
	st.Switch = switchStateLookup[src.Switch.State]
	if src.IsThermostat() {
		m.mapThermostat(st, src)
	}
	target.State = st
}

var windowStateLookup = map[string]string{
	"0": "CLOSED",
	"1": "OPEN",
}

func (m *mapper) mapThermostat(target *State, src *fritz.Device) {
	tc := &TemperatureControl{}
	tc.Goal = src.Thermostat.FmtGoalTemperature()
	tc.Saving = src.Thermostat.FmtSavingTemperature()
	tc.Comfort = src.Thermostat.FmtComfortTemperature()
	if src.Thermostat.NextChange.Goal != "" {
		m.mapNextChange(tc, src)
	}
	tc.BatteryLow = src.Thermostat.BatteryLow
	tc.BatteryChargeLevel = src.Thermostat.BatteryChargeLevel
	tc.Window = windowStateLookup[src.Thermostat.WindowOpen]
	if src.Thermostat.WindowOpenEnd != 0 {
		tc.WindowOpenEnd = time.Unix(src.Thermostat.WindowOpenEnd, 0).Format((time.RFC3339))
	}
	tc.Boost = src.Thermostat.Boost
	if src.Thermostat.BoostEnd != 0 {
		tc.BoostEnd = time.Unix(src.Thermostat.BoostEnd, 0).Format((time.RFC3339))
	}
	target.TemperatureControl = tc

	switch src.Thermostat.BatteryLow {
	case "0":
		target.BatteryState = "OK"
	case "1":
		target.BatteryState = "LOW"
	}

	if chargeLevelPct, err := strconv.ParseFloat(src.Thermostat.BatteryChargeLevel, 64); err == nil {
		chargeLevelNormalized := chargeLevelPct / 100.0
		target.BatteryChargeLevel = strconv.FormatFloat(chargeLevelNormalized, 'f', -1, 64)
	}
}

func (m *mapper) mapNextChange(target *TemperatureControl, src *fritz.Device) {
	nc := &NextChange{}
	nc.Goal = src.Thermostat.NextChange.FmtGoalTemperature()
	t, err := strconv.ParseInt(src.Thermostat.NextChange.TimeStamp, 10, 64)
	if err == nil {
		nc.At = time.Unix(t, 0).Format(time.RFC3339)
	}
	target.NextChange = nc
}
