package manifest

import (
	"strconv"

	"github.com/bpicode/fritzctl/fritz"
)

func convertDevicelist(l *fritz.Devicelist) *Plan {
	var p Plan
	for _, s := range l.Switches() {
		p.Switches = append(p.Switches, convertSwitch(&s))
	}
	for _, t := range l.Thermostats() {
		p.Thermostats = append(p.Thermostats, convertThermostat(&t))
	}

	return &p
}

func convertSwitch(d *fritz.Device) Switch {
	var s Switch
	s.Name = d.Name
	s.State, _ = strconv.ParseBool(d.Switch.State)
	return s
}

func convertThermostat(d *fritz.Device) Thermostat {
	var t Thermostat
	t.Name = d.Name
	goalTimesTwo, _ := strconv.ParseFloat(d.Thermostat.Goal, 64)
	t.Temperature = goalTimesTwo * 0.5
	return t
}
