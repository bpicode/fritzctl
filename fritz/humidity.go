package fritz

// Humidity models a humidity measurement.
type Humidity struct {
	RelHumidity string `xml:"rel_humidity"` // Relative humidity measured as full percentile.
}

// FmtRelativeHumidity formats the value of p.RelHumidity as obtained on the http interface as a string, units are percentile.
func (p *Humidity) FmtRelativeHumidity() string {
	return p.RelHumidity
}
