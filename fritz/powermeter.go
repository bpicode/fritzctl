package fritz

import "strconv"

// Powermeter models a power measurement
type Powermeter struct {
	Power  string `xml:"power"`  // Current power, refreshed approx every 2 minutes
	Energy string `xml:"energy"` // Absolute energy consumption since the device started operating
}

// FmtPowerW formats the value of p.Power as obtained on the http interface as a string, units are W.
func (p *Powermeter) FmtPowerW() string {
	f, err := strconv.ParseFloat(p.Power, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatFloat(f*0.001, 'f', -1, 64)
}

// FmtEnergyWh formats the value of p.Energy as obtained on the http interface as a string, units are Wh.
func (p *Powermeter) FmtEnergyWh() string {
	return p.Energy
}
