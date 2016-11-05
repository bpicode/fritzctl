package fritz

// Temperature models a temperature measurement.
type Temperature struct {
	Celsius string `xml:"celsius"` // Current power, refreshed approx every 2 minutes
	Offset  string `xml:"offset"`  // Absolute energy consuption since the device started operating
}
